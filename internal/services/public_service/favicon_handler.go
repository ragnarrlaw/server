package publicservice

import (
	"net/http"
)

func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/home/bobby/playground/go/json_server/public/favicon.ico")
}
