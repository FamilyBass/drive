package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/familybass/drive/internal/domain/service"
	"github.com/go-chi/chi/v5"
)

// DriveHandler gère les endpoints de gestion des fichiers
type DriveHandler struct {
	driveService *service.DriveService
}

// NewDriveHandler crée un nouveau gestionnaire de fichiers
func NewDriveHandler(driveService *service.DriveService) *DriveHandler {
	return &DriveHandler{
		driveService: driveService,
	}
}

// Upload télécharge un fichier
func (h *DriveHandler) Upload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user_id").(string)

	// Limiter à 200MB
	r.Body = http.MaxBytesReader(w, r.Body, 200<<20)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "file too large", http.StatusBadRequest)
		return
	}

	file, fh, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "missing file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	uploadedFile, err := h.driveService.UploadFile(ctx, userID, fh.Filename, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":   uploadedFile.ID,
		"path": uploadedFile.Path,
	})
}

// List liste les fichiers accessibles
func (h *DriveHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user_id").(string)
	isAdmin := ctx.Value("is_admin").(bool)

	files, err := h.driveService.ListFiles(ctx, userID, isAdmin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshaller les résultats
	type FileDTO struct {
		ID        string `json:"ID"`
		OwnerID   string `json:"OwnerID"`
		Filename  string `json:"Filename"`
		Size      int64  `json:"Size"`
		CreatedAt int64  `json:"CreatedAt"`
	}

	var dtos []FileDTO
	for _, f := range files {
		dtos = append(dtos, FileDTO{
			ID:        f.ID,
			OwnerID:   f.OwnerID,
			Filename:  f.Filename,
			Size:      f.Size,
			CreatedAt: f.CreatedAt.Unix(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dtos)
}

// Download télécharge un fichier
func (h *DriveHandler) Download(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user_id").(string)
	isAdmin := ctx.Value("is_admin").(bool)
	fileID := chi.URLParam(r, "id")

	file, err := h.driveService.DownloadFile(ctx, fileID, userID, isAdmin)
	if err != nil {
		if err.Error() == "forbidden" {
			http.Error(w, "forbidden", http.StatusForbidden)
		} else if err.Error() == "file not found" {
			http.Error(w, "not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Ouvrir et servir le fichier
	f, err := http.Dir(filepath.Dir(file.Path)).Open(filepath.Base(file.Path))
	if err != nil {
		http.Error(w, "file error", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Length", strconv.FormatInt(file.Size, 10))
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.Filename))
	http.ServeContent(w, r, file.Filename, file.CreatedAt, f)
}
