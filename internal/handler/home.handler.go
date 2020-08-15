package handler

import (
	"fmt"
	"net/http"
)

// Home function for home handler
func Home(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "%s", "Hello World")
}
