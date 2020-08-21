package main

import (
	"log"
	"net/http"

	"github.com/dsa0x/go-social-network/internal/router"
)

func main() {
	port := "8080"
	staticDir := "/public/"
	router.APIRouter.
		PathPrefix(staticDir).
		Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))
	log.Println("Listening on :8080...")
	http.ListenAndServe(":"+port, router.APIRouter)
}
