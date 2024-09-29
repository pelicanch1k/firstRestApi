package psql

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/pelicanch1k/homework-http/pkg/logger/sl"
)

type Storage struct {
	Db  *sql.DB
	Log *slog.Logger
}

func errorsHandler(err error, log *slog.Logger) bool {
	const op = "internal.storage.psql"

	if err != nil {
		log.Error("failed to init storage",
			sl.Err(err),
			slog.Attr{
				Key:   "path",
				Value: slog.StringValue(op),
			},
		)

		return true
	}

	return false
}

func (r Storage) getId() int {
	var id int

	row := r.Db.QueryRow("select max(id) from books")
	err := row.Scan(&id)
	errorsHandler(err, r.Log)

	return id + 1
}

func (r Storage) getBookByName(name string) {}

func (r Storage) CreateBook(name, description string) (int64, error) {
	// добавление новой книги
	result, err := r.Db.Exec("insert into books (id, name, description) values ($1, $2, $3)",
		r.getId(),
		name,
		description)

	if answer := errorsHandler(err, r.Log); answer != false {
		return 0, err
	}

	num, err := result.RowsAffected()

	if answer := errorsHandler(err, r.Log); answer != false {
		return 0, err
	}

	return num, nil
}

func (r Storage) DeleteBook(id int) (int64, error) {
	result, err := r.Db.Exec("DELETE FROM books WHERE id = $1", id)
	if answer := errorsHandler(err, r.Log); answer != false {
		return 0, err
	}

	num, err := result.RowsAffected()

	if answer := errorsHandler(err, r.Log); answer != false {
		return 0, err
	}

	return num, nil
}

func (r Storage) updateBook(id int, mapValue map[string]string) (int64, error) {
	for key, value := range mapValue {
		fmt.Println(key, value)
		result, err := r.Db.Exec(fmt.Sprintf("UPDATE books SET %s = $1 WHERE id = $2", key), value, id)
		if err != nil {
			return 0, err
		}

		num, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}

		return num, nil
	}

	return 0, errors.New("No args")
}

func (r Storage) UpdateBook(id int, name, description string) (map[string]int64, error) {
	if id != 0 {
		answer := make(map[string]int64)

		if name != "" {
			answerCode, err := r.updateBook(id, map[string]string{"name": name})
			if err != nil {
				return nil, err
			}

			answer["name"] = answerCode
		}

		if description != "" {
			answerCode, err := r.updateBook(id, map[string]string{"description": description})
			if err != nil {
				return nil, err
			}

			answer["description"] = answerCode
		}

		return answer, nil
	} else {
		return nil, errors.New("Id is 0")
	}
}

func (r Storage) GetBook(id string) *sql.Row {
	// получение одной книги
	row := r.Db.QueryRow("select * from books where id = $1", id)

	return row
}

func (r Storage) GetAllBooks() (*sql.Rows, error) {
	// получение всех книг
	rows, err := r.Db.Query("select * from books")
	if err != nil {
		return nil, err
	}
	return rows, nil
}
