package service

import (
	"archive/zip"
	"mime/multipart"
)

func (r *ArchiveService) SendFileRepo(files []*multipart.FileHeader, zipWriter *zip.Writer) error {
	return nil
}
