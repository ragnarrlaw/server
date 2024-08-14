package main

import (
	"fmt"

	"github.com/raganrrlaw/server/server"
)

/*
   Limit this to a single call, running the server
*/

func main() {
  fmt.Println(">>>>Proc_Server Serve>>>>")
  server := server.NewServer(":9999")
  server.Run()
}
