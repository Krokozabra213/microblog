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
	"microblog/internal/service"
	"microblog/internal/storage"
)

func main() {
	// Инициализация storage слоя
	userStorage := storage.NewUserStorage()
	postStorage := storage.NewPostStorage()

	// Инициализация logger
	eventLogger := logger.NewEventLogger()

	// Инициализация service слоя (с логгером)
	userService := service.NewUserService(userStorage, eventLogger)
	postService := service.NewPostService(postStorage, userStorage, eventLogger)

	// Инициализация очереди лайков
	likeQueue := queue.NewLikeQueue(postService, eventLogger)

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

	port := ":8080"
	server := &http.Server{
		Addr:    port,
		Handler: nil,
	}

	// Канал для ошибок сервера
	serverErrors := make(chan error, 1)

	go func() {
		fmt.Printf("Server is running on http://localhost%s\n", port)
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
