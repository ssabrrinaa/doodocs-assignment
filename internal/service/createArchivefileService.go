package service

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"path/filepath"
)

func (r *ArchiveService) CreateArchive(files []*multipart.FileHeader, zipWriter *zip.Writer) error {
	for _, fileHeader := range files {
		// fmt.Println(fileHeader.Filename, "filename")
		is, err := isValidMIMEType(fileHeader)
		if err != nil {
			// http.Error(w, fmt.Sprintf("Failed to create archive files: %s", err), http.StatusInternalServerError) //check the error
			return err
		}
		if is {
			file, err := fileHeader.Open()
			if err != nil {
				log.Println("Error opening file:", err)
				return err
			}
			defer file.Close()

			zipFile, err := zipWriter.Create(fileHeader.Filename)
			if err != nil {
				log.Println("Error creating ZIP file:", err)
				return err
			}

			_, err = io.Copy(zipFile, file)
			if err != nil {
				log.Println("Error copying file to ZIP archive:", err)
				return err
			}
		} else {
			// log.Println()
			// http.Error(w, "Invalid MIME type:", http.StatusBadRequest)
			return fmt.Errorf("Invalid MIME type:%v", fileHeader.Filename)
		}
	}
	return nil
}

func isValidMIMEType(fileHeader *multipart.FileHeader) (bool, error) {
	// Здесь определите разрешенные MIME-типы файлов, которые вы хотите архивировать
	allowedMIMETypes := map[string]bool{
		"image/jpeg":      true,
		"image/png":       true,
		"application/xml": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	}
	mimeType := mime.TypeByExtension(filepath.Ext(fileHeader.Filename))

	// file, err := fileHeader.Open()
	// if err != nil {
	// 	fmt.Println("here", err)
	// 	return false, err
	// }
	// defer file.Close()

	// // Читаем первые 512 байт файла, чтобы определить MIME-тип
	// buffer := make([]byte, 512)
	// _, err = file.Read(buffer)
	// if err != nil {
	// 	return false, err
	// }

	// mimeType := http.DetectContentType(buffer)
	fmt.Println(mimeType, "is valid mimetype")
	return allowedMIMETypes[mimeType], nil
}
