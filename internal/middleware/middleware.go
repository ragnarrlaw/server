package middleware

import (
	"log"
	"net/http"

	"github.com/fatih/color"
)

type Middleware func(http.Handler) http.Handler

func AppendMiddlewareStack(fs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(fs) - 1; i >= 0; i-- {
			next = fs[i](next)
		}
		return next
	}
}

func LogRequestDetailsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		method := r.Method
		endpoint := r.URL.Path

		var coloredMethod string
		switch method {
		case http.MethodGet:
			coloredMethod = color.New(color.BgGreen, color.FgBlack).Sprint(method)
		case http.MethodPost:
			coloredMethod = color.New(color.BgBlue, color.FgWhite).Sprint(method)
		case http.MethodPut:
			coloredMethod = color.New(color.BgYellow, color.FgBlack).Sprint(method)
		case http.MethodDelete:
			coloredMethod = color.New(color.BgRed, color.FgWhite).Sprint(method)
		default:
			coloredMethod = method
		}

		log.Printf("%s\t%s\t%s\n", clientIP, coloredMethod, endpoint)
		next.ServeHTTP(w, r)
	})
}
