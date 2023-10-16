package handlers

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"strings"
)

func (h *ArchiveHandler) SendFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		errorHandler(w, r, http.StatusMethodNotAllowed, "")

		return
	}

	// Parse the incoming multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get the file and file header from the form
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get the list of emails from the form
	emailsStr := r.FormValue("emails")
	emails := strings.Split(emailsStr, ",")

	// Check if the file MIME type is valid
	if !isValidFileType(fileHeader) {
		http.Error(w, "Invalid file type", http.StatusBadRequest)
		return
	}

	// Send the file to the specified emails
	err = sendFileToEmails(file, fileHeader.Filename, emails)
	if err != nil {
		http.Error(w, "Failed to send file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "File sent successfully")
}

// Check if the file MIME type is valid
func isValidFileType(fileHeader *multipart.FileHeader) bool {
	allowedTypes := map[string]bool{
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"application/pdf": true,
	}

	// mimeType := mime.TypeByExtension(filepath.Ext(fileHeader.Filename))

	return allowedTypes[fileHeader.Header.Get("Content-Type")]
}

// Send the file to the specified emails using SMTP
func sendFileToEmails(file io.Reader, filename string, emails []string) error {
	// Replace these with your actual email credentials
	from := "bakirova200024@gmail.com"
	password := "54321bsb"
	smtpServer := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpServer)

	// Create the email message
	// message := createMessage(from, emails, "File Subject", "File Body")

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add the file part to the message
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	// Add the emails part to the message
	writer.WriteField("emails", strings.Join(emails, ","))

	writer.Close()

	// Send the email with the file attachment
	err = smtp.SendMail(smtpServer+":"+smtpPort, auth, from, emails, body.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// Create the email message
func createMessage(from string, to []string, subject string, body string) string {
	message := "From: " + from + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Subject: " + subject + "\n\n" +
		body
	return message
}
