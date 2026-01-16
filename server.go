package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/kamuridesu/gomechan/core/response"
	"github.com/kamuridesu/gomechan/core/routes"
	"github.com/kamuridesu/gomechan/core/templates"
)

var (
	Templates, _ = templates.LoadTemplateFolder("./template")
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	responseWriter := response.New(&w, r)
	var headers map[string]string

	if r.URL.Path != "/" {
		headers = map[string]string{
			"content-type": "text/html",
		}
		notFound := Templates.LoadHTML("404.tmpl", map[string]any{"message": "Error! Page not found"})
		responseWriter.SetHeaders(headers).Build(http.StatusNotFound, []byte(notFound)).Send()
		return
	}

	hasUsername := r.URL.Query().Has("username")
	if !hasUsername {
		headers = map[string]string{
			"x-missing-field": "username",
			"content-type":    "text/html",
		}
		notFound := Templates.LoadHTML("404.tmpl", map[string]any{"message": "Error! Page not found"})
		responseWriter.SetHeaders(headers).Build(http.StatusBadRequest, []byte(notFound)).Send()
		return
	}

	username := r.URL.Query().Get("username")
	var user *User
	var err error

	user, err = mainDatabase.searchForUser(username)
	if err != nil {
		user, err = mainDatabase.insertUserIntoDB(username)
		if err != nil {
			responseWriter.SendAsJson(http.StatusInternalServerError, map[string]any{
				"message": err.Error(),
			})
			return
		}
	}

	headers = map[string]string{
		"Cache-Control": "no-cache, no-store, must-revalidate",
		"Vary":          "Accept-Encoding",
		"Pragma":        "no-cache",
		"Expires":       "0",
		"content-type":  "image/svg+xml",
	}
	responseWriter.SetHeaders(headers).Build(http.StatusOK, generateSVG(user.counter, loadedImages)).Send()
	mainDatabase.updateUserViewCount(user)
}

func serve() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	routes.AddHealthCheck(mux)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/fonts"))))

	slog.Info("Listening on 0.0.0.0:8080")
	err := http.ListenAndServe("0.0.0.0:8080", mux)

	if errors.Is(err, http.ErrServerClosed) {
		slog.Error("Server closed")
	} else if err != nil {
		slog.Error(fmt.Sprintf("Unknown error: %s", err))
		os.Exit(1)
	}
}
