package avenue

import (
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/avenue-contrib/avenue/context"
	"github.com/avenue-contrib/avenue/mux"
)

type Server struct {
	mux.Mux
	http.Server
	listener  net.Listener
	signal    chan os.Signal
	OpenConns sync.WaitGroup
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func New() *Server {
	m := mux.New("/")
	return &Server{
		m,
		http.Server{},
		nil,
		make(chan os.Signal, 1),
		sync.WaitGroup{},
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

func (s *Server) Run(addr string) error {
	log.Printf("Listening on: %s\n\n", addr)
	s.Addr = addr
	s.Handler = s
	s.ReadTimeout = time.Second * 2
	s.WriteTimeout = time.Second * 2
	s.ConnState = s.ConnHandler

	go s.watchSignals()

	if s.Addr == "" {
		s.Addr = ":http"
	}
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		panic(err)
	}
	s.listener = tcpKeepAliveListener{ln.(*net.TCPListener)}
	return s.Serve(s.listener)
}

func (s *Server) ConnHandler(conn net.Conn, state http.ConnState) {
	switch state {
	case http.StateNew:
		s.OpenConns.Add(1)
	case http.StateIdle:
		log.Println("Idling")
	case http.StateHijacked, http.StateClosed:
		s.OpenConns.Done()
	}
}

func (s *Server) GracefulShutdown(restart bool) {
	// shutdown keep-alives
	s.SetKeepAlivesEnabled(false)
	s.OpenConns.Wait()
	s.listener.Close()

	if restart {
		log.Println("Restarting....")
		// start another process
		// for now, program will terminate unless other locks held
	} else {
		log.Println("Terminating...")
	}
}

func (s *Server) watchSignals() {
	signal.Notify(s.signal, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGUSR1)
	for {
		sig := <-s.signal
		switch sig {
		case syscall.SIGINT, syscall.SIGKILL:
			log.Printf("Halting due to: %s\n", sig)
			s.listener.Close()
		case syscall.SIGTERM, syscall.SIGUSR1:
			log.Println("Starting graceful shutdown")
			s.GracefulShutdown(sig == syscall.SIGUSR1)
		}
	}
}
