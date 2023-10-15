package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

func (h *ArchiveHandler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound, "")
		return
	}

	if r.Method != http.MethodGet {
		errorHandler(w, r, http.StatusMethodNotAllowed, "")
		return
	}

	t, err := template.ParseFiles("templates/home.html")
	if err != nil {
		fmt.Println("here")
		errorHandler(w, r, http.StatusInternalServerError, "")
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}
