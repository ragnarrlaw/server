package publicservice

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/raganrrlaw/server/internal/config"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: config.WELCOME_MESSAGE,
	})
	if err != nil {
		log.Fatalf("Error occurred while marshalling the data: %s\n", err.Error())
		http.Error(w, "Internal server error occurred", http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
