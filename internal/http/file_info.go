package http

import (
	"net/http"

	"github.com/Mylamuu/drift/internal/storage"
)

func (s *Server) handleGetFileInfo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	fileMetadata, ok := storage.Get(id)
	if !ok {
		writeJSONError(w, http.StatusBadRequest, "Invalid file ID")
		return
	}

	writeJSON(w, http.StatusOK, fileMetadata)
}
