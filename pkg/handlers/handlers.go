package main

import (
	"net/http"
	"net/smtp"
	"os"

	"github.com/gin-gonic/gin"
)

func sendEmail(c *gin.Context) {
	var newUser user

	if err := c.BindJSON(&newUser); err != nil {
		respondWithError(c, http.StatusBadRequest, "Error in binding JSON")
		return
	}

	verseName, verseText, err := getBibleVerse(newUser.Verse)

	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to fetch verse")
		return
	}

	from := os.Getenv("FROM_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	to := []string{newUser.Email}
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")

	message := generateEmailTemplate(newUser, verseName, verseText)

	msg := []byte("To: " + newUser.Email + "\r\n" +
		"Subject: Your Verse!\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n\r\n" +
		message)

	auth := smtp.PlainAuth("", smtpUsername, password, smtpHost)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to send email")
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "Email sent successfully",
	})
}
