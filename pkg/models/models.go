package main

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
