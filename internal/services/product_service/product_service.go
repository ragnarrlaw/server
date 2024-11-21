package productservice

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	productrepository "github.com/raganrrlaw/server/internal/services/product_service/product_repository"
	"github.com/raganrrlaw/server/internal/types"
)

type ProductService struct {
	productRepository productrepository.ProductRepository
}

func NewProductService(productRepository productrepository.ProductRepository) *ProductService {
	return &ProductService{
		productRepository: productRepository,
	}
}

func (ps *ProductService) AddProductHandler(w http.ResponseWriter, r *http.Request) {
	payload := &types.ProductPayload{}
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	/** TODO: PAYLOAD VALIDATION */
	/** TODO: GENERIC RESPONSE ERROR STRUCTURE ORGANIZATION */

	if storeId := r.PathValue("storeId"); storeId != "" {
		if product, err := ps.productRepository.Add(ctx, storeId, payload); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			if b, err := json.Marshal(product); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				w.Write(b)
			}
		}
	} else {
		http.Error(w, "Invalid store id", http.StatusBadRequest)
	}
}

func (ps *ProductService) GetAllProductsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if storeId := r.PathValue("storeId"); storeId != "" {
		if products, err := ps.productRepository.GetAll(ctx, storeId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			if b, err := json.Marshal(products); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(b)
			}
		}
	} else {
		http.Error(w, "Invalid store id", http.StatusBadRequest)
	}
}

func (ps *ProductService) GetProductHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if storeId := r.PathValue("storeId"); storeId != "" {
		if productId := r.PathValue("productId"); productId != "" {
			if product, err := ps.productRepository.GetById(ctx, storeId, productId); err != nil {
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
	} else {
		http.Error(w, "Invalid store id", http.StatusBadRequest)
	}
}

func (ps *ProductService) UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	payload := &types.ProductPayload{}
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if storeId := r.PathValue("storeId"); storeId != "" {
		if productId := r.PathValue("productId"); productId != "" {
			if product, err := ps.productRepository.Update(ctx, storeId, productId, payload); err != nil {
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
	} else {
		http.Error(w, "Invalid store id", http.StatusBadRequest)
	}
}

func (ps *ProductService) RemoveProductHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if storeId := r.PathValue("storeId"); storeId != "" {
		if productId := r.PathValue("productId"); productId != "" {
			if err := ps.productRepository.Remove(
				ctx,
				storeId,
				[]string{productId},
			); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusNoContent)
			}
		} else {
			http.Error(w, "Invalid product id", http.StatusBadRequest)
		}
	} else {
		http.Error(w, "Invalid store id", http.StatusBadRequest)
	}
}

// product discount handlers
func (ps *ProductService) GetProductDiscountsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if storeId := r.PathValue("storeId"); storeId != "" {
		if productId := r.PathValue("productId"); productId != "" {
			if discounts, err := ps.productRepository.GetProductDiscounts(
				ctx,
				storeId,
				productId,
			); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				if b, err := json.Marshal(discounts); err != nil {
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
	} else {
		http.Error(w, "Invalid store id", http.StatusBadRequest)
	}
}

func (ps *ProductService) AddProductDiscountHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	payload := &types.DiscountPayload{}

	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		if storeId := r.PathValue("storeId"); storeId != "" {
			if productId := r.PathValue("productId"); productId != "" {
				if discount, err := ps.productRepository.AddProductDiscount(
					ctx,
					storeId,
					productId,
					payload,
				); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				} else {
					if b, err := json.Marshal(discount); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					} else {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusCreated)
						w.Write(b)
					}
				}
			} else {
				http.Error(w, "Invalid product id", http.StatusBadRequest)
			}
		} else {
			http.Error(w, "Invalid store id", http.StatusBadRequest)
		}
	}
}

func (ps *ProductService) UpdateProductDiscountHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	payload := &types.DiscountPayload{}

	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if storeId := r.PathValue("storeId"); storeId != "" {
		if productId := r.PathValue("productId"); productId != "" {
			if discountId := r.PathValue("discountId"); discountId != "" {
				if discount, err := ps.productRepository.UpdateDiscount(
					ctx,
					storeId,
					productId,
					discountId,
					payload,
				); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				} else {
					if b, err := json.Marshal(discount); err != nil {
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
		} else {
			http.Error(w, "Invalid product id", http.StatusBadRequest)
		}
	} else {
		http.Error(w, "Invalid store id", http.StatusBadRequest)
	}

}

func (ps *ProductService) RemoveProductDiscountHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if storeId := r.PathValue("storeId"); storeId != "" {
		if productId := r.PathValue("productId"); productId != "" {
			if discountId := r.PathValue("discountId"); discountId != "" {
				if err := ps.productRepository.RemoveDiscount(
					ctx,
					storeId,
					productId,
					[]string{discountId},
				); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				} else {
					w.WriteHeader(http.StatusNoContent)
				}
			} else {
				http.Error(w, "Invalid discount id", http.StatusBadRequest)
			}
		} else {
			http.Error(w, "Invalid product id", http.StatusBadRequest)
		}
	} else {
		http.Error(w, "Invalid store id", http.StatusBadRequest)
	}
}

func (ps *ProductService) RemoveAllProductDiscountHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if storeId := r.PathValue("storeId"); storeId != "" {
		if productId := r.PathValue("productId"); productId != "" {
			if err := ps.productRepository.RemoveAllDiscounts(
				ctx,
				storeId,
				productId,
			); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusNoContent)
			}
		} else {
			http.Error(w, "Invalid product id", http.StatusBadRequest)
		}
	} else {
		http.Error(w, "Invalid store id", http.StatusBadRequest)
	}
}

// product categories
func (ps *ProductService) GetAllProductCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if categories, err := ps.productRepository.ProductCategories(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		if b, err := json.Marshal(categories); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(b)
		}
	}
}
