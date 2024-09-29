package get

import (
	"database/sql"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
)

type Response struct {
	Id          int
	Name        string
	Description string
}

type Get interface {
	GetAllBooks() (*sql.Rows, error)
	GetBook(id string) *sql.Row
}

// @Summary      Get all books
// @Description  Get a list of books
// @Tags         books
// @Produce      json
// @Success      200  {array}   helpers.Book
// @Failure      500  {object}  response.Error
// @Router       /books/ [get]
func GetAll(log *slog.Logger, get Get) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.handler.get.GetAll"

		log = log.With(
			slog.String("op", op),
		)

		rows, err := get.GetAllBooks()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		defer rows.Close()
		var responseBooks []Response

		for rows.Next() {
			var books Response
			rows.Scan(&books.Id, &books.Name, &books.Description)
			log.Info(strconv.Itoa(books.Id),
				slog.Any("name", books.Name),
				slog.Any("description", books.Description),
			)

			responseBooks = append(responseBooks, books)
		}

		responseOK(w, r, responseBooks)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, responseBooks []Response) {
	w.WriteHeader(http.StatusOK)

	render.JSON(w, r, responseBooks)
}
