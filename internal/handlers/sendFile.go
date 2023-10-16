package handlers

import (
	"fmt"
	"net/http"
	"strings"
)

func (h *ArchiveHandler) SendFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorHandler(w, r, http.StatusMethodNotAllowed, "")
		return
	}

	// Parse the incoming multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		errorHandler(w, r, http.StatusBadRequest, "Unable to parse form")
		return
	}

	// Get the file and file header from the form
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		errorHandler(w, r, http.StatusBadRequest, "Unable to get file")
		return
	}
	defer file.Close()

	// Get the list of emails from the form
	emailsStr := r.FormValue("emails")
	emails := strings.Split(emailsStr, ",")

	err = h.archiveService.SendFile(file, fileHeader, emails)
	if err != nil {
		errorHandler(w, r, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "File sent successfully")
}
