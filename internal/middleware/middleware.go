package middleware

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/fatih/color"
	"github.com/ragnarrlaw/server/internal/types/search"
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
		case http.MethodPatch:
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

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight (OPTIONS) requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Invalid content type. Only application/json is allowed", http.StatusUnsupportedMediaType)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func QueryParameterParsingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			paginateParams := search.Paginate{}

			limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
			if err != nil {
				paginateParams.Limit = 10
			} else {
				paginateParams.Limit = limit
			}

			offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
			if err != nil {
				paginateParams.Offset = 0
			} else {
				paginateParams.Offset = offset
			}
			context := context.WithValue(r.Context(), search.SearchKey, paginateParams)
			r = r.WithContext(context)
		}
		next.ServeHTTP(w, r)
	})
}
