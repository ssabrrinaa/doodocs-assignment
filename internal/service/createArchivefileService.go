package service

import (
	"archive/zip"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"path/filepath"
)

func (r *ArchiveService) CreateArchive(files []*multipart.FileHeader, zipWriter *zip.Writer) error {
	for _, fileHeader := range files {
		valid, err := isValidMIMEType(fileHeader)
		if err != nil {
			return err
		}

		if valid {
			// Open the file for reading.
			file, err := fileHeader.Open()
			if err != nil {
				return err
			}
			defer file.Close()

			// Create a new file within the ZIP archive with the same name as the original file.
			zipFile, err := zipWriter.Create(fileHeader.Filename)
			if err != nil {
				return err
			}

			// Copy the contents of the opened file to the ZIP file.
			_, err = io.Copy(zipFile, file)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("invalid mimetype:%v", fileHeader.Filename)
		}
	}

	return nil
}

func isValidMIMEType(fileHeader *multipart.FileHeader) (bool, error) {
	allowedMIMETypes := map[string]bool{
		"image/jpeg":      true,
		"image/png":       true,
		"application/xml": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	}
	mimeType := mime.TypeByExtension(filepath.Ext(fileHeader.Filename))

	return allowedMIMETypes[mimeType], nil
}
