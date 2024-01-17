package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

var BaseURL = "https://query.getbible.net/v2/kjv/"

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

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
