package helpers

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/render"
	"github.com/pelicanch1k/homework-http/pkg/api/response"
	"github.com/pelicanch1k/homework-http/pkg/logger/sl"
)

func IdToInt(w http.ResponseWriter, r *http.Request) (int, error) {
	path := strings.Split(r.URL.Path, "/")

	id, err := strconv.Atoi(path[2])
	if err != nil {
		render.JSON(w, r, response.Error("invalid value"))
		return 0, err
	}

	return id, nil
}

type Request struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func NewRequest(w http.ResponseWriter, r *http.Request, log *slog.Logger) (Request, error) {
	var req Request

	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		log.Error("failed to decode request body", sl.Err(err))

		render.JSON(w, r, response.Error("failed to decode request"))
	}

	return req, err
}
