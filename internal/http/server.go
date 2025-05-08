package http

import (
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"
)

// Based configuration pattern
type Config struct {
	BindAddress      string
	Port             int
	StoragePath      string
	MaxFileSize      int64
	AllowedFileTypes []string
	KeepTime         time.Duration
}

type Server struct {
	config     Config
	router     *http.ServeMux
	httpServer *http.Server
}

func NewServer(config Config) *Server {
	mux := http.NewServeMux()
	server := &Server{
		config: config,
		router: mux,
	}

	mux.Handle("POST /api/upload", http.HandlerFunc(server.handleUploadFile))
	mux.Handle("GET /api/files", http.HandlerFunc(server.handleListFiles))
	mux.Handle("GET /api/files/{id}", http.HandlerFunc(server.handleGetFileInfo))

	server.httpServer = &http.Server{
		Addr:    net.JoinHostPort(config.BindAddress, strconv.Itoa(config.Port)),
		Handler: mux,
	}

	return server
}

func (s *Server) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) Start() error {
	slog.Info("Starting HTTP server", "address", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}
