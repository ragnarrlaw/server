package productservice

import "net/http"

func (ps *ProductService) RegisterRoutes(router *http.ServeMux) {
	// Product routes
	router.HandleFunc("GET /product", ps.GetProductsHandler)
	router.HandleFunc("GET /product/{productId}", ps.GetProductHandler)
	router.HandleFunc("GET /product/search", ps.GetSearchProductByNameHandler)
}
