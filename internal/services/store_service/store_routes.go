package storeservice

import "net/http"

func (ss *StoreService) RegisterRoutes(router *http.ServeMux) {
	// store routes
	router.HandleFunc("GET /store", ss.GetAllStoresHandler)
	router.HandleFunc("GET /store/{storeId}", ss.GetStoreHandler)
	router.HandleFunc("GET /store/stat", ss.GetStoreStatHandler)
	router.HandleFunc("POST /store/search", ss.PostStoreSearchHandler)

	// store product routes
	router.HandleFunc("GET /store/{storeId}/product", ss.GetStoreProductsHandler)
	router.HandleFunc("POST /store/{storeId}/product", ss.PostStoreProductHandler)
	router.HandleFunc("PATCH /store/{storeId}/product/{productId}", ss.UpdateStoreProductHandler)
	router.HandleFunc("DELETE /store/{storeId}/product", ss.DeleteStoreProductsHandler)

	// store discount routes
	router.HandleFunc("GET /store/{storeId}/discount", ss.GetStoreDiscountHandler)
	router.HandleFunc("POST /store/{storeId}/discount", ss.PostStoreDiscountHandler)
	router.HandleFunc("PATCH /store/{storeId}/discount/{discountId}", ss.UpdateStoreDiscountHandler)
	router.HandleFunc("DELETE /store/{storeId}/discount", ss.DeleteStoreDiscountHandler)
}
