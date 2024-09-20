package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

func loadTemplate(filename string) (string, error) {
	templateFolderName := "template"
	templateFolder, err := os.ReadDir(templateFolderName)
	if err != nil {
		return "", err
	}
	for _, file := range templateFolder {
		if file.Name() == filename {
			content, err := os.ReadFile(path.Join(templateFolderName, file.Name()))
			if err != nil {
				slog.Error(fmt.Sprintf("Error loading %s: %s", filename, err))
			}
			return string(content), nil
		}
	}
	return "", fmt.Errorf("template %s not found in template folder", filename)
}

func parseTemplate(filename string, variables map[string]interface{}) (string, error) {
	content, err := loadTemplate(filename)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("template").Parse(content)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %s", err)
	}

	buff := new(strings.Builder)

	err = tmpl.Execute(buff, variables)
	if err != nil {
		return "", err
	}
	return buff.String(), nil
}

type ResponseWriter struct {
	startTimeRequest time.Time
	w                *http.ResponseWriter
	r                *http.Request
	status           int
	body             string
	headers          map[string]string
}

func NewResponseWriter(w *http.ResponseWriter, r *http.Request) ResponseWriter {
	return ResponseWriter{
		startTimeRequest: time.Now(),
		status:           http.StatusOK,
		headers:          map[string]string{},
		body:             "",
		w:                w,
		r:                r,
	}
}

func (r *ResponseWriter) Build(body string, status int, headers map[string]string) {
	r.body = body
	r.status = status
	r.headers = headers
}

func (r *ResponseWriter) BuildAndSend(body string, status int, headers map[string]string) {
	r.Build(body, status, headers)
	r.Send()
}

func (r *ResponseWriter) Send() {
	for k, v := range r.headers {
		(*r.w).Header().Add(k, v)
	}
	(*r.w).WriteHeader(r.status)
	io.WriteString((*r.w), r.body)
	requestTime := time.Since(r.startTimeRequest)
	slog.Info(fmt.Sprintf("| %-3d | %-12v | %-15s | %-6s | %-30s",
		r.status, requestTime, strings.Split(r.r.RemoteAddr, ":")[0], r.r.Method, r.r.URL))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	responseWriter := NewResponseWriter(&w, r)
	hasUsername := r.URL.Query().Has("username")
	var headers map[string]string

	if !hasUsername {
		headers = map[string]string{
			"x-missing-field": "username",
			"content-type":    "text/html",
		}
		notFound, _ := parseTemplate("404.tmpl", map[string]interface{}{"message": "Error! Missing username param"})
		responseWriter.BuildAndSend(notFound, http.StatusBadRequest, headers)
		return
	}

	username := r.URL.Query().Get("username")
	var user User
	var err error

	user, err = mainDatabase.searchForUser(username)
	if err != nil {
		user, err = mainDatabase.insertUserIntoDB(username)
		if err != nil {
			headers = map[string]string{
				"content-type": "application/json",
			}
			responseWriter.BuildAndSend(fmt.Sprintf("{\"message\": \"%s\"}", err), http.StatusInternalServerError, headers)
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
	mainDatabase.updateUserViewCount(user)
	responseWriter.BuildAndSend(generateSVG(user.counter, loadedImages).String(), http.StatusOK, headers)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	writer := NewResponseWriter(&w, r)
	writer.BuildAndSend("{\"status\": \"up\"}", http.StatusOK, map[string]string{"Content-Type": "application/json"})
}

func serve() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/health", healthCheck)
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
