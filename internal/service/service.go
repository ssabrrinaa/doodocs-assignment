package service

import (
	"archive/zip"
	"io"
	"mime/multipart"

	"test/internal/repository"
	"test/models"
)

type ArchiveService struct {
	archiveRepository *repository.ArchiveRepository
}

func NewArchiveService(r *repository.ArchiveRepository) *ArchiveService {
	return &ArchiveService{
		archiveRepository: r,
	}
}

func (s *ArchiveService) GetArchiveInfo(file io.Reader, header *multipart.FileHeader) (*models.ArchiveInfo, error) {
	return s.archiveRepository.ExtractAndSave(file, header)
}

func (s *ArchiveService) CreateArchive(files []*multipart.FileHeader, zipWriter *zip.Writer) error {
	return s.archiveRepository.CreateArchiveFile(files, zipWriter)
}
