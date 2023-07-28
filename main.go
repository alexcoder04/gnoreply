package main

import (
	"log"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PostData struct {
	Token       string                  `form:"token" binding:"required"`
	Title       string                  `form:"title" binding:"required"`
	Recepient   string                  `form:"recipient" binding:"required"`
	Body        string                  `form:"body" binding:"required"`
	Attachments []*multipart.FileHeader `form:"attachments"`
}

func Respond(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"message": message,
	})
}

func main() {
	config := LoadConfig()

	r := gin.New()

	r.POST("/send", func(c *gin.Context) {
		var data PostData

		// max 5MB
		if err := c.Request.ParseMultipartForm(5 * 1024 * 1024); err != nil {
			Respond(c, http.StatusBadRequest, "Unable to parse form data")
			return
		}

		// bind fields
		if err := c.ShouldBind(&data); err != nil {
			Respond(c, http.StatusBadRequest, "Fields provided are invalid")
			return
		}

		// check token
		for _, user := range config.Users {
			if user.Token == data.Token {
				err := SendMail(MailConfig{
					Password:    config.Password,
					FromAddress: user.Address,
					FromName:    user.Name,
					ToAddress:   data.Recepient,
					Title:       data.Title,
					Body:        data.Body,
					Attachments: data.Attachments,
				})
				if err != nil {
					log.Println(err.Error())
					Respond(c, http.StatusInternalServerError, "Failed to send mail")
					return
				}

				Respond(c, http.StatusOK, "Mail sent successfully")
				return
			}
		}

		Respond(c, http.StatusForbidden, "Invalid token")
	})

	r.Run(":" + config.Port)
}
