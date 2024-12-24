package main

import (
	"log"
	"os"

	"github.com/ragnarrlaw/server/config"
	"github.com/ragnarrlaw/server/internal/server"
)

/*
   Limit this to a single call, running the server
*/

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(config.EXIT_FAILURE)
	} else {
		log.Println(">>>> JSON Server Serve")
		server := server.NewServer(":9999")
		server.Run()
	}
}
