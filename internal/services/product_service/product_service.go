package productservice

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	productrepository "github.com/raganrrlaw/server/internal/services/product_service/product_repository"
	"github.com/raganrrlaw/server/internal/types/search"
)

type ProductService struct {
	productRepository productrepository.ProductRepository
}

func NewProductService(productRepository productrepository.ProductRepository) *ProductService {
	return &ProductService{
		productRepository: productRepository,
	}
}

func (ps *ProductService) GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if products, err := ps.productRepository.Get(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		if totalProducts, err := ps.productRepository.GetTotalNumberOfProducts(ctx); err != nil {
		} else {
			if b, err := json.Marshal(search.Payload{
				Data: products,
				Meta: search.Meta{
					Total: totalProducts,
				},
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(b)
			}
		}
	}
}

func (ps *ProductService) GetProductHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if productId := r.PathValue("productId"); productId != "" {
		if product, err := ps.productRepository.GetById(ctx, productId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			if b, err := json.Marshal(product); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(b)
			}
		}
	} else {
		http.Error(w, "Invalid product id", http.StatusBadRequest)
	}
}

func (ps *ProductService) GetSearchProductByNameHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if term := r.URL.Query().Get("term"); term != "" {
		if products, err := ps.productRepository.SearchProductByName(ctx, term); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			if b, err := json.Marshal(search.Payload{
				Data: products,
				Meta: search.Meta{
					Total: uint(len(*products)),
				},
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(b)
			}
		}
	} else {
		http.Error(w, "Invalid product id", http.StatusBadRequest)
	}
}
