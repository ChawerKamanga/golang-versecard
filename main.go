package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"os"

	"github.com/gin-gonic/gin"
)

type user struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Verse string `json:"verse"`
}

type VerseResponse struct {
	Data map[string]struct {
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
	} `json:"kjv_19_2"`
}

func testFun(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "test is done",
	})
}

func getBibleVerse(verse string) (string, error) {
	url := fmt.Sprintf("https://query.getbible.net/v2/kjv/%s", verse)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var verseResponse VerseResponse
	// Adjust this based on the actual API response
	err = json.Unmarshal(body, &verseResponse)
	if err != nil {
		return "", err
	}

	var verseText string
	for _, data := range verseResponse.Data {
		for _, verse := range data.Verses {
			verseText = verse.Text
			break // assuming you need only the first verse
		}
		break // assuming only one key in the map
	}

	return verseText, nil
}

func sendEmail(c *gin.Context) {
	var newUser user

	if err := c.BindJSON(&newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Error in binding JSON",
		})
		return
	}

	verseText, err := getBibleVerse(newUser.Verse)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch verse",
		})
		return
	}

	// Create an HTML email template with the verse
	emailTemplate := fmt.Sprintf(`
        <html>
            <body>
                <h1>Welcome %s!</h1>
                <p>Your verse: %s</p>
                <p>%s</p> <!-- Insert verse text here -->
            </body>
        </html>`, newUser.Name, newUser.Verse, verseText)

	fromEmail := os.Getenv("FROM_EMAIL")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	to := []string{newUser.Email}
	smtpHost := "sandbox.smtp.mailtrap.io"
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")

	message := []byte("To: " + newUser.Email + "\r\n" +
		"Subject: Welcome!\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n\r\n" +
		emailTemplate)

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, fromEmail, to, message)
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to send email",
		})
		return
	}

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
