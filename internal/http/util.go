package http

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"net/http"
)

// Might be easier to just return a fallback value here instead of an ok value
func randomId(n int) (string, bool) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		slog.Error("Failed to generate random bytes for ID", "error", err)
		return "", false
	}

	return hex.EncodeToString(b), true
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			slog.Error("Couldn't write data as response", "error", err)
		}
	}
}

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
