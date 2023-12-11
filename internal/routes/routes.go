package routes

import (
	"fmt"
	"net/http"
	"strings"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/echo", echoHandler)
	mux.HandleFunc("/echo/", echoHandler)

	return mux
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the homepage!")
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	body := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, "/echo"), "/")
	fmt.Fprintln(w, body)
}
