package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Mylamuu/drift/internal/storage"
)

func (s *Server) handleDownloadFile(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	fileMetadata, ok := storage.Get(id)
	if !ok || time.Now().After(fileMetadata.ExpiresAt) {
		writeJSONError(w, http.StatusBadRequest, "Invalid file ID")
		return
	}

	file, err := os.Open(fileMetadata.StoragePath)
	if err != nil {
		slog.Error("Unable to open file from disk", "path", fileMetadata.StoragePath, "error", err)
		writeJSONError(w, http.StatusInternalServerError, "Unable to read file from disk")
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", fileMetadata.ContentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileMetadata.OriginalName))

	http.ServeContent(w, r, fileMetadata.OriginalName, fileMetadata.UploadedAt, file)
}
