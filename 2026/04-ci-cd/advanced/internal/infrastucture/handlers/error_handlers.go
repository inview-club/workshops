package handlers

import (
	"net/http"
)

func NotAllowedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Method not allowed"
		http.Error(w, errorMessage, http.StatusMethodNotAllowed)
		return
	})

}

func NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Page not found"
		http.Error(w, errorMessage, http.StatusNotFound)
		return
	})
}
