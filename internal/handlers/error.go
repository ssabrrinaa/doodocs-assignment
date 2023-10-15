package handlers

import (
	"html/template"
	"net/http"
)

type PageData struct {
	StatusCode int
	StatusText string
	ErrText    string
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int, errtext string) {
	text := http.StatusText(status)
	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	tmpl.Execute(w, PageData{StatusCode: status, StatusText: text, ErrText: errtext})
}
