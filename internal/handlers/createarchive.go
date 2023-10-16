package handlers

import (
	"archive/zip"
	"net/http"
)

func (h *ArchiveHandler) CreateArchiveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorHandler(w, r, http.StatusMethodNotAllowed, "")
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		errorHandler(w, r, http.StatusBadRequest, "Unable to parse form")
		return
	}

	// Extract the uploaded files from the multipart form
	files := r.MultipartForm.File["files[]"]
	if len(files) == 0 {
		errorHandler(w, r, http.StatusBadRequest, "No files uploaded")
		return
	}

	// Create a new ZIP writer that will write the ZIP archive to the HTTP response writer 'w'.

	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	err = h.archiveService.CreateArchive(files, zipWriter)
	if err != nil {
		errorHandler(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// Set response headers to indicate that the response is a ZIP file for download.
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=output.zip")

}
