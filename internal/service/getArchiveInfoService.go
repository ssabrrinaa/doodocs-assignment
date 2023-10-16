package service

import (
	"archive/zip"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"

	"test/models"
)

func (r *ArchiveService) GetArchiveInfo(file io.Reader, header *multipart.FileHeader) (*models.ArchiveInfo, error) {
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

		// Fill the info about each file
		fileInfo := models.FileInfo{
			FilePath: filePath_extracted,
			Size:     int64(size),
			MimeType: mime.TypeByExtension(filepath.Ext(file_extracted.Name)),
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

// func GetMimeType(zipFile *zip.File) (string, error) {
// 	file, err := zipFile.Open()
// 	if err != nil {
// 		return "", err
// 	}
// 	defer file.Close()

// 	buffer := make([]byte, 512) // Read the first 512 bytes to determine the MIME type
// 	_, err = file.Read(buffer)
// 	if err != nil && err != io.EOF {
// 		return "", err
// 	}

// 	// Determine the MIME type
// 	mimeType := http.DetectContentType(buffer)

// 	return mimeType, nil
// }
