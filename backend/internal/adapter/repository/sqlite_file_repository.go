package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/familybass/drive/internal/domain/entity"
)

// SQLiteFileRepository implémente FileRepository avec SQLite
type SQLiteFileRepository struct {
	db *sql.DB
}

// NewSQLiteFileRepository crée un nouveau repository SQLite pour les fichiers
func NewSQLiteFileRepository(db *sql.DB) *SQLiteFileRepository {
	return &SQLiteFileRepository{db: db}
}

// Create insère un nouveau fichier
func (r *SQLiteFileRepository) Create(ctx context.Context, file *entity.File) error {
	query := `INSERT INTO files(id,owner_id,filename,size,path,created_at) 
	          VALUES(?,?,?,?,?,?)`
	_, err := r.db.ExecContext(ctx, query,
		file.ID, file.OwnerID, file.Filename,
		file.Size, file.Path, file.CreatedAt.Unix(),
	)
	if err != nil {
		return fmt.Errorf("failed to insert file: %w", err)
	}
	return nil
}

// GetByID récupère un fichier par son ID
func (r *SQLiteFileRepository) GetByID(ctx context.Context, id string) (*entity.File, error) {
	file := &entity.File{}

	query := `SELECT id,owner_id,filename,size,path,created_at FROM files WHERE id=?`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&file.ID, &file.OwnerID, &file.Filename,
		&file.Size, &file.Path, (*unixTime)(&file.CreatedAt),
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("file not found")
		}
		return nil, fmt.Errorf("failed to query file: %w", err)
	}

	return file, nil
}

// ListByOwner liste les fichiers d'un propriétaire
func (r *SQLiteFileRepository) ListByOwner(ctx context.Context, ownerID string, limit int) ([]*entity.File, error) {
	query := `SELECT id,owner_id,filename,size,path,created_at FROM files 
	          WHERE owner_id=? ORDER BY created_at DESC LIMIT ?`
	rows, err := r.db.QueryContext(ctx, query, ownerID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query files: %w", err)
	}
	defer rows.Close()

	var files []*entity.File
	for rows.Next() {
		file := &entity.File{}
		if err := rows.Scan(
			&file.ID, &file.OwnerID, &file.Filename,
			&file.Size, &file.Path, (*unixTime)(&file.CreatedAt),
		); err != nil {
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}
		files = append(files, file)
	}

	return files, rows.Err()
}

// ListAll liste tous les fichiers
func (r *SQLiteFileRepository) ListAll(ctx context.Context, limit int) ([]*entity.File, error) {
	query := `SELECT id,owner_id,filename,size,path,created_at FROM files 
	          ORDER BY created_at DESC LIMIT ?`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query files: %w", err)
	}
	defer rows.Close()

	var files []*entity.File
	for rows.Next() {
		file := &entity.File{}
		if err := rows.Scan(
			&file.ID, &file.OwnerID, &file.Filename,
			&file.Size, &file.Path, (*unixTime)(&file.CreatedAt),
		); err != nil {
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}
		files = append(files, file)
	}

	return files, rows.Err()
}
