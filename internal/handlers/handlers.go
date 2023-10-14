// handlers/archive_handler.go
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

func (ah *ArchiveHandler) HandleArchiveInformation(w http.ResponseWriter, r *http.Request) {
	// Check the POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Recieve a file form a request
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "No file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check if the file is zipped
	if !strings.HasSuffix(header.Filename, ".zip") {
		http.Error(w, "Invalid file format. Please provide a ZIP archive", http.StatusBadRequest)
		return
	}

	// Use a service to get/extract/read info about achive
	archiveInfo, err := ah.archiveService.GetArchiveInfo(file, header)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get archive information: %s", err), http.StatusInternalServerError)
		return
	}

	// Return the result in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(archiveInfo)
}
