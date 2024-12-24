package storeservice

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	uuid "github.com/google/uuid"
	storerepository "github.com/ragnarrlaw/server/internal/services/store_service/store_repository"
	"github.com/ragnarrlaw/server/internal/types/search"
	"github.com/ragnarrlaw/server/internal/types/store"
)

type StoreService struct {
	storeRepository storerepository.StoreRepository
}

func NewStoreService(storeRepository storerepository.StoreRepository) *StoreService {
	return &StoreService{
		storeRepository: storeRepository,
	}
}

// store handlers-------------------------------------------------------------------------
func (ss *StoreService) GetAllStoresHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if stores, err := ss.storeRepository.Get(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		if totalRecords, err := ss.storeRepository.GetTotalNumberOfStores(ctx); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			if b, err := json.Marshal(search.Payload{
				Data: stores,
				Meta: search.Meta{
					Total: totalRecords,
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

func (ss *StoreService) GetStoreStatHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	stats, err := ss.storeRepository.Stats(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	total, err := ss.storeRepository.GetTotalNumberOfStores(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(search.Payload{
		Data: stats,
		Meta: search.Meta{
			Total: total,
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (ss *StoreService) GetStoreHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if storeId := r.PathValue("storeId"); storeId != "" {
		store, err := ss.storeRepository.GetById(ctx, storeId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			if b, err := json.Marshal(search.Payload{
				Data: store,
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(b)
			}
		}
	}
}

func (ss *StoreService) PostStoreSearchHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	payload := &store.FilterCriteria{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stores, err := ss.storeRepository.SearchStores(ctx, payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(search.Payload{
		Data: stores,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// store product handlers-----------------------------------------------------------------
func (ss *StoreService) GetStoreProductsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	if storeId := r.PathValue("storeId"); storeId != "" {
		if products, err := ss.storeRepository.GetStoreProducts(ctx, storeId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			if total, err := ss.storeRepository.GetTotalNumberOfProductsInStore(ctx, storeId); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				if b, err := json.Marshal(search.Payload{
					Data: products,
					Meta: search.Meta{
						Total: total,
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
	} else {
		http.Error(w, "store id required", http.StatusBadRequest)
	}
}

func (ss *StoreService) PostStoreProductHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if storeId := r.PathValue("storeId"); storeId != "" {
		payload := &store.StoreProductPayload{}
		if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			if payload.ProductId.String() != "" {
				if data, err := ss.storeRepository.CreateStoreProduct(ctx, storeId, payload); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				} else {
					if b, err := json.Marshal(search.Payload{
						Data: data,
					}); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					} else {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						w.Write(b)
					}
				}
			} else {
				http.Error(w, "product id required", http.StatusBadRequest)
			}
		}
	} else {
		http.Error(w, "store id required", http.StatusBadRequest)
	}
}

func (ss *StoreService) UpdateStoreProductHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if storeId := r.PathValue("storeId"); storeId != "" {
		payload := &store.StoreProductPayload{}
		if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			if payload.ProductId.String() != "" {
				if data, err := ss.storeRepository.UpdateStoreProduct(ctx, storeId, payload.ProductId.String(), payload); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				} else {
					if b, err := json.Marshal(search.Payload{
						Data: data,
					}); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					} else {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						w.Write(b)
					}
				}
			} else {
				http.Error(w, "product id required", http.StatusBadRequest)
			}
		}
	} else {
		http.Error(w, "store id required", http.StatusBadRequest)
	}
}

func (ss *StoreService) DeleteStoreProductsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	if storeId := r.URL.Query().Get("storeId"); storeId != "" {
		if productIds := r.URL.Query()["productIds"]; len(productIds) > 0 {
			if err := ss.storeRepository.DeleteStoreProducts(ctx, storeId, productIds); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusNoContent)
			}
		} else {
			http.Error(w, "product ids required", http.StatusBadRequest)
		}
	} else {
		http.Error(w, "store id required", http.StatusBadRequest)
	}
}

// store discount handlers----------------------------------------------------------------
func (ss *StoreService) GetStoreDiscountHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if storeId := r.PathValue("storeId"); storeId != "" {
		if id, err := uuid.Parse(storeId); err != nil {
			http.Error(w, "store id required", http.StatusBadRequest)
		} else {
			if discounts, err := ss.storeRepository.GetStoreDiscounts(ctx, id); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				if total, err := ss.storeRepository.GetTotalNumberOfDiscountsInStore(ctx, id); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				} else {
					if b, err := json.Marshal(search.Payload{
						Data: discounts,
						Meta: search.Meta{
							Total: total,
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

	} else {
		http.Error(w, "store id required", http.StatusBadRequest)
	}
}

func (ss *StoreService) PostStoreDiscountHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if storeId := r.PathValue("storeId"); storeId != "" {
		payload := &store.StoreDiscountPayload{}
		if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			if id, err := uuid.Parse(storeId); err != nil {
				http.Error(w, "invalid store id format", http.StatusBadRequest)
			} else {
				if data, err := ss.storeRepository.CreateStoreDiscount(ctx, id, payload); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				} else {
					if b, err := json.Marshal(search.Payload{
						Data: data,
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
	} else {
		http.Error(w, "store id required", http.StatusBadRequest)
	}
}

func (ss *StoreService) UpdateStoreDiscountHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if storeId := r.PathValue("storeId"); storeId != "" {
		payload := &store.StoreDiscountPayload{}
		if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			if id, err := uuid.Parse(storeId); err != nil {
				http.Error(w, "invalid store id format", http.StatusBadRequest)
			} else {
				if data, err := ss.storeRepository.UpdateStoreDiscount(ctx, id, payload); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				} else {
					if b, err := json.Marshal(search.Payload{
						Data: data,
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
	} else {
		http.Error(w, "store id required", http.StatusBadRequest)
	}
}

func (ss *StoreService) DeleteStoreDiscountHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if storeId := r.URL.Query().Get("storeId"); storeId != "" {
		if discountIds := r.URL.Query()["discountIds"]; len(discountIds) > 0 {
			if id, err := uuid.Parse(storeId); err != nil {
				http.Error(w, "invalid store id format", http.StatusBadRequest)
			} else {
				ids := make([]uuid.UUID, 0, len(discountIds))
				for i, id := range discountIds {
					if uuid, err := uuid.Parse(id); err != nil {
						http.Error(w, "invalid discount id format", http.StatusBadRequest)
					} else {
						ids[i] = uuid
					}
				}
				if err := ss.storeRepository.DeleteStoreDiscounts(ctx, id, ids); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				} else {
					w.WriteHeader(http.StatusNoContent)
				}
			}
		} else {
			http.Error(w, "discount ids required", http.StatusBadRequest)
		}
	} else {
		http.Error(w, "store id required", http.StatusBadRequest)
	}
}
