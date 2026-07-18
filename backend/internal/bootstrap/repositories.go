package bootstrap

import (
	"database/sql"

	adaptrepo "github.com/familybass/drive/internal/adapter/repository"
	domrepo "github.com/familybass/drive/internal/domain/repository"
)

// Repositories crée les repositories
func Repositories(db *sql.DB) (domrepo.UserRepository, domrepo.FileRepository) {
	userRepo := adaptrepo.NewSQLiteUserRepository(db)
	fileRepo := adaptrepo.NewSQLiteFileRepository(db)
	return userRepo, fileRepo
}
