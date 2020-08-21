package router

import (
	"github.com/dsa0x/go-social-network/internal/handler"
)

func init() {

	APIRouter.HandleFunc("/user/{id}", handler.Auth(handler.GetUser)).Methods("GET", "POST")
	APIRouter.HandleFunc("/user/{id}/follow", handler.Auth(handler.FollowUser)).Methods("GET", "POST")
	APIRouter.HandleFunc("/user/{id}/unfollow", handler.Auth(handler.UnfollowUser)).Methods("GET", "POST")
	APIRouter.HandleFunc("/users", handler.Auth(handler.GetAllUsers)).Methods("GET", "POST")

}
