package http

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Mylamuu/drift/internal/storage"
)

func (s *Server) handleUploadFile(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(s.config.MaxFileSize); err != nil {
		slog.Warn("Unable to parse mutilpart form data", "error", err)
		writeJSONError(w, http.StatusBadRequest, "Unable to parse multipart form data")
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		slog.Error("Unable to get form file", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to read uploaded file")
	}
	defer file.Close()

	if fileHeader.Size > s.config.MaxFileSize {
		writeJSONError(w, http.StatusBadRequest, "File is too large")
		return
	}

	// DetectContentType can use up to 512 bytes so we might as well give it all of it
	buf := make([]byte, 512)
	n, _ := io.ReadFull(file, buf)
	contentType := http.DetectContentType(buf[:n])

	if len(s.config.AllowedFileTypes) > 0 {
		allowed := false

		for _, typ := range s.config.AllowedFileTypes {
			if strings.HasPrefix(contentType, typ) {
				allowed = true
				break
			}
		}

		if !allowed {
			writeJSONError(w, http.StatusBadRequest, "Invalid file type")
			return
		}
	}

	reader := io.MultiReader(bytes.NewReader(buf[:n]), file)
	id, ok := randomId(12)

	if !ok {
		slog.Error("Unable to generate random ID for file")
		writeJSONError(w, http.StatusInternalServerError, "Unable to generate unique ID")
		return
	}

	dst, err := os.Create(storage.Path(id))
	if err != nil {
		slog.Error("Failed to create upload on disk", "path", storage.Path(id), "error", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to create upload on disk")
		return
	}
	defer dst.Close()

	written, err := io.Copy(dst, reader)
	if err != nil {
		slog.Error("Unable to write file to disk", "path", storage.Path(id), "error", err)
		writeJSONError(w, http.StatusInternalServerError, "Unable to write file to disk")
		return
	}

	fileMetadata := storage.File{
		ID:           id,
		OriginalName: fileHeader.Filename,
		Size:         written,
		ContentType:  contentType,
		UploadedAt:   time.Now(),
		ExpiresAt:    time.Now().Add(s.config.KeepTime),
		AccessCount:  0,
		StoragePath:  storage.Path(id),
	}

	storage.Save(fileMetadata)

	writeJSON(w, http.StatusCreated, map[string]storage.File{
		"file": fileMetadata,
	})
}
