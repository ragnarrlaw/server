package server

import (
	"log"
	"net/http"

	"github.com/raganrrlaw/server/database"
	"github.com/raganrrlaw/server/middleware"
	"github.com/raganrrlaw/server/services"
)

/*
  Create the server and export it.
  Use the imports and other function to make the router.
*/

type Server struct {
	address string
	dbPool *database.DBPool 
}

func NewServer(addr string) *Server {

  cfg, err := database.NewDatabaseConfig()
  if err != nil {
    log.Fatal(err)
    return nil
  }

  store, err := database.NewDBStorage(*cfg)
  if err != nil {
    log.Fatal(err)
    return nil
  }

	return &Server{
		address: addr,
		dbPool:  store,
	}
}

func (s *Server) Run() {
  router := http.NewServeMux()

  userService := services.NewUserService(s.dbPool)
  userService.RegisterRoutes(router)

  routerV1 := http.NewServeMux()
  routerV1.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

  use := middleware.Use(middleware.LogginMiddleware,)

	server := http.Server{
		Addr:    s.address,
		Handler: use(router),
	}

	log.Printf("Server is listening on: %s\n", s.address)
	log.Fatalf("Server failed: %s\n", server.ListenAndServe().Error())
}

