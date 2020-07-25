package main

import (
	"github.com/KairoBoni/boltons/pkg/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

type Server struct {
	route    *echo.Echo
	handlers *Handler
}

//NewServer create a new server of the REST-API
func NewServer(dbConfigFilepath string) (*Server, error) {
	store, err := database.NewStore(dbConfigFilepath)
	if err != nil {
		return nil, err
	}
	s := &Server{
		route: echo.New(),
		handlers: &Handler{
			db: store,
		},
	}
	s.setupRoutes()

	return s, nil

}

func (s *Server) setupRoutes() {
	s.route.Use(middleware.Recover())

	s.route.GET("/nfe/amount/:accessKey", s.handlers.getNfeAmount)
}

//Run starts the server
func (s *Server) Run() error {
	s.route.HideBanner = true
	log.Info().Msg("Server start on localhost:5002/")
	return s.route.Start(":5002")
}
