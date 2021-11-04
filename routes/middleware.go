package routes

import (
	"fmt"
	"net/http"
)

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Trove Test App - Version 0.0255\n")
}
