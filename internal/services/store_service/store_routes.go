package storeservice

import "net/http"

func (ss *StoreService) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /store", ss.GetAllStoresHandler)
	router.HandleFunc("GET /store/{id}", ss.GetStoreHandler)
	router.HandleFunc("POST /store", ss.AddStoreHandler)
	router.HandleFunc("PATCH /store/{id}", ss.UpdateStoreHandler)
	router.HandleFunc("DELETE /store/{id}", ss.RemoveStoreHandler)
}
