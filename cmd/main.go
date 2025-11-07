package main

import (
	"fmt"
	"log"
	"microblog/internal/handlers"
	"microblog/internal/service"
	"microblog/internal/storage"
	"net/http"
	"strings"
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

	// Роутинг для /posts (точное совпадение)
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

	// Роутинг для /posts/{id} и /posts/{id}/like
	http.HandleFunc("/posts/", func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, заканчивается ли на /like
		if strings.HasSuffix(r.URL.Path, "/like") {
			postHandler.LikePost(w, r)
			return
		}
		// Иначе это GetPost
		postHandler.GetPost(w, r)
	})

	// Запуск сервера
	port := ":8080"
	fmt.Printf("Server is running on http://localhost%s\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  POST   /users              - Register user")
	fmt.Println("  POST   /posts              - Create post")
	fmt.Println("  GET    /posts              - Get all posts")
	fmt.Println("  GET    /posts/{id}         - Get post by ID")
	fmt.Println("  POST   /posts/{id}/like    - Like post")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
