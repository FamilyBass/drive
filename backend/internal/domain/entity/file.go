package entity

import "time"

// File représente un fichier stocké
type File struct {
	ID        string
	OwnerID   string
	Filename  string
	Size      int64
	Path      string
	CreatedAt time.Time
}

// NewFile crée un nouveau fichier
func NewFile(id, ownerID, filename, path string, size int64, createdAt time.Time) *File {
	return &File{
		ID:        id,
		OwnerID:   ownerID,
		Filename:  filename,
		Size:      size,
		Path:      path,
		CreatedAt: createdAt,
	}
}
