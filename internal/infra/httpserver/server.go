package httpserver

import (
	"context"
	"net"
	"net/http"
	"time"
)

type Server struct {
	http *http.Server
}

type Config struct {
	Addr string
}

func NewServer(cfg Config, h *Health) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", h.Healthz)
	mux.HandleFunc("/readyz", h.Readyz)

	serv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	return &Server{http: serv}
}

func (s *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", s.http.Addr)
	if err != nil {
		return err
	}
	return s.http.Serve(ln)
}

// Graceful http server shutdown
func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
