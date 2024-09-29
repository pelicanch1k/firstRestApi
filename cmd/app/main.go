package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/pelicanch1k/homework-http/docs"
	httpSwagger "github.com/swaggo/http-swagger" // http-swagger middleware

	"github.com/pelicanch1k/homework-http/internal/config"
	"github.com/pelicanch1k/homework-http/internal/handler/url/delete"
	"github.com/pelicanch1k/homework-http/internal/handler/url/get"
	"github.com/pelicanch1k/homework-http/internal/handler/url/save"
	"github.com/pelicanch1k/homework-http/internal/handler/url/update"

	"github.com/pelicanch1k/homework-http/pkg/database"
	"github.com/pelicanch1k/homework-http/pkg/logger/handlers/slogpretty"
	"github.com/pelicanch1k/homework-http/pkg/logger/sl"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.

// @host localhost:8080
// @BasePath /

func main() {
	// init config
	cfg := config.MustLoad()

	// init logger
	log := setupLogger(cfg.Env)

	log.Info("starting")
	log.Debug("debug")

	// init storage
	storage, err := database.NewPostgresConnection(
		database.NewConnectionInfo(cfg), log)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage

	// init router
	router := mux.NewRouter()

	router.HandleFunc("/swagger/{*}", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	)).Methods(http.MethodGet)

	router.HandleFunc("/books/", loggerMiddleware(save.New(log, storage))).Methods(http.MethodPost)
	router.HandleFunc("/books/", loggerMiddleware(get.GetAll(log, storage))).Methods(http.MethodGet)

	router.HandleFunc("/book/{id:[0-9]+}/", loggerMiddleware(get.GetBook(log, storage))).Methods(http.MethodGet)
	router.HandleFunc("/book/{id:[0-9]+}/", loggerMiddleware(delete.Delete(log, storage))).Methods(http.MethodDelete)
	router.HandleFunc("/book/{id:[0-9]+}/", loggerMiddleware(update.Update(log, storage))).Methods(http.MethodPut)

	// run server
	srv := &http.Server{
		Addr:         cfg.Http_server.Address,
		Handler:      router,
		ReadTimeout:  cfg.Http_server.Timeout,
		WriteTimeout: cfg.Http_server.Timeout,
		IdleTimeout:  cfg.Http_server.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
		// log = slog.New(
		// 	slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		// )
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}

func loggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] - %s\n", r.Method, r.URL)
		next(w, r)
	}
}
