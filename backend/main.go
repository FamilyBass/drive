package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	adapthttp "github.com/familybass/drive/internal/adapter/http"
	"github.com/familybass/drive/internal/bootstrap"
)

func main() {
	// Configuration
	dataDir := os.Getenv("DATA_DIR")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Bootstrap: Initialiser la base de données
	db, err := bootstrap.Database(dataDir)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	// Bootstrap: Créer les repositories
	userRepo, fileRepo := bootstrap.Repositories(db)

	// Bootstrap: Créer les services
	authService, driveService, err := bootstrap.Services(userRepo, fileRepo)
	if err != nil {
		log.Fatalf("failed to initialize services: %v", err)
	}

	// Handlers
	authHandler := adapthttp.NewAuthHandler(authService)
	driveHandler := adapthttp.NewDriveHandler(driveService)

	// Router
	r := chi.NewRouter()

	// Middlewares globaux
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	// Routes publiques
	r.Post("/api/register", authHandler.Register)
	r.Post("/api/login", authHandler.Login)

	// Routes protégées
	r.Group(func(r chi.Router) {
		r.Use(adapthttp.AuthMiddleware(authService))

		// Drive
		r.Post("/api/drive/upload", driveHandler.Upload)
		r.Get("/api/drive/list", driveHandler.List)
		r.Get("/api/drive/download/{id}", driveHandler.Download)

		// Admin
		r.Group(func(r chi.Router) {
			r.Use(adapthttp.AdminMiddleware)
			r.Post("/api/admin/validate", authHandler.ValidateUser)
		})
	})

	// Démarrer le serveur
	addr := ":" + port
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("listen error: %v", err)
	}
}
