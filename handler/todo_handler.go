package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/kil-san/simple-todo-server/connection"
	"github.com/kil-san/simple-todo-server/factory"
	"github.com/kil-san/simple-todo-server/model"
	"github.com/kil-san/simple-todo-server/service"
)

type TodoHandler struct {
	repoFactory factory.RepoFactory
}

func NewTodoHandler(repoFactory factory.RepoFactory) TodoHandler {
	return TodoHandler{
		repoFactory: repoFactory,
	}
}

func (h TodoHandler) Create(ctx context.Context, todo model.Todo) (error, int) {
	db, err := connection.NewSqliteConnection("sqlite.db")
	if err != nil {
		return err, http.StatusInternalServerError
	}
	defer db.Close()
	repo := h.repoFactory.CreateSqlRepo(db)
	svc := service.NewTodoService(repo)

	err = svc.CreateTodo(ctx, todo)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusCreated
}

func (h TodoHandler) Get(ctx context.Context, id string) ([]byte, error, int) {
	var response []byte
	db, err := connection.NewSqliteConnection("sqlite.db")
	if err != nil {
		return response, err, http.StatusInternalServerError
	}
	defer db.Close()
	repo := h.repoFactory.CreateSqlRepo(db)
	svc := service.NewTodoService(repo)

	todo, err := svc.GetTodo(ctx, id)
	if err != nil {
		return response, err, http.StatusInternalServerError
	}

	response, err = json.Marshal(todo)
	if err != nil {
		return response, err, http.StatusInternalServerError
	}

	return response, err, http.StatusOK
}
