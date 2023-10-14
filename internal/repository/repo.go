// repository/archive_repository.go
package repository

import (
	"archive/zip"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

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
	for _, file_ex := range reader.File {
		filePath := filepath.Join("extracted_files", file_ex.Name)
		size := file_ex.UncompressedSize64
		totalSize += int64(size)

		mimetype, err := GetMimeType(file)
		if err != nil {
			return archiveInfo, err
		}
		// Create an info for the file
		fileInfo := models.FileInfo{
			FilePath: filePath,
			Size:     int64(size),
			MimeType: mimetype,
		}
		filesInfo = append(filesInfo, fileInfo)
	}

	// Fil the info about archive
	archiveInfo.Filename = filepath.Base(filePath)
	archiveInfo.ArchiveSize, err = getCompressedSize(filePath)
	archiveInfo.TotalSize = float64(totalSize)
	archiveInfo.TotalFiles = float64(len(filesInfo))
	archiveInfo.Files = filesInfo

	return archiveInfo, nil
}

func getCompressedSize(filePath string) (float64, error) {
	// Open the archived file
	archiveFile, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer archiveFile.Close()

	// Get file information
	fileInfo, err := archiveFile.Stat()
	if err != nil {
		return 0, err
	}

	// Retrieve the compressed size from file information
	return float64(fileInfo.Size()), nil
}

func GetMimeType(file io.Reader) (string, error) {
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	contentType := http.DetectContentType(buffer[:n])

	return contentType, nil
}
