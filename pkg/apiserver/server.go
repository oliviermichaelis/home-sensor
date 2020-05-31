package apiserver

import "net/http"

type Server struct {
	debug        bool
	repositories []Repository
	mux          *http.ServeMux
	addr         string
}

func (s *Server) Start() error {
	// TODO wait for databases
	// TODO healthcheck /health
	panic("not implemented")
}

func NewServer(debug bool, mux *http.ServeMux, addr string, repositories ...Repository) Server {
	return Server{
		debug:        debug,
		mux:          mux,
		addr:         addr,
		repositories: repositories,
	}
}
