package main

import (
	"log"
	"mime/multipart"
	"net/smtp"
)

const (
	BOUNDARY = "boundary875923523489756345"
	HOST     = "alexcoder04.de"
)

type MailConfig struct {
	Password string

	FromAddress string
	FromName    string

	ToAddress string

	Title string
	Body  string

	Attachments []*multipart.FileHeader
}

func SendMail(config MailConfig) error {
	auth := smtp.PlainAuth("", "noreply@"+HOST, config.Password, "mail."+HOST)

	message := "From: " + config.FromName + "<noreply@" + HOST + ">\r\n" +
		"To: " + config.ToAddress + "\r\n" +
		"Reply-To: " + config.FromName + "<" + config.FromAddress + ">\r\n" +
		"Subject: " + config.Title + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: multipart/mixed; boundary=" + BOUNDARY + "\r\n\r\n" +
		"--" + BOUNDARY + "\r\n" +
		"Content-Type: text/plain\r\n\r\n" +
		config.Body + "\r\n\r\n"

	log.Println("built message")

	for _, attachment := range config.Attachments {
		file, err := attachment.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		message += "--" + BOUNDARY + "\r\n" +
			"Content-Type: application/octet-stream\r\n" +
			"Content-Disposition: attachment; filename=" + attachment.Filename + "\r\n\r\n"

		buffer := make([]byte, attachment.Size)
		file.Read(buffer)
		message += string(buffer) + "\r\n\r\n"
	}

	message += "--" + BOUNDARY + "--"

	return smtp.SendMail("mail."+HOST+":587", auth, "noreply@"+HOST, []string{config.ToAddress}, []byte(message))
}
