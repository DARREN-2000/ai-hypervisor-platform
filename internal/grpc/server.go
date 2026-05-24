package grpc

import (
	"context"
	"net"
	"time"
)

// Options configures a future gRPC server.
type Options struct {
	Address         string
	ReadTimeout     time.Duration
	ShutdownTimeout time.Duration
}

// Server is a lightweight scaffold for gRPC service wiring.
type Server struct {
	options Options
}

// NewServer returns a new scaffold server.
func NewServer(options Options) *Server {
	if options.ShutdownTimeout == 0 {
		options.ShutdownTimeout = 10 * time.Second
	}
	return &Server{options: options}
}

// ListenAddr returns the bind address for the server.
func (s *Server) ListenAddr() string {
	return s.options.Address
}

// Start opens a listener and returns it to the caller.
func (s *Server) Start(ctx context.Context) (net.Listener, error) {
	return net.Listen("tcp", s.options.Address)
}

// Shutdown is a placeholder for graceful gRPC shutdown wiring.
func (s *Server) Shutdown(ctx context.Context) error {
	return nil
}
