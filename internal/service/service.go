// service/archive_service.go
package service

import (
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

func (as *ArchiveService) GetArchiveInfo(file io.Reader, header *multipart.FileHeader) (*models.ArchiveInfo, error) {
	return as.archiveRepository.ExtractAndSave(file, header)
}
