package handlers

import (
	"archive/zip"
	"net/http"
)

func (h *ArchiveHandler) CreateArchiveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		errorHandler(w, r, http.StatusMethodNotAllowed, "")

		return
	}
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		// http.Error(w, "Unable to parse form", http.StatusBadRequest)
		errorHandler(w, r, http.StatusBadRequest, "Unable to parse form")

		return
	}

	files := r.MultipartForm.File["files[]"]
	if len(files) == 0 {
		errorHandler(w, r, http.StatusBadRequest, "No files uploaded")
		// http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	err = h.archiveService.CreateArchive(files, zipWriter)
	if err != nil {
		// fmt.Println("no mime")
		errorHandler(w, r, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=output.zip")
}
