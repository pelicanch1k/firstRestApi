package main

import (
	"log"
	"net/http"
)

const PORT = ":8080"

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /book", loggerMiddleware(getBook))
	mux.HandleFunc("GET /book/delete", loggerMiddleware(deleteBook))
	mux.HandleFunc("GET /book/create", loggerMiddleware(createBook))
	mux.HandleFunc("/book/update/", loggerMiddleware(updateBook))

	mux.HandleFunc("/books", loggerMiddleware(getBooks))

	go func() {
		log.Println("server is running...", "\n")
	}()

	if err := http.ListenAndServe(PORT, mux); err != nil {
		panic(err)
	}
}

func loggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] - %s\n", r.Method, r.URL)
		next(w, r)
	}
}
