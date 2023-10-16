package main

import (
	"fmt"
	"net/smtp"
)

func main() {
	sendMail()
}

func sendMail() {
	auth := smtp.PlainAuth(
		"",
		"bakirova200024@gmail.com",
		"zevghlaxkgwpkdic",
		"smtp.gmail.com",
	)

	msg := "Subject: Send a file\nThis is the body of the email"

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"bakirova200024@gmail.com",
		[]string{"bakirova200024@gmail.com"},
		[]byte(msg),
	)
	if err != nil {
		fmt.Println(err)
	}
}
