package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/cors"
	"gorm.io/gorm/clause"
)

// Author Struct
type Author struct {
	gorm.Model
	Firstname string
	Lastname  string
	Books     []Book
}

// Book Struct
type Book struct {
	gorm.Model
	Isbn     string
	Title    string
	AuthorID uint
}

var db *gorm.DB

var err error

// Init Books varaiable as slice Book struct and some mock data
var (
	authors = []Author{
		{Firstname: "John", Lastname: "Doe"},
		{Firstname: "Steve", Lastname: "Smith"},
	}
	books = []Book{
		{Isbn: "448743", Title: "Book one", AuthorID: 1},
		{Isbn: "875468", Title: "Book two", AuthorID: 2},
	}
)

// Get All Authors and if any books
func getAuthorsExtra(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	db.Find(&authors)
	db.Preload("Books").Preload(clause.Associations).Find(&authors)
	json.NewEncoder(w).Encode(&authors)
	log.Printf("Get all books issued")

}

// Get All Authors
func getAllAuthors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	db.Find(&authors)
	json.NewEncoder(w).Encode(&authors)
}

// Get Single book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) //get params
	var book Book

	if err := db.First(&book, params["id"]).Error; err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	db.First(&book, params["id"])
	json.NewEncoder(w).Encode(&book)
}

// Get all books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db.Find(&books)
	json.NewEncoder(w).Encode(&books)
}

// Delete Single book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) //get params
	var book Book

	if err := db.First(&book, params["id"]).Error; err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	db.Delete(&book)
	var books []Book
	db.Find(&books)
	json.NewEncoder(w).Encode(&books)
}

//Create an author
func createAuthor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var author Author
	_ = json.NewDecoder(r.Body).Decode(&author)
	params := r.URL.Query()
	firstname := params.Get("firstname")
	lastname := params.Get("lastname")
	if len(firstname) == 0 || len(lastname) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	db.Create(&Author{Firstname: firstname, Lastname: lastname})
	db.Find(&authors)
	json.NewEncoder(w).Encode(&authors)
}

//Create a book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	params := r.URL.Query()
	isbn := params.Get("isbn")
	title := params.Get("title")
	if len(isbn) == 0 || len(title) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	db.Create(&Book{Isbn: isbn, Title: title})
	db.Find(&books)
	json.NewEncoder(w).Encode(&books)
}

//InitDB initialization
func InitDB() {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres sslmode=disable password=docker")
	if err != nil {
		panic("failed to connect database")
	}

}

func main() {
	log.SetPrefix("LOG: ")
	log.SetFlags(log.Ldate | log.Ltime)
	log.Println("Server Started...")
	isReady := &atomic.Value{}
	isReady.Store(false)
	go func() {
		log.Printf("Readyz probe is negative by default...")
		time.Sleep(10 * time.Second)
		isReady.Store(true)
		log.Printf("Readyz probe is positive.")
	}()
	router := mux.NewRouter()
	InitDB()
	defer db.Close()
	db.DropTableIfExists(&Book{})
	db.DropTableIfExists(&Author{})

	db.AutoMigrate(&Book{})

	db.AutoMigrate(&Author{})

	for index := range books {
		db.Create(&books[index])
	}

	for index := range authors {
		db.Create(&authors[index])
	}

	// Route Handlers / Endpoints
	router.HandleFunc("/api/authorsextra", getAuthorsExtra).Methods("GET")
	router.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/api/books", getBooks).Methods("GET")
	router.HandleFunc("/api/authors", getAllAuthors).Methods("GET")
	router.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	router.HandleFunc("/api/authors/", createAuthor).Methods("POST")
	router.HandleFunc("/api/books/", createBook).Methods("POST")

	router.HandleFunc("/healthz", healthz)
	router.HandleFunc("/readyz", readyz(isReady))

	log.Print("Server Started function...")
	handler := cors.Default().Handler(router)

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", handler))

}
