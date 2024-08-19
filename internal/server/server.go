package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/raganrrlaw/server/internal/config"
	database "github.com/raganrrlaw/server/internal/db"
	"github.com/raganrrlaw/server/internal/middleware"
	publicservice "github.com/raganrrlaw/server/internal/services/public_service"
	userservice "github.com/raganrrlaw/server/internal/services/user_service"
	userrepository "github.com/raganrrlaw/server/internal/services/user_service/user_repository"
)

type Server struct {
	address string
	storage *database.Storage
	server  *http.Server
}

func NewServer(addr string) *Server {
	err := config.Init()
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(config.EXIT_FAILURE)
	}

	cfg, err := database.NewStorageConfig()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	store, err := database.NewStorage(*cfg)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return &Server{
		address: addr,
		storage: store,
	}
}

func (s *Server) Run() {
	// initialize the router/nux here
	router := http.NewServeMux()

	// initialize the repositories here
	userRepo := userrepository.NewUserRepo(s.storage)

	// initialize the user services
	userService := userservice.NewUserService(userRepo)
	userService.RegisterRoutes(router)

	// initialize the public services here
	publicservice.RegisterRoutes(router)

	// initialize the services here
	routerV1 := http.NewServeMux()
	routerV1.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	// initialize the middleware

	s.server = &http.Server{
		Addr: s.address,
		Handler: middleware.CorsMiddleware(
			middleware.LogRequestDetailsMiddleware(
				router,
			),
		), // wrap the router with the middleware here use(router)
	}

	// listen for interrupt or terminate signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf(">>>> Server is listening on: %s\n", s.address)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf(">>>> Server failed: %s\n", err)
		}
	}()

	<-stop

	if err := s.Shutdown(context.Background()); err != nil {
		log.Fatalf(">>>> Server Shutdown Failed:%+v", err)
	}
	log.Println(">>>> Server exited properly")
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println(">>>> Shutting down server...")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	if err := s.storage.Shutdown(); err != nil {
		return err
	}
	return nil
}
