package userservice

import "net/http"

func (us *UserService) RegisterRoutes(router *http.ServeMux) {
	// user routes
	router.HandleFunc("GET /user", us.GetAllUsersHandler)
	router.HandleFunc("GET /user/{id}", us.GetUserHandler)
}
