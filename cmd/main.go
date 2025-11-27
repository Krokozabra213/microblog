package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"microblog/internal/handlers"
	"microblog/internal/service"
	"microblog/internal/storage"
)

func main() {
	// Инициализация storage слоя
	userStorage := storage.NewUserStorage()
	postStorage := storage.NewPostStorage()

	// Инициализация service слоя
	userService := service.NewUserService(userStorage)
	postService := service.NewPostService(postStorage, userStorage)

	// Инициализация handlers
	userHandler := handlers.NewUserHandler(userService)
	postHandler := handlers.NewPostHandler(postService)

	// Настройка роутинга
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

	// Создаём сервер
	port := ":8080"
	server := &http.Server{
		Addr:    port,
		Handler: nil,
	}

	// Канал для ошибок сервера
	serverErrors := make(chan error, 1)

	// Запускаем сервер в горутине
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

	// Слушаем сигналы остановки
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Ждём либо ошибки, либо сигнала остановки
	select {
	case err := <-serverErrors:
		log.Fatalf("Failed to start server: %v", err)

	case sig := <-shutdown:
		log.Printf("Received signal: %v. Starting graceful shutdown...", sig)

		// Даём 30 секунд на завершение текущих запросов
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Graceful shutdown
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)
			log.Println("Forcing shutdown...")
			server.Close()
		}

		log.Println("Server stopped gracefully")
	}
}
