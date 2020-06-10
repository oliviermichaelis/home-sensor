package apiserver

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct {
	ClimateEndpoint
	mux          *http.ServeMux
	addr         string
	events       chan error		// send true if step succeeded, false if step failed
	healthy      bool
	debug        bool
}

func (s *Server) Start() error {
	go s.health()
	go s.connectRepositories()

	return s.listen()
}

func (s *Server) listen() error {
	s.mux.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		status := http.StatusNotFound
		if s.healthy {
			status = http.StatusOK
		}
		writer.WriteHeader(status)
	})

	s.mux.HandleFunc("/measurements/climate", s.climateHandler)

	// Signal that application is listening on http
	s.events <- nil
	return http.ListenAndServe(s.addr, s.mux)
}

// Change healthy variable if all functionalities reported to by health. Repositories + http endpoint
func (s *Server) health() {
	defer close(s.events)
	for i := 0; i < len(s.Repositories)+1; i++ {
	if e := <-s.events; e != nil {
			log.Fatalf("health: received fatal error: %v", e)
		}
	}

	s.healthy = true
}

func (s *Server) connectRepositories() {
	timeout := time.Minute * 1
	deadline := time.Now().Add(timeout)

	for _, r := range s.Repositories {
		retries := 0
		var repoErr RepoError
		success := false
		for ; time.Now().Before(deadline); retries++ {
			if repoErr = r.attemptConnect(); repoErr.Unwrap() == nil {
				success = true
				log.Printf("connectRepositories: successfully connected to %s", repoErr.Url)
				break
			}

			log.Printf("connectRepository: failed to connect to %v: %v", repoErr.Url, repoErr.Err)
			time.Sleep(time.Second * 5)
		}

		if success {
			s.events <- nil
		} else {
			s.events <- fmt.Errorf("connectRepository: failed after %d retries: %v", retries, repoErr.Err)
		}
	}
}

func NewServer(debug bool, mux *http.ServeMux, addr string, repositories ...Repository) Server {
	return Server{
		debug:        debug,
		mux:          mux,
		addr:         addr,
		ClimateEndpoint: ClimateEndpoint{Repositories: repositories},
		events:       make(chan error),
		healthy:      false,
	}
}
