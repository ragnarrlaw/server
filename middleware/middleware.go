package middleware

import "net/http"

type Middlware func (http.Handler) http.Handler

func Use(fs ...Middlware) Middlware {
  return func (next http.Handler) http.Handler {
    for i := len(fs) - 1; i >= 0; i-- {
      f := fs[i]
      next = f(next)
    }
    return next
  }
}
