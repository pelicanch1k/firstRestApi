package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/pelicanch1k/homework-http/server/pkg/database"
)

type Book struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var db = database.NewRepository("postgres", "root", "postgres")

func getBook(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("x-id")
	book := Book{}

	dbBook := db.GetBook(id)

	err := dbBook.Scan(&book.Id, &book.Name, &book.Description)
	if err != nil {
		panic(err)
	}
	fmt.Println(book.Id, book.Name, book.Description)

	resp, err := json.Marshal(book)
	if err != nil {
		panic(nil)
	}

	w.Write(resp)
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	var books []Book
	dbBooks := db.GetAllBooks()

	for dbBooks.Next() {
		book := Book{}

		err := dbBooks.Scan(&book.Id, &book.Name, &book.Description)
		if err != nil {
			panic(err)
		}

		books = append(books, book)
	}

	resp, err := json.Marshal(books)
	if err != nil {
		panic(nil)
	}

	w.Write(resp)
}

func getAnswer(n int64) string {
	var answer string

	if n == int64(1) {
		answer = "ok"
	} else {
		answer = "not ok"
	}

	println(answer)
	return answer
}

func createBook(w http.ResponseWriter, r *http.Request) {
	name := r.Header.Get("name")
	desc := r.Header.Get("desc")
	println(name, desc)

	dbBook := db.CreateBook(name, desc)
	answer := getAnswer(dbBook)

	fmt.Fprintf(w, answer)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.Header.Get("id"))
	if err != nil {
		panic(err)
	}

	dbBook := db.DeleteBook(id)
	answer := getAnswer(dbBook)

	fmt.Fprintf(w, answer)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	// Декодируем JSON из тела запроса
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Received book: %+v\n", book)

	// Здесь можно добавить логику обработки книги (например, сохранить в базу данных)
	db.UpdateBook(book.Id, book.Name, book.Description)

	// Отправляем успешный ответ
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book) // Возвращаем созданную книгу в ответе
}
