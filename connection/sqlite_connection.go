package connection

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func NewSqliteConnection(filename string) (*sql.DB, error) {
	path, err := filepath.Abs("db/" + filename)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		file.Close()
	}

	sqliteDatabase, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	createTodoTableSQL := `CREATE TABLE IF NOT EXISTS todo (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"title" TEXT,
		"status" TEXT
	  );`

	stmt, err := sqliteDatabase.Prepare(createTodoTableSQL)
	if err != nil {
		return nil, err
	}
	stmt.Exec()

	return sqliteDatabase, nil
}
