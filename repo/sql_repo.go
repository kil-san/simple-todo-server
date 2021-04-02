package repo

import (
	"context"
	"database/sql"
	"fmt"

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
	stmt, _ := tx.Prepare("insert into todo (title,status) values (?,?)")
	_, err := stmt.Exec(data.Title, data.Status)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func (r SqlRepo) Get(ctx context.Context, id string) (model.Todo, error) {
	var todo model.Todo
	row, err := r.db.Query("SELECT * FROM todo WHERE id=" + id)
	if err != nil {
		return todo, err
	}
	defer row.Close()

	for row.Next() {
		row.Scan(&todo.Id, &todo.Title, &todo.Status)
		return todo, nil
	}

	return todo, fmt.Errorf("Todo not Found\n")
}

func (r SqlRepo) Delete(ctx context.Context, id string) error {
	return nil
}

func (r SqlRepo) Update(ctx context.Context, id string, data model.Todo) error {
	return nil
}
