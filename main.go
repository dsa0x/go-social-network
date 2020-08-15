package main

import (
	"net/http"

	"github.com/dsa0x/go-social-network/internal/router"
)

func main() {

	port := "8080"

	http.ListenAndServe(":"+port, router.APIRouter)
}
