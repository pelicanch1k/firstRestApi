package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func errorsHandler(err error) {
	if err != nil {
		panic(err)
	}
}

type Repository struct {
	User     string
	Password string
	Dbname   string
	db       *sql.DB
}

func NewRepository(user, password, dbname string) *Repository {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	errorsHandler(err)

	return &Repository{user, password, dbname, db}
}

func (r Repository) getId() int {
	var id int
	row := r.db.QueryRow("select max(id) from books")
	err := row.Scan(&id)
	errorsHandler(err)

	return id + 1

}

func (r Repository) CreateBook(name, description string) int64 {
	// добавление новой книги
	result, err := r.db.Exec("insert into books (id, name, description) values ($1, $2, $3)",
		r.getId(),
		name,
		description)

	errorsHandler(err)

	num, err := result.RowsAffected()
	errorsHandler(err)
	return num
}

func (r Repository) DeleteBook(id int) int64 {
	result, err := r.db.Exec("DELETE FROM books WHERE id = $1", id)
	errorsHandler(err)

	num, err := result.RowsAffected()
	errorsHandler(err)
	return num
}

func (r Repository) updateBook(id int, mapValue map[string]string) int64 {
	for key, value := range mapValue {
		fmt.Println(key, value)
		result, err := r.db.Exec(fmt.Sprintf("UPDATE books SET %s = $1 WHERE id = $2", key), value, id)
		errorsHandler(err)

		num, err := result.RowsAffected()
		errorsHandler(err)
		return num
	}

	return 0
}

func (r Repository) UpdateBook(id int, name, description string) {
	if id != 0 {
		answer := make([]int64, 1)

		if name != "" {
			answer = append(answer, r.updateBook(id, map[string]string{"name": name}))
		}

		if description != "" {
			answer = append(answer, r.updateBook(id, map[string]string{"description": description}))
		}
	}
}

func (r Repository) GetBook(id string) *sql.Row {
	// получение одной книги
	row := r.db.QueryRow("select * from books where id = $1", id)

	return row
}

func (r Repository) GetAllBooks() *sql.Rows {
	// получение всех книг
	rows, err := r.db.Query("select * from books")
	if err != nil {
		log.Fatal(err)
	}

	return rows
}
