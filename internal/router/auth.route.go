package router

import (
	"github.com/dsa0x/go-social-network/internal/handler"
)

func init() {

	APIRouter.HandleFunc("/login", handler.Login).Methods("GET", "POST")
	APIRouter.HandleFunc("/logout", handler.Login).Methods("GET", "POST")
	APIRouter.HandleFunc("/signup", handler.SignUp).Methods("GET", "POST")

}
