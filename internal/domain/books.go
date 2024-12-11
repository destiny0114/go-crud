package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
)

type Book struct {
	ID          int32  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type BookService struct {
	db *pgx.Conn
}

type APIResponse[T any] struct {
	Data   T   `json:"data"`
	Status int `json:"status"`
}

type APIError struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	APIError `json:"error"`
	Status   int `json:"status"`
}

var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

func NewBookService(db *pgx.Conn) *BookService {
	return &BookService{
		db: db,
	}
}

func (s *BookService) HandleGetBooks(w http.ResponseWriter, req *http.Request) {
	rows, err := s.db.Query(context.Background(), "select * from books")
	if err != nil {
		fmt.Printf("Query error: %v", err)
		JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	books, err := pgx.CollectRows(rows, pgx.RowToStructByName[Book])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := JSONResponse(w, http.StatusOK, APIResponse[[]Book]{
		Data:   books,
		Status: http.StatusOK,
	}); err != nil {
		JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *BookService) HandleGetBooksById(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	if len(id) == 0 {
		JSONError(w, "ID path not exist", http.StatusInternalServerError)
		return
	}
	rows, err := s.db.Query(context.Background(), "select * from books where id = $1", id)
	if err != nil {
		fmt.Printf("Query error: %v", err)
		JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	books, err := pgx.CollectRows(rows, pgx.RowToStructByName[Book])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := JSONResponse(w, http.StatusOK, APIResponse[[]Book]{
		Data:   books,
		Status: http.StatusOK,
	}); err != nil {
		JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *BookService) HandleCreateBook(w http.ResponseWriter, req *http.Request) {
	if contentType := req.Header.Get("Content-Type"); contentType != "application/json" {
		JSONError(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}

	type Input struct {
		Title       string `json:"title" validate:"required"`
		Description string `json:"description" validate:"required"`
	}

	var book Input

	if err := json.NewDecoder(req.Body).Decode(&book); err != nil {
		JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := validate.Struct(book); err != nil {
		JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := s.db.Exec(context.Background(), "insert into books (title, description) values ($1, $2)", book.Title, book.Description); err != nil {
		fmt.Printf("Insert error: %v", err)
		JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := JSONResponse(w, http.StatusOK, APIResponse[Input]{
		Data:   book,
		Status: http.StatusOK,
	}); err != nil {
		JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *BookService) HandleUpdateBook(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	if len(id) == 0 {
		JSONError(w, "ID path not exist", http.StatusInternalServerError)
		return
	}

	if contentType := req.Header.Get("Content-Type"); contentType != "application/json" {
		JSONError(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}

	type Input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	var book Input

	if err := json.NewDecoder(req.Body).Decode(&book); err != nil {
		JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := validate.Struct(book); err != nil {
		JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := s.db.Exec(context.Background(),
		"UPDATE books SET title = COALESCE(NULLIF($1, ''), title), description = COALESCE(NULLIF($2, ''), description) WHERE id = $3",
		book.Title, book.Description, id); err != nil {
		fmt.Printf("Update error: %v", err)
		JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := JSONResponse(w, http.StatusOK, APIResponse[Input]{
		Data:   book,
		Status: http.StatusOK,
	}); err != nil {
		JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *BookService) HandleDeleteBook(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	if len(id) == 0 {
		JSONError(w, "ID path not exist", http.StatusInternalServerError)
		return
	}
	result, err := s.db.Exec(context.Background(), "delete from books where id = $1", id)
	if err != nil {
		fmt.Printf("Delete error: %v", err)
		JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if result.RowsAffected() != 1 {
		JSONError(w, "No row found to delete", http.StatusInternalServerError)
		return
	}

	if err := JSONResponse(w, http.StatusOK, APIResponse[map[string]interface{}]{
		Data: map[string]interface{}{
			"message": fmt.Sprintf("Deleted with id %v", req.PathValue("id")),
		},
		Status: http.StatusOK,
	}); err != nil {
		JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func JSONResponse[T any](w http.ResponseWriter, status int, v T) error {
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func JSONError(w http.ResponseWriter, error string, status int) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(ErrorResponse{APIError: APIError{Message: error}, Status: status}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
