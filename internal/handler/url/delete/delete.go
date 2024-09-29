package delete

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
	helpers "github.com/pelicanch1k/homework-http/internal/handler/url"
	"github.com/pelicanch1k/homework-http/pkg/api/response"
)

const op = "internal.handler.delete.Delete"

type Response struct {
	response.Response
	Id string `json:"id,omitempty"`
}

type Deleter interface {
	DeleteBook(id int) (int64, error)
}

// @Summary      Delete a book by ID
// @Description  Delete a book by its ID
// @Tags         books
// @Produce      json
// @Param        id     path      int             true  "Book ID"
// @Success      200    {object}  response.Success
// @Failure      400    {object}  response.Error
// @Failure      404    {object}  response.Error
// @Router       /book/{id}/ [delete]
func Delete(log *slog.Logger, deleter Deleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("op", op),
		)

		id, err := helpers.IdToInt(w, r)
		if err != nil {
			return
		}

		log.Info("id is correct", slog.Any("id", id))

		deleted_id, err := deleter.DeleteBook(id)
		if err != nil {
			render.JSON(w, r, response.Error("failed to delete book"))

			return
		}

		if deleted_id == 0 {
			render.JSON(w, r, response.Error("the book has already been deleted"))

			return
		}

		log.Info("book deleted", slog.Any("id", id))

		responseOK(w, r, id)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, id int) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		Id:       strconv.Itoa(id),
	})
}
