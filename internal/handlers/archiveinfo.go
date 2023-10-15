package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"test/internal/service"
)

type ArchiveHandler struct {
	archiveService *service.ArchiveService
}

func NewArchiveHandler(s *service.ArchiveService) *ArchiveHandler {
	return &ArchiveHandler{
		archiveService: s,
	}
}

func (h *ArchiveHandler) ArchiveInformationHandler(w http.ResponseWriter, r *http.Request) {
	// Check the POST request
	if r.Method != http.MethodPost {
		errorHandler(w, r, http.StatusMethodNotAllowed, "")
		return
	}

	// Recieve a file form a request
	file, header, err := r.FormFile("file")
	if err != nil {
		errorHandler(w, r, http.StatusBadRequest, "No file provided")
		return
	}
	defer file.Close()

	// Check if the file is zipped
	if !strings.HasSuffix(header.Filename, ".zip") {
		errorHandler(w, r, http.StatusBadRequest, "Invalid file format. Please provide a ZIP archive")
		return
	}

	// Use a service to get/extract/read info about archive
	archiveInfo, err := h.archiveService.GetArchiveInfo(file, header)
	if err != nil {
		errorHandler(w, r, http.StatusBadRequest, fmt.Sprintf("Failed to get archive information: %s", err)) //check the type of error

		return
	}

	// Return the result in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(archiveInfo)
}
