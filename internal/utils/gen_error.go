package utils

import (
	"log"
	"net/http"
	"runtime"
)

/**
  Add the generic error handlers for the repositories
  These are added to maintain the structure of the repositories
  despite the data storage changes
*/

func ErrorHandler(w http.ResponseWriter, err error, statusCode int, message string) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Printf("Error in  %s at %s:%d: %v\n", message, file, line, err)
	} else {
		log.Printf("Error in %s: %v\n", message, err)
	}
  /**
    TODO: Make the response a JSON response
  */
	http.Error(w, message, statusCode)
}
