package server

import (
	"atlas/config"
	"atlas/services"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type Server struct {
	router   *echo.Echo
	config   config.Config
	services *services.Service

	authMiddleware echo.MiddlewareFunc
}

// Validator wraps go-playground/validator and implements echo.Validator interface
type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i any) error {
	return v.validator.Struct(i)
}

func NewServer(config config.Config) (*Server, error) {

	services := services.NewService(config)
	server := &Server{

		config:   config,
		services: services,
	}
	server.setupRouter()
	return server, nil
}
func (s *Server) setupRouter() {
	router := echo.New()

	// Validator
	router.Validator = &Validator{validator: validator.New()}

	router.Use(echoMiddleware.Recover())
	router.Use(echoMiddleware.CORS())
	router.Use(echoMiddleware.RequestLogger())

	// Health check
	router.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	router.POST("/utterance", s.HandleSpeech)

	s.router = router
}

func (s *Server) Start(address string) error {
	return s.router.Start(address)
}
