package repository

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"test/models"
)

type ArchiveRepository struct{}

func NewArchiveRepository() *ArchiveRepository {
	return &ArchiveRepository{}
}

func (r *ArchiveRepository) ExtractAndSave(file io.Reader, header *multipart.FileHeader) (*models.ArchiveInfo, error) {
	tempDir := "extracted_files"
	archiveInfo := &models.ArchiveInfo{}
	filesInfo := []models.FileInfo{}
	totalSize := int64(0)

	// Create a temp directory for extracted files
	err := os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		return archiveInfo, err
	}
	defer os.RemoveAll(tempDir)

	// Save a file to the disk
	filePath := path.Join(tempDir, header.Filename)
	fmt.Println(filePath, "filepath")
	outFile, err := os.Create(filePath)
	if err != nil {
		return archiveInfo, err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		return archiveInfo, err
	}

	// Open an extracted file
	reader, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// Iterate over the files
	for _, file_extracted := range reader.File {
		// Check if the entry is a directory and Skip directories
		if file_extracted.FileInfo().IsDir() {
			continue
		}
		filePath_extracted := path.Join("extracted_files", file_extracted.Name)
		size := file_extracted.UncompressedSize64
		totalSize += int64(size)

		mimetype, err := GetMimeType(file_extracted)
		if err != nil {
			return archiveInfo, err
		}

		// Fill the info about each file
		fileInfo := models.FileInfo{
			FilePath: filePath_extracted,
			Size:     int64(size),
			MimeType: mimetype,
		}
		filesInfo = append(filesInfo, fileInfo)
	}

	// Fil the info about archive
	archiveInfo.Filename = filepath.Base(filePath)
	archiveInfo.ArchiveSize = float64(header.Size)
	archiveInfo.TotalSize = float64(totalSize)
	archiveInfo.TotalFiles = float64(len(filesInfo))
	archiveInfo.Files = filesInfo

	return archiveInfo, nil
}

// func GetMimeType(filePath string) (string, error) {
// 	ext := strings.ToLower(filepath.Ext(filePath))
// 	mimeType := mime.TypeByExtension(ext)
// 	if mimeType == "" {
// 		mimeType = "application/octet-stream"
// 	}
// 	return mimeType, nil
// }

func GetMimeType(zipFile *zip.File) (string, error) {
	file, err := zipFile.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := make([]byte, 512) // Read the first 512 bytes to determine the MIME type
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	// Determine the MIME type
	mimeType := http.DetectContentType(buffer)

	return mimeType, nil
}

func (r *ArchiveRepository) CreateArchiveFile(files []*multipart.FileHeader, zipWriter *zip.Writer) error {
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

	file, err := fileHeader.Open()
	if err != nil {
		fmt.Println("here", err)
		return false, err
	}
	defer file.Close()

	// Читаем первые 512 байт файла, чтобы определить MIME-тип
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return false, err
	}

	mimeType := http.DetectContentType(buffer)
	fmt.Println(mimeType)
	return allowedMIMETypes[mimeType], err
}
