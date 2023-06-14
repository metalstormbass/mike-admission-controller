package router

import (
	"github.com/gorilla/mux"
	"github.com/metalstormbass/mike-admission-controller/src/webhook"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", webhook.Validate).Methods("POST")
	return router
}
