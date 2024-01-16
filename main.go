package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var BaseURL = "https://query.getbible.net/v2/kjv/"

var imageUrls = []string{
	"https://images.pexels.com/photos/3225517/pexels-photo-3225517.jpeg",
	"https://images.pexels.com/photos/459038/pexels-photo-459038.jpeg",
	"https://images.pexels.com/photos/3571576/pexels-photo-3571576.jpeg",
	"https://images.pexels.com/photos/1671230/pexels-photo-1671230.jpeg",
	"https://images.pexels.com/photos/573863/pexels-photo-573863.jpeg",
	"https://images.pexels.com/photos/3244513/pexels-photo-3244513.jpeg",
	"https://images.pexels.com/photos/1547992/pexels-photo-1547992.jpeg",
	"https://images.pexels.com/photos/624015/pexels-photo-624015.jpeg",
	"https://images.pexels.com/photos/3225517/pexels-photo-3225517.jpeg",
}

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

func getRandomImageUrl() string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	return imageUrls[r.Intn(len(imageUrls))]
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func respondWithError(c *gin.Context, statusCode int, message string) {
	c.IndentedJSON(statusCode, gin.H{"message": message})
}

func getBibleVerse(verse string) (string, string, error) {
	url := fmt.Sprintf(BaseURL+"%s", verse)

	resp, err := http.Get(url)
	checkError(err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	checkError(err)

	var verseResponse VerseResponse
	err = json.Unmarshal(body, &verseResponse)
	checkError(err)

	var verseName, verseText string
	for _, v := range verseResponse {
		verseName = v.Verses[0].Name
		verseText = v.Verses[0].Text
		break
	}

	return verseName, verseText, nil
}

func generateEmailTemplate(newUser user, verseName string, verseText string) string {
	randomImageUrl := getRandomImageUrl()

	fmt.Println(randomImageUrl)
	return `
    <!DOCTYPE html>
    <html>
    <head>
    <style>
        * body {
            font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
        }

        .email-container {
            background-image: url('` + randomImageUrl + `');
            background-size: cover;
            background-position: center;
            padding: 10px;
            text-align: center;
            color: #ffffff;

            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
        }

        .email-content {
            background: rgba(0, 0, 0, 0.7);
            padding: 20px;
        }

        h1 {
            color: #fff;
        }

        p {
            color: #fff;
        }
    </style>
    </head>
    <body>
    <div class="email-container">
        <div class="email-content">
            <p>Hello ` + newUser.Name + `</p>
            <h1>Here is Your Verse!</h1>
            <p><strong>` + verseName + `</strong></p>
            <p>` + verseText + `</p>
        </div>
    </div>
    </body>
    </html>
    `
}

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

func main() {
	router := gin.Default()
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	router.POST("/user", sendEmail)

	router.Run()
}
