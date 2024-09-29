package get

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/pelicanch1k/homework-http/pkg/api/response"
)

// @Summary      Get a book by ID
// @Description  Get a book by its ID
// @Tags         books
// @Produce      json
// @Param        id     path      int             true  "Book ID"
// @Success      200    {object}  helpers.Book
// @Failure      400    {object}  response.Error
// @Failure      404    {object}  response.Error
// @Router       /book/{id}/ [get]
func GetBook(log *slog.Logger, get Get) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.handler.get.GetBook"

		log = log.With(
			slog.String("op", op),
		)

		id := strings.Split(r.URL.Path, "/")[2]
		if id == "" {
			render.JSON(w, r, response.Error("invalid id"))

			return
		}

		row := get.GetBook(id)
		var book Response

		row.Scan(&book.Id, &book.Name, &book.Description)

		log.Info(id,
			slog.Any("name", book.Name),
			slog.Any("description", book.Description),
		)

		responseOK(w, r, []Response{book})
	}
}
