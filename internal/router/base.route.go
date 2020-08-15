package router

import (
	"github.com/dsa0x/go-social-network/internal/handler"
	"github.com/gorilla/mux"
)

// APIRouter is the main router
var APIRouter = mux.NewRouter()

func init() {

	APIRouter.Get("/").HandlerFunc(handler.Home)

}
