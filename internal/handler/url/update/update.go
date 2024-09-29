package update

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	helpers "github.com/pelicanch1k/homework-http/internal/handler/url"
	"github.com/pelicanch1k/homework-http/pkg/api/response"
)

const op = "internal.handler.update.Update"

type Response struct {
	response.Response
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Updater interface {
	UpdateBook(id int, name, description string) (map[string]int64, error)
}

// @Summary      Update a book by ID
// @Description  Update a book by its ID with the input payload
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id     path      int             true  "Book ID"
// @Param        book   body      helpers.Book    true  "Book"
// @Success      200    {object}  helpers.Book
// @Failure      400    {object}  response.Error
// @Failure      404    {object}  response.Error
// @Router       /book/{id}/ [put]
func Update(log *slog.Logger, updater Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("op", op),
		)

		// get id from url
		id, err := helpers.IdToInt(w, r)
		if err != nil {
			render.JSON(w, r, response.Error("invalid id"))

			return
		}

		log.Info("id is correct", slog.Any("id", id))

		// get json
		req, err := helpers.NewRequest(w, r, log)
		if err != nil {
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		answer, err := updater.UpdateBook(id, req.Name, req.Description)
		if err != nil {
			log.Error(err.Error())
			render.JSON(w, r, response.Error("failed to update book"))

			return
		}

		args := make([]slog.Attr, 2)
		for key, value := range answer {
			args = append(args, slog.Any(key, value))
		}

		log.Info("book updated", args[0], args[1])

		responseOK(w, r, req)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, req helpers.Request) {
	render.JSON(w, r, Response{
		Response:    response.OK(),
		Name:        req.Name,
		Description: req.Description,
	})
}
