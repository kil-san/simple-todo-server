package main

import (
	"bytes"
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kil-san/simple-todo-server/connection"
	"github.com/kil-san/simple-todo-server/controller"
	"github.com/kil-san/simple-todo-server/factory"
	"github.com/kil-san/simple-todo-server/handler"
)

var testDB = "test.db"

func init() {
	if _, err := os.Stat("db/" + testDB); os.IsNotExist(err) {
		return
	}
	os.Remove("db/" + testDB)
}

func setup() (*chi.Mux, *sql.DB, error) {
	var r *chi.Mux
	var db *sql.DB
	db, err := connection.NewSqliteConnection(testDB)
	if err != nil {
		return r, db, err
	}

	repoFactory := factory.NewRepoFactory()
	todoHandler := handler.NewTodoHandler(db, repoFactory)
	todoController := controller.NewTodoController(todoHandler)
	r = chi.NewRouter()
	r.Mount("/todos", todoController.Routes())

	return r, db, nil
}

func TestInvalidTodo(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/todos/10", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	rec := httptest.NewRecorder()

	r, db, err := setup()
	if err != nil {
		t.Fatalf("could not create database connection: %v", err)
	}
	defer db.Close()

	r.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status NOTFOUND; got %v", res.Status)
	}
}

func TestCreateTodo(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var data = []byte(`{"title":"Read a book","status":"done"}`)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/todos", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	rec := httptest.NewRecorder()

	r, db, err := setup()
	if err != nil {
		t.Fatalf("could not create database connection: %v", err)
	}
	defer db.Close()

	r.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected status CREATED; got %v", res.Status)
	}
}

func TestGetTodo(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/todos/1", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	rec := httptest.NewRecorder()

	r, db, err := setup()
	if err != nil {
		t.Fatalf("could not create database connection: %v", err)
	}
	defer db.Close()

	r.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status OK; got %v", res.Status)
	}
}
