package http

import (
	"net/http"

	"github.com/Mylamuu/drift/internal/storage"
)

func (s *Server) handleListFiles(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, storage.List())
}
