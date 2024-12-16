package storeservice

import "net/http"

func (ss *StoreService) RegisterRoutes(router *http.ServeMux) {
	// store routes
	router.HandleFunc("GET /store", ss.GetAllStoresHandler)
	router.HandleFunc("GET /store/{id}", ss.GetStoreHandler)
	router.HandleFunc("GET /store/stat", ss.GetStoreStatHandler)
}
