package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/familybass/drive/internal/domain/entity"
	"github.com/familybass/drive/internal/domain/repository"
	"github.com/google/uuid"
)

// DriveService gère les opérations de fichiers
type DriveService struct {
	fileRepo repository.FileRepository
	mediaDir string
}

// NewDriveService crée un nouveau service de gestion des fichiers
func NewDriveService(fileRepo repository.FileRepository, mediaDir string) *DriveService {
	return &DriveService{
		fileRepo: fileRepo,
		mediaDir: mediaDir,
	}
}

// UploadFile télécharge un fichier
func (s *DriveService) UploadFile(ctx context.Context, userID, filename string, content io.Reader) (*entity.File, error) {
	// Générer l'ID et le chemin
	fileID := uuid.New().String()
	filepath := filepath.Join(s.mediaDir, fileID+"_"+filepath.Base(filename))

	// Créer le fichier
	out, err := os.Create(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Copier le contenu (avec limite pour la RAM)
	limitedReader := io.LimitReader(content, 200<<20) // 200MB max
	n, err := io.Copy(out, limitedReader)
	if err != nil {
		os.Remove(filepath)
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	// Créer l'entité
	file := entity.NewFile(fileID, userID, filename, filepath, n, time.Now())

	// Persister
	if err := s.fileRepo.Create(ctx, file); err != nil {
		os.Remove(filepath)
		return nil, fmt.Errorf("failed to save file metadata: %w", err)
	}

	return file, nil
}

// DownloadFile récupère un fichier
func (s *DriveService) DownloadFile(ctx context.Context, fileID, userID string, isAdmin bool) (*entity.File, error) {
	file, err := s.fileRepo.GetByID(ctx, fileID)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	// Vérifier l'accès
	if file.OwnerID != userID && !isAdmin {
		return nil, fmt.Errorf("forbidden")
	}

	// Vérifier que le fichier existe physiquement
	if _, err := os.Stat(file.Path); err != nil {
		return nil, fmt.Errorf("file not available: %w", err)
	}

	return file, nil
}

// ListFiles liste les fichiers accessibles
func (s *DriveService) ListFiles(ctx context.Context, userID string, isAdmin bool) ([]*entity.File, error) {
	var files []*entity.File
	var err error

	if isAdmin {
		files, err = s.fileRepo.ListAll(ctx, 100)
	} else {
		files, err = s.fileRepo.ListByOwner(ctx, userID, 100)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	return files, nil
}
