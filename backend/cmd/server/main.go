package main

import (
	"context"
	"flashpoll/internal/db"
	"flashpoll/internal/http/handlers"
	"flashpoll/internal/middleware"
	"flashpoll/internal/repository"
	"flashpoll/internal/service"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	// Standard Production Middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.Logging)
	r.Use(chimiddleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000", "*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Get("/api/events", broker.ServeHTTP)
	r.Get("/api/stats", pollHandler.GetStats)

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

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Graceful Shutdown Logic
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-done
	log.Print("Server Stopping...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Cleanly")
}
