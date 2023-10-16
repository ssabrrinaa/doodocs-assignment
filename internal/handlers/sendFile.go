package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"os"
	"strings"

	"github.com/scorredoira/email"
)

func (h *ArchiveHandler) SendFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		errorHandler(w, r, http.StatusMethodNotAllowed, "")

		return
	}

	// Create the "uploads" directory if it doesn't exist
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to create 'uploads' directory", http.StatusInternalServerError)
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

	// Create a temporary file to store the uploaded file
	tempFile, err := os.CreateTemp("uploads", "upload-*.txt")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to create temporary file", http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	// Copy the file to the temporary file
	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, "Error copying file", http.StatusInternalServerError)
		return
	}

	// Send email with the file attachment
	err = sendFileToEmails(fileHeader.Filename, tempFile.Name(), emails)
	if err != nil {
		http.Error(w, "Error sending email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "File sent successfully")
}

const (
	smtpServer   = "smtp.gmail.com"
	smtpPort     = "587"
	smtpUsername = "bakirova200024@gmail.com"
	smtpPassword = "zevghlaxkgwpkdic"
)

// Send the file to the specified emails using SMTP
func sendFileToEmails(fileName, filePath string, emails []string) error {
	// Set up authentication information.
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer)

	// Create an email message.
	m := email.NewMessage("Subject", "Body")

	// Add attachments.
	if err := m.Attach(filePath); err != nil {
		panic(err)
	}

	// if err := m.AttachFile("path/to/file2.jpg"); err != nil {
	// 	panic(err)
	// }

	// Send the email.
	err := email.Send("smtp.example.com:587", auth, m)
	if err != nil {
		panic(err)
	}
	return nil
}

// Send the email
// err := smtp.SendMail(
// 	"smtp.gmail.com:587",
// 	auth,
// 	"bakirova200024@gmail.com",
// 	[]string{"bakirova200024@gmail.com"},
// 	// emails,
// 	buffer.Bytes(),
// )
// if err != nil {
// 	fmt.Println(err)
// }

// Check if the file MIME type is valid
func isValidFileType(fileHeader *multipart.FileHeader) bool {
	allowedTypes := map[string]bool{
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"application/pdf": true,
	}

	// mimeType := mime.TypeByExtension(filepath.Ext(fileHeader.Filename))

	return allowedTypes[fileHeader.Header.Get("Content-Type")]
}

// // Set up the MIME headers
// headers := make(map[string]string)
// headers["From"] = smtpUsername
// headers["To"] = strings.Join(emails, ", ")
// headers["Subject"] = "File Attachment"

// // Create the MIME message
// var message bytes.Buffer

// for key, value := range headers {
// 	message.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
// }

// // Attach the file
// message.WriteString("MIME-version: 1.0\r\n")
// message.WriteString("Content-Type: multipart/mixed; boundary=foo\r\n\r\n")

// message.WriteString("--foo\r\n")
// message.WriteString(fmt.Sprintf("Content-Type: text/plain\r\n\r\n"))
// message.WriteString("Please find the attached file.\r\n\r\n")

// message.WriteString("--foo\r\n")
// message.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", fileName))
// message.WriteString("Content-Type: application/octet-stream\r\n\r\n")

// // Read the file content
// fileContents, err := os.ReadFile(filePath)
// if err != nil {
// 	return err
// }

// message.Write(fileContents)
// message.WriteString("\r\n\r\n--foo--\r\n")

// // Connect to the server, authenticate, and send the email
// err = smtp.SendMail(fmt.Sprintf("%s:%s", smtpServer, smtpPort), auth, smtpUsername, emails, message.Bytes())
// if err != nil {
// 	return err
// }

// return nil
// }
