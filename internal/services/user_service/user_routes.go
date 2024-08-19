package userservice

import "net/http"

func (us *UserService) RegisterRoutes(router *http.ServeMux) {
  router.HandleFunc("GET /user", us.GetAllUsersHandler)
  router.HandleFunc("GET /user/{id}", us.GetUserHandler)
  router.HandleFunc("POST /user", us.AddUserHandler)
  router.HandleFunc("PATCH /user/{id}", us.UpdateUserHandler)
  router.HandleFunc("DELETE /user/{ids}", us.RemoveUserHandler)
}
