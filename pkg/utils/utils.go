package main

import (
	"log"
	"math/rand"
	"time"
)

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
