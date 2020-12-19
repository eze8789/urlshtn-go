package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Routes() http.Handler {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", home).Methods(http.MethodGet)

	return rtr
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to urlshtn-go, this is the supported endpoint list")
}
