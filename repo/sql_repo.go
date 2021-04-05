package repo

import (
	"context"
	"database/sql"

	"github.com/kil-san/simple-todo-server/model"
	_ "github.com/mattn/go-sqlite3"
)

type SqlRepo struct {
	db *sql.DB
}

func NewSqlRepo(db *sql.DB) SqlRepo {
	return SqlRepo{
		db: db,
	}
}

func (r SqlRepo) Create(ctx context.Context, data model.Todo) error {
	tx, _ := r.db.Begin()
	stmt, _ := tx.Prepare("INSERT INTO todo (title,status) VALUES (?,?)")
	_, err := stmt.Exec(data.Title, data.Status)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (r SqlRepo) Get(ctx context.Context, id string) (model.Todo, error) {
	var todo model.Todo
	row := r.db.QueryRow("SELECT id, title, status FROM todo WHERE id=?", id)

	err := row.Scan(&todo.Id, &todo.Title, &todo.Status)
	if err != nil {
		return todo, err
	}

	return todo, nil
}

func (r SqlRepo) Delete(ctx context.Context, id string) error {
	tx, _ := r.db.Begin()
	stmt, _ := tx.Prepare("DELETE FROM todo WHERE id=?")
	_, err := stmt.Exec(id)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (r SqlRepo) Update(ctx context.Context, id string, data model.Todo) error {
	tx, _ := r.db.Begin()
	stmt, _ := tx.Prepare("UPDATE todo SET title=?, status=? WHERE id=?")
	_, err := stmt.Exec(data.Title, data.Status, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}
