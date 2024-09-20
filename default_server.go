package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

func serveStatic(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RequestURI)
	// files, err := os.ReadDir("static/fonts")
	// if err != nil {
	// 	panic("Fail to read static files")
	// }

	// for _, file := range files {

	// }

}

func main_handler(w http.ResponseWriter, r *http.Request) {
	hasUsername := r.URL.Query().Has("username")
	if !hasUsername {
		w.Header().Add("x-missing-field", "username")
		w.WriteHeader(http.StatusBadRequest)
		// io.WriteString()
	}
	// username := r.URL.Query().Get("username")

}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, "{\"status\": \"up\"}")
}

func default_serve() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", main_handler)
	mux.HandleFunc("/health", healthCheck)
	mux.Handle("/static/", http.FileServer(http.Dir("static/fonts")))

	err := http.ListenAndServe(":8080", mux)

	if errors.Is(err, http.ErrServerClosed) {
		slog.Error("Server closed")
	} else if err != nil {
		slog.Error(fmt.Sprintf("Unknown error: %s", err))
		os.Exit(1)
	}

}

func main() {
	default_serve()
}
