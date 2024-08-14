package middleware

import (
  "net/http"
  "log"
)

func LogginMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
    log.Printf("Client: %s\t, Request Route: %s\t\n\n", r.Header,  r.URL.Path)
    next.ServeHTTP(w, r)
  })
}
