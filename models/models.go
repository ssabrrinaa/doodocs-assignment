// models/models.go
package models

type FileInfo struct {
	FilePath string `json:"file_path"`
	Size     int64  `json:"size"`
	MimeType string `json:"mimetype"`
}

type ArchiveInfo struct {
	Filename    string     `json:"filename"`
	ArchiveSize float64    `json:"archive_size"`
	TotalSize   float64    `json:"total_size"`
	TotalFiles  float64    `json:"total_files"`
	Files       []FileInfo `json:"files"`
}
