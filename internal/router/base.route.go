package router

import (
	"github.com/dsa0x/go-social-network/internal/handler"
	"github.com/gorilla/mux"
)

// APIRouter is the main router
var APIRouter = mux.NewRouter()

func init() {

	APIRouter.HandleFunc("/", handler.Auth(handler.Home)).Methods("GET", "POST")
	// APIRouter.HandleFunc("/", handler.Home).Methods("GET")
	// APIRouter.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("../../public"))))
	// APIRouter.HandleFunc("/logged", handler.Auth(handler.Home)).Methods("GET", "POST")

}
