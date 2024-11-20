package productservice

import "net/http"

func (ps *ProductService) RegisterRoutes(router *http.ServeMux) {
	// Product routes
	router.HandleFunc("GET /store/{storeId}/product", ps.GetAllProductsHandler)
	router.HandleFunc("GET /store/{storeId}/product/{productId}", ps.GetProductHandler)
	router.HandleFunc("POST /store/{storeId}/product", ps.AddProductHandler)
	router.HandleFunc("PATCH /store/{storeId}/product/{productId}", ps.UpdateProductHandler)
	router.HandleFunc("DELETE /store/{storeId}/product/{productId}", ps.RemoveProductHandler)

	// Discount routes
	router.HandleFunc("POST /store/{storeId}/product/{productId}/discount", ps.AddProductDiscountHandler)
	router.HandleFunc("PATCH /store/{storeId}/product/{productId}/discount/{discountId}", ps.UpdateProductDiscountHandler)
	router.HandleFunc("DELETE /store/{storeId}/product/{productId}/discount/{discountId}", ps.RemoveProductDiscountHandler)
	router.HandleFunc("DELETE /store/{storeId}/product/{productId}/discount", ps.RemoveAllProductDiscountHandler)

	// Category routes
	router.HandleFunc("GET /product/category", ps.GetAllProductCategoriesHandler)
}
