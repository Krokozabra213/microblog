package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"microblog/internal/handlers"
	"microblog/internal/logger"
	"microblog/internal/queue"
	postservice "microblog/internal/service/post-service"
	userservice "microblog/internal/service/user-service"
	"microblog/internal/storage"
)

// без конфига
const (
	gracefullTimeout  = 5 * time.Second
	eventBufferLogger = 100
	likeQueueBuffer   = 100
	serverAddr        = ":8080"
)

func main() {
	// Инициализация storage слоя
	userStorage := storage.NewUserStorage()
	postStorage := storage.NewPostStorage()

	// Инициализация logger
	eventLogger := logger.NewEventLogger(eventBufferLogger)
	defer eventLogger.GracefullShutdown(gracefullTimeout)

	// Инициализация service слоя (с логгером)
	userService := userservice.NewUserService(eventLogger, userStorage)
	postService := postservice.NewPostService(eventLogger, userStorage, postStorage)

	// Инициализация очереди лайков
	likeQueue := queue.NewLikeQueue(eventLogger, postService, likeQueueBuffer)
	defer likeQueue.GracefullShutdown(gracefullTimeout)

	// Инициализация handlers (с очередью)
	userHandler := handlers.NewUserHandler(userService)
	postHandler := handlers.NewPostHandler(postService, likeQueue)

	http.HandleFunc("/users", userHandler.RegisterUser)

	http.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/posts" {
			http.NotFound(w, r)
			return
		}
		if r.Method == http.MethodGet {
			postHandler.GetAllPosts(w, r)
			return
		}
		if r.Method == http.MethodPost {
			postHandler.CreatePost(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/posts/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/like") {
			postHandler.LikePost(w, r)
			return
		}
		postHandler.GetPost(w, r)
	})

	server := &http.Server{
		Addr:    serverAddr,
		Handler: nil,
	}

	// Канал для ошибок сервера
	serverErrors := make(chan error, 1)

	go func() {
		fmt.Printf("Server is running on http://localhost%s\n", serverAddr)
		fmt.Println("Available endpoints:")
		fmt.Println("  POST   /users              - Register user")
		fmt.Println("  POST   /posts              - Create post")
		fmt.Println("  GET    /posts              - Get all posts")
		fmt.Println("  GET    /posts/{id}         - Get post by ID")
		fmt.Println("  POST   /posts/{id}/like    - Like post")
		fmt.Println("\nPress Ctrl+C to stop")

		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("Failed to start server: %v", err)

	case sig := <-shutdown:
		log.Printf("Received signal: %v. Starting graceful shutdown...", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)
			log.Println("Forcing shutdown...")
			server.Close()
		}

		log.Println("Server stopped gracefully")
	}
}
