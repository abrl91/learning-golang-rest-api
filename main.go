package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}

func getBooks(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, books)
}

// return a pointer to a book or an error
func getBookById(id string) (*book, error) {
	for i, book := range books {
		if book.ID == id {
			return &books[i], nil
		}
	}
	return nil, errors.New("Book not found")
}

func getBook(context *gin.Context) {
	id := context.Param("id")
	book, err := getBookById(id)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}
	context.IndentedJSON(http.StatusOK, book)
}

func createBook(context *gin.Context) {
	var newBook book
	if err := context.ShouldBindJSON(&newBook); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	books = append(books, newBook)
	context.IndentedJSON(http.StatusCreated, newBook)
}

func updateBook(context *gin.Context) {
	id := context.Param("id")
	book, err := getBookById(id)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}
	if err := context.ShouldBindJSON(book); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.IndentedJSON(http.StatusOK, book)
}

func deleteBook(context *gin.Context) {
	id := context.Param("id")
	for i, book := range books {
		if book.ID == id {
			books = append(books[:i], books[i+1:]...)
			break
		}
	}
	context.Status(http.StatusNoContent)
}

func checkoutBook(context *gin.Context) {
	id := context.Param("id")
	book, err := getBookById(id)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}
	if book.Quantity == 0 {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Book not available"})
		return
	}
	book.Quantity--
	context.IndentedJSON(http.StatusOK, book)
}

func main() {
	router := gin.Default()

	//curl localhost:8080/books
	router.GET("/books", getBooks)
	// curl localhost:8080/books -i -H Content-Type:\ application/json -d @body.json -X POST
	router.POST("/books", createBook)
	// curl localhost:8080/books/1 -i
	router.PUT("/books/:id", updateBook)
	// curl localhost:8080/books/4 -i -X DELETE
	router.DELETE("/books/:id", deleteBook)
	// curl localhost:8080/books/1/checkout -i
	router.POST("/books/:id/checkout", checkoutBook)

	router.Run("localhost:8080")
}
