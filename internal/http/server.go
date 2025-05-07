package http

import (
	"log/slog"
	"net"
	"net/http"
	"strconv"
)

// Based configuration pattern
type Config struct {
	BindAddress      string
	Port             int
	StoragePath      string
	MaxFileSizeMB    int64
	AllowedFileTypes []string
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

	// Routes for the future ;P
	// mux.Handle("/api/upload", http.HandleFunc())
	// mux.Handle("/api/files/", http.HandleFunc())
	// mux.Handle("/api/files", http.HandleFunc())
	// mux.Handle("/f/", http.HandleFunc())

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
