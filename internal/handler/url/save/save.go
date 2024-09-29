package save

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	helpers "github.com/pelicanch1k/homework-http/internal/handler/url"
	"github.com/pelicanch1k/homework-http/pkg/api/response"
)

const op = "internal.handler.save.New"

type Response struct {
	response.Response
	Name string `json:"name,omitempty"`
}

type Saver interface {
	CreateBook(name, description string) (int64, error)
}

// @Summary      Create a new book
// @Description  Create a new book with the input payload
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        book  body      helpers.Book  true  "Book"
// @Success      200   {object}  helpers.Book
// @Failure      400   {object}  response.Error
// @Failure      500   {object}  response.Error
// @Router       /books/ [post]
func New(log *slog.Logger, saver Saver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("op", op),
		)

		req, err := helpers.NewRequest(w, r, log)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			render.JSON(w, r, response.Error("invalid id"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		id, err := saver.CreateBook(req.Name, req.Description)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			render.JSON(w, r, response.Error("failed to add book"))

			return
		}

		log.Info("book added", slog.Int64("id", id))

		responseOK(w, r, req.Name)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, name string) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		Name:     name,
	})
}
