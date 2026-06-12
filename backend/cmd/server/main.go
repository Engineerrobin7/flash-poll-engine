package main

import (
	"flashpoll/internal/db"
	"flashpoll/internal/http/handlers"
	"flashpoll/internal/middleware"
	"flashpoll/internal/repository"
	"flashpoll/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"log"
	"net/http"
	"os"
)

func main() {
	dbPath := os.Getenv("DATABASE_URL")
	if dbPath == "" {
		dbPath = "flashpoll.db"
	}

	database, err := db.InitDB(dbPath)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	if err := db.Migrate(database); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	repo := repository.NewPollRepository(database)
	broker := service.NewBroker()
	pollService := service.NewPollService(repo, broker)
	pollHandler := handlers.NewPollHandler(pollService)

	r := chi.NewRouter()

	r.Use(middleware.Logging)
	r.Use(chi.Middleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("CRITICAL RECOVERY: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000", "*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/api/events", broker.ServeHTTP)

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.RateLimit)
		r.Get("/polls", pollHandler.GetPolls)
		r.Post("/polls", pollHandler.CreatePoll)
		r.Patch("/polls/{id}/vote", pollHandler.Vote)
		r.Delete("/polls/{id}", pollHandler.DeletePoll)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
