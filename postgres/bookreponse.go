package main

// BookResponse Struct
type BookResponse struct {
	ID     string          `json:"id"`
	Isbn   string          `json:"isbn"`
	Title  string          `json:"title"`
	Author *AuthorResponse `json:"author"`
}

// AuthorResponse Struct
type AuthorResponse struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}
