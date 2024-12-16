package recommenderservice

import "net/http"

func (rs *RecommenderService) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /recommender/search", rs.SearchHandler)
	router.HandleFunc("POST /recommender/search/line-string", rs.LineStringSearchHandler)
}
