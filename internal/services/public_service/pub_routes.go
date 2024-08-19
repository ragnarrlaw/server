package publicservice

import (
	"net/http"
)

func RegisterRoutes(router *http.ServeMux) {
  router.HandleFunc("GET /", IndexHandler)
  router.HandleFunc("GET /favicon.ico", FaviconHandler)
}
