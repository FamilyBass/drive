package repository

import (
	"context"
	"github.com/familybass/drive/internal/domain/entity"
)

// UserRepository définit l'interface pour la persistance des utilisateurs
type UserRepository interface {
	// Create insère un nouvel utilisateur
	Create(ctx context.Context, user *entity.User) error

	// GetByID récupère un utilisateur par son ID
	GetByID(ctx context.Context, id string) (*entity.User, error)

	// GetByEmail récupère un utilisateur par son email
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	// Update met à jour un utilisateur
	Update(ctx context.Context, user *entity.User) error

	// Exists vérifie si un utilisateur existe
	Exists(ctx context.Context, email string) (bool, error)
}
