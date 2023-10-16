package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/smtp"
	"os"
	"path/filepath"

	"github.com/scorredoira/email"
)

const (
	smtpServer   = "smtp.gmail.com"
	smtpPort     = "587"
	smtpUsername = "bakirova200024@gmail.com"
	smtpPassword = "zevghlaxkgwpkdic"
)

func (r *ArchiveService) SendFile(file multipart.File, fileHeader *multipart.FileHeader, emails []string) error {

	// Check if the file MIME type is valid
	if !isValidFileType(fileHeader) {
		// errorHandler(w, r, http.StatusBadRequest, "Invalid file type")
		return fmt.Errorf("Invalid MIME type:%v", fileHeader.Filename)
	}

	// Создаем временную директорию
	tempDir := os.TempDir()

	// Определяем путь для сохранения файла во временной директории
	filePath := filepath.Join(tempDir, fileHeader.Filename)

	// Создаем новый файл для сохранения во временной директории
	destination, err := os.Create(filePath)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	defer destination.Close()

	// Копируем содержимое файла из запроса в созданный файл
	_, err = io.Copy(destination, file)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	// Send email with the file attachment
	err = sendFileToEmails(filePath, emails)
	if err != nil {
		// fmt.Println(err)
		// http.Error(w, "Error sending email", http.StatusInternalServerError)
		return err
	}
	return nil
}

// Send the file to the specified emails using SMTP
func sendFileToEmails(filePath string, emails []string) error {
	// Set up authentication information.
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer)

	// Create an email message
	m := email.NewMessage("Subject: File submission", "Please see a file attached bellow")
	m.To = emails

	// Add attachments
	if err := m.Attach(filePath); err != nil {
		return err
	}

	if err := email.Send("smtp.gmail.com:587", auth, m); err != nil {
		return err
	}

	return nil
}

// Check if the file MIME type is valid
func isValidFileType(fileHeader *multipart.FileHeader) bool {
	allowedTypes := map[string]bool{
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"application/pdf": true,
	}
	return allowedTypes[fileHeader.Header.Get("Content-Type")]
}
