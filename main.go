package main

import (
	"fmt"
	"net/http"
	"net/smtp"

	"github.com/gin-gonic/gin"
)

type user struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func testFun(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "test is done",
	})
}
func sendEmail(c *gin.Context) {
	var newUser user

	if err := c.BindJSON(&newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Error in binding JSON",
		})
		return
	}

	// Setup the email details
	from := "a135b9a056bcc2@sandbox.smtp.mailtrap.io" // Using the username as the from address
	password := "4d04cd8c684981"
	to := []string{newUser.Email}
	smtpHost := "sandbox.smtp.mailtrap.io"
	smtpPort := "2525" // You can also try "25", "465", or "587" if "2525" doesn't work

	// Email message
	message := []byte("To: " + newUser.Email + "\r\n" +
		"Subject: Welcome!\r\n" +
		"\r\n" +
		"Hello " + newUser.Name + ", welcome to our service.")

	// Authenticate with the SMTP server
	auth := smtp.PlainAuth("", "a135b9a056bcc2", password, smtpHost) // Use the provided username for authentication

	// Send the email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to send email",
		})
		return
	}

	// Confirm success
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "Email sent successfully",
	})
}

func main() {
	router := gin.Default()

	router.GET("/test", testFun)
	router.POST("/user", sendEmail)

	router.Run()
}
