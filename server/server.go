package server

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zewolfe/hermes/handlers/cellcube"
	"github.com/zewolfe/hermes/pkg/rapidpro"
)

type Server interface {
	Start() error
	Stop() error

	AddHandler() error
	AddRouter(string, RouteHandler)
	Router() Router
}

type server struct {
	router     chi.Router
	port       string
	httpServer *http.Server
}

type Router = chi.Router

type handler func(http.ResponseWriter, http.Request)
type RouteHandler func(r chi.Router)

func New() Server {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CleanPath)
	r.Use(middleware.Heartbeat("/health"))

	//TODO: change the timeout to be part of config
	r.Use(middleware.Timeout(60 * time.Second))
	baseRouter := chi.NewRouter()

	//TODO: Setup Config
	rapidProUrl := os.Getenv("RAPIDPRO_URL")
	channelId := os.Getenv("RAPIDPRO_CHANNEL_ID")
	token := os.Getenv("RAPIDPRO_TOKEN")

	rpS := rapidpro.New(rapidProUrl, channelId, rapidpro.WithToken(token))

	h := cellcube.Initialise(rpS)
	r.Mount("/u/", baseRouter)
	r.Route("/h/", h)

	s := &server{
		router: r,
	}

	s.httpServer = &http.Server{
		Addr:         "localhost:9090", //TODO: Ah beg clean this up
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 45 * time.Second,
		IdleTimeout:  90 * time.Second,
	}

	return s
}

func (s *server) Start() error {
	err := s.httpServer.ListenAndServe()
	return err
}

func (s *server) Stop() error {
	err := s.httpServer.Shutdown(context.Background())
	return err
}

// TODO: ???
func (s *server) AddHandler() error {
	return nil
}

func (s *server) AddRouter(path string, r RouteHandler) {
	s.router.Route(path, r)
}

func (s *server) Router() chi.Router {
	return s.router
}
