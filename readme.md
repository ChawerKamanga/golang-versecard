# Go VerseCardAPI

Hey there, so this project is all about letting people request a personalized email with their favorite Bible verse in it. You give it your name, email, and the specific verse you want, and it generates a fancy email with the verse and a cool background image. It uses [GetBible](https://query.getbible.net/) to make sure the verses are accurate and [Pexels](https://www.pexels.com/) to find pretty pictures.

## Getting Started

These instructions will guide you through getting a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Go (Version 1.21.6)
- Gin (Version 1.9.1)
- GoDotEnv

### Installing

1. **Clone the Repository**

   ```bash
   $ git clone https://github.com/ChawerKamanga/golang-versecard

2. **Enter the directory**

   ```bash
   $ cd golang-versecard

3. **Install Gin**

   ```bash
   $ go get -u github.com/gin-gonic/gin

4. **Install DotEnv**

   ```bash
   $ go get github.com/joho/godotenv

5. **Run the code**

   ```bash
    $ go run main.go

6. **Test the API Using Curl**

   ```bash
    $ curl localhost:8080/user --include --header "Content-Type: application/json" -d @body.json --request "POST"

## License

This project is licensed under the MIT License
