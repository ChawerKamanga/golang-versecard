package main

import (
	"math/rand"
	"time"
)

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

func generateEmailTemplate(newUser user, verseName string, verseText string) string {
	randomImageUrl := getRandomImageUrl()

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

func getRandomImageUrl() string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	return imageUrls[r.Intn(len(imageUrls))]
}
