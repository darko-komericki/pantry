package api

import (
	"context"
	"log/slog"
)

// Pinger is the one thing the handler needs from the database.
// *pgxpool.Pool satisfies it without declaring so; tests use a fake.
type Pinger interface {
	Ping(ctx context.Context) error
}

// Server implements StrictServerInterface. Dependencies are plain fields
// set by the constructor; there is no global state.
type Server struct {
	db     Pinger
	logger *slog.Logger
}

// Compile-time check that Server covers every operation in the contract.
// If the spec gains an endpoint and `make gen` is run, the build fails here
// until the handler implements it.
var _ StrictServerInterface = (*Server)(nil)

func NewServer(db Pinger, logger *slog.Logger) *Server {
	return &Server{db: db, logger: logger}
}

func (s *Server) GetHealth(ctx context.Context, _ GetHealthRequestObject) (GetHealthResponseObject, error) {
	if err := s.db.Ping(ctx); err != nil {
		// The client only needs to know the DB is down; the cause goes to the log.
		s.logger.ErrorContext(ctx, "health check: ping db", "error", err)
		return GetHealth503JSONResponse{Status: Unavailable}, nil
	}
	return GetHealth200JSONResponse{Status: Ok}, nil
}
