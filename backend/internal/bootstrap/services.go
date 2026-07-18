package bootstrap

import (
	"context"
	"log"
	"os"

	domrepo "github.com/familybass/drive/internal/domain/repository"
	"github.com/familybass/drive/internal/domain/service"
)

// Services crée et injecte les services
func Services(userRepo domrepo.UserRepository, fileRepo domrepo.FileRepository) (*service.AuthService, *service.DriveService, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev_secret_change_me"
	}

	authService := service.NewAuthService(userRepo, jwtSecret)
	driveService := service.NewDriveService(fileRepo, GetMediaDir(os.Getenv("DATA_DIR")))

	// Créer l'admin initial si les variables d'environnement sont définies
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPass := os.Getenv("ADMIN_PASSWORD")
	log.Printf("admin email: %s, admin password: %s", adminEmail, adminPass)
	if adminEmail != "" && adminPass != "" {
		ctx := context.Background()
		if err := authService.EnsureAdmin(ctx, adminEmail, adminPass); err != nil {
			log.Printf("warning: failed to ensure admin: %v", err)
		}
	}

	return authService, driveService, nil
}
