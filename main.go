package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/smtp"

	"github.com/gin-gonic/gin"
)

type user struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Verse string `json:"verse"`
}

type VerseResponse map[string]struct {
	Translation  string   `json:"translation"`
	Abbreviation string   `json:"abbreviation"`
	Lang         string   `json:"lang"`
	Language     string   `json:"language"`
	Direction    string   `json:"direction"`
	Encoding     string   `json:"encoding"`
	BookNr       int      `json:"book_nr"`
	BookName     string   `json:"book_name"`
	Chapter      int      `json:"chapter"`
	Name         string   `json:"name"`
	Ref          []string `json:"ref"`
	Verses       []struct {
		Chapter int    `json:"chapter"`
		Verse   int    `json:"verse"`
		Name    string `json:"name"`
		Text    string `json:"text"`
	} `json:"verses"`
}

func testFun(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "test is done",
	})
}

func getBibleVerse(verse string) (string, string, error) {
	url := fmt.Sprintf("https://query.getbible.net/v2/kjv/%s", verse)

	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var verseResponse VerseResponse
	err = json.Unmarshal(body, &verseResponse)
	if err != nil {
		return "", "", err
	}

	var verseName, verseText string
	for _, v := range verseResponse {
		for _, verse := range v.Verses {
			verseName = verse.Name
			verseText = verse.Text
			break // assuming you need only the first verse
		}
		break // break after handling the first key-value pair
	}

	return verseName, verseText, nil
}

func sendEmail(c *gin.Context) {
	var newUser user

	if err := c.BindJSON(&newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Error in binding JSON",
		})
		return
	}

	verseName, verseText, err := getBibleVerse(newUser.Verse)

	fmt.Println(verseName)
	fmt.Println(verseText)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch verse",
		})
		return
	}

	// Create an HTML email template with the verse
	from := "a135b9a056bcc2@sandbox.smtp.mailtrap.io" // Using the username as the from address
	password := "4d04cd8c684981"
	to := []string{newUser.Email}
	smtpHost := "sandbox.smtp.mailtrap.io"
	smtpPort := "2525" // You can also try "25", "465", or "587" if "2525" doesn't work

	// Email message
	message := []byte("To: " + newUser.Email + "\r\n" +
		"Subject: Welcome!\r\n" +
		"\r\n" +
		"Hello " + newUser.Name + ", here is your verse. \n" + verseName + "\n" + verseText)

	// Authenticate with the SMTP server
	auth := smtp.PlainAuth("", "a135b9a056bcc2", password, smtpHost) // Use the provided username for authentication

	// Send the email
	newerr := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if newerr != nil {
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
