package handlers

import (
	"test/internal/service"
)

type ArchiveHandler struct {
	archiveService *service.ArchiveService
}

func NewArchiveHandler(s *service.ArchiveService) *ArchiveHandler {
	return &ArchiveHandler{
		archiveService: s,
	}
}
