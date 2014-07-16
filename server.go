package avenue

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"avenue/context"
	"avenue/mux"
)

type Server struct {
	mux.Mux
}

func New() *Server {
	m := mux.New("/")
	return &Server{
		m,
	}
}

func (s *Server) Use(m context.Handler) {
	// TODO: checks
	s.Mux.Handlers = append(s.Mux.Handlers, m)
}

func (s *Server) Templates(path, delims string, extensions []string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		panic(errors.New("Template root does not exist: " + path))
	} else if err != nil {
		// TODO: don't panic
		panic(err)
	} else {
		context.TemplatePath = path
	}

	if split := strings.Split(delims, " "); len(split) != 2 {
		panic(errors.New("Failed to parse delims: " + delims))
	} else {
		context.TemplateDelims = delims
	}

	if len(extensions) == 0 {
		context.TemplateExtensions = []string{"html"}
	} else {
		context.TemplateExtensions = extensions
	}
}

func (s *Server) ServeFiles(path string, root http.FileSystem) {
	s.Router().ServeFiles(path, root)
}

// ServeHTTP makes the router implement the http.Handler interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.Router().ServeHTTP(context.Wrap(w), req)
}

func (s *Server) Run(addr string) {
	log.Printf("Listening on: %s\n\n", addr)

	http.ListenAndServe(addr, s)
}
