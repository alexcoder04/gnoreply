package main

import (
	"bytes"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
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
	var encoded bytes.Buffer

	encoder := quotedprintable.NewWriter(&encoded)

	_, err := encoder.Write([]byte(config.Body))
	if err != nil {
		return err
	}

	encoder.Close()

	message := "From: " + config.FromName + "<noreply@" + HOST + ">\r\n" +
		"To: " + config.ToAddress + "\r\n" +
		"Reply-To: " + config.FromName + "<" + config.FromAddress + ">\r\n" +

		"Subject: " + mime.QEncoding.Encode("utf-8", config.Title) + "\r\n" +

		"MIME-Version: 1.0\r\n" +
		"Content-Type: multipart/mixed; boundary=\"" + BOUNDARY + "\"\r\n\r\n" +
		"--" + BOUNDARY + "\r\n" +
		"Content-Type: text/plain; charset=utf-8\r\n" +
		"Content-Transfer-Encoding: quoted-printable\r\n\r\n" +
		encoded.String() + "\r\n\r\n"

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

	auth := smtp.PlainAuth("", "noreply@"+HOST, config.Password, "mail."+HOST)
	return smtp.SendMail("mail."+HOST+":587", auth, "noreply@"+HOST, []string{config.ToAddress}, []byte(message))
}
