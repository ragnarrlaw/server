package userservice

import "net/http"

func (us *UserService) RegisterRoutes(router *http.ServeMux) {
	// user routes
	router.HandleFunc("GET /user", us.GetAllUsersHandler)
	router.HandleFunc("GET /user/{id}", us.GetUserHandler)
	router.HandleFunc("POST /user", us.AddUserHandler)
	router.HandleFunc("PATCH /user/{id}", us.UpdateUserHandler)
	router.HandleFunc("DELETE /user/{id}", us.RemoveUserHandler)

	// user preferences
	router.HandleFunc("GET /user/{id}/preferences", us.GetUserPreferencesHandler)
	router.HandleFunc("PATCH /user/{id}/preferences", us.UpdateUserPreferencesHandler) // clear the list of preferences if an empty list is received - delete also happens in this route

	// user added product lists
	router.HandleFunc("GET /user/{id}/products-list", nil)            // - get all
	router.HandleFunc("GET /user/{id}/products-list/{listId}", nil)   // - get one
	router.HandleFunc("POST /user/{id}/products-list", nil)           // - add list
	router.HandleFunc("PATCH /user/{id}/products-list/{listId}", nil) // - update list - if an empty list is received then don't update
	router.HandleFunc("DELETE /user/{id}/products-list", nil)         // - delete lists e.g. : /user/{id}/products-list?listId=uuid1&listId=uuid2
}
