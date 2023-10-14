// repository/archive_repository.go
package repository

import (
	"archive/zip"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"test/models"
)

type ArchiveRepository struct{}

func NewArchiveRepository() *ArchiveRepository {
	return &ArchiveRepository{}
}

func (ar *ArchiveRepository) ExtractAndSave(file io.Reader, header *multipart.FileHeader) (*models.ArchiveInfo, error) {
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
	filePath := filepath.Join(tempDir, header.Filename)
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
	for _, file := range reader.File {
		filePath := filepath.Join("extracted_files", file.Name)
		size := file.UncompressedSize64
		totalSize += int64(size)

		// Create an info for the file
		fileInfo := models.FileInfo{
			FilePath: filePath,
			Size:     int64(size),
			MimeType: GetMimeType(file.Name),
		}
		filesInfo = append(filesInfo, fileInfo)
	}

	// Fil the info about archive
	archiveInfo.Filename = filepath.Base(filePath)
	archiveInfo.ArchiveSize = float64(totalSize)
	archiveInfo.TotalSize = float64(totalSize)
	archiveInfo.TotalFiles = float64(len(filesInfo))
	archiveInfo.Files = filesInfo

	return archiveInfo, nil
}

func GetMimeType(fileName string) string {
	// Пример простой логики для определения MIME-типа
	if strings.HasSuffix(fileName, ".jpg") {
		return "image/jpeg"
	} else if strings.HasSuffix(fileName, ".docx") {
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	}
	return "application/octet-stream"
}
