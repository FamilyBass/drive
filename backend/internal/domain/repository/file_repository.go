package repository

import (
	"context"
	"github.com/familybass/drive/internal/domain/entity"
)

// FileRepository définit l'interface pour la persistance des fichiers
type FileRepository interface {
	// Create insère un nouveau fichier
	Create(ctx context.Context, file *entity.File) error

	// GetByID récupère un fichier par son ID
	GetByID(ctx context.Context, id string) (*entity.File, error)

	// ListByOwner liste les fichiers d'un propriétaire
	ListByOwner(ctx context.Context, ownerID string, limit int) ([]*entity.File, error)

	// ListAll liste tous les fichiers (admin)
	ListAll(ctx context.Context, limit int) ([]*entity.File, error)
}
