package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Server struct {
	Router   *mux.Router
	Basepath string
	Port     string
}

func NewServer(basepath, port string) (*Server, error) {
	server := &Server{
		Basepath: basepath,
		Router:   mux.NewRouter(),
		Port:     port,
	}

	return server, nil
}

// AddRoute adds a new route to the server
func (s *Server) AddRoute(pattern string, handler http.HandlerFunc) {
	s.Router.HandleFunc(pattern, handler)
}

// Start initiates the server to start listening on the specified port
func (s *Server) Start() error {
	fmt.Printf("Server starting on port %s\n", s.Port)

	// Setup CORS
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	// TODO: Adjust this to allow specific origins when moving to production
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// Apply the CORS middleware to our router
	return http.ListenAndServe(":"+s.Port, handlers.CORS(originsOk, headersOk, methodsOk)(s.Router))
}
