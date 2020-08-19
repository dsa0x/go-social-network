package main

import (
	"net/http"

	"github.com/dsa0x/go-social-network/internal/router"
)

func main() {

	port := "8080"
	staticDir := "/public/"
	// fs := http.StripPrefix("/public/", http.FileServer(http.Dir("./public")))
	// fmt.Println(fs)

	// // http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("./public"))))
	// router.APIRouter.Handle("/public/", fs)
	// // router.APIRouter.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("./public"))))
	// http.ListenAndServe(":"+port, router.APIRouter)
	// fs := http.FileServer(http.Dir("./public"))
	// router.APIRouter.Handle("/public/", http.StripPrefix("/public/", fs))
	router.APIRouter.
		PathPrefix(staticDir).
		Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))

	http.ListenAndServe(":"+port, router.APIRouter)
	// log.Println("Listening on :8080...")
	// err := http.ListenAndServe(":8080", nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
