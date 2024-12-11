package main

import (
	"go-crud/internal/db"
	"go-crud/internal/domain"
	"go-crud/internal/rate"
	"log"
	"net/http"
)

func main() {
	conn, err := db.NewNeonDatabase()
	if err != nil {
		log.Fatal(err)
	}

	bookService := domain.NewBookService(conn)
	http.HandleFunc("GET /books", middleware(bookService.HandleGetBooks))
	http.HandleFunc("GET /books/{id}", middleware(bookService.HandleGetBooksById))
	http.HandleFunc("POST /books/create", middleware(bookService.HandleCreateBook))
	http.HandleFunc("PUT /books/{id}", middleware(bookService.HandleUpdateBook))
	http.HandleFunc("DELETE /books/{id}", middleware(bookService.HandleDeleteBook))

	http.ListenAndServe(":8080", nil)
}

func middleware(next http.HandlerFunc) http.HandlerFunc {
	limiter := rate.NewRateLimiter(3, 1)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if limiter.IsRequestAllow() {
			next(w, r)
			return
		}

		domain.JSONError(w, "Too Many Requests", http.StatusTooManyRequests)
	}
}
