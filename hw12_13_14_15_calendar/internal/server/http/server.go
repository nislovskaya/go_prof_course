package internalhttp

import (
	"context"
)

type Server struct { // TODO
}

type Logger interface { // TODO
}

type Application interface { // TODO
}

func NewServer(logger Logger, app Application) *Server { //revive:disable
	return &Server{}
}

func (s *Server) Start(ctx context.Context) error {
	// TODO
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error { //revive:disable
	// TODO
	return nil
}

// TODO
