package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/kil-san/simple-todo-server/factory"
	"github.com/kil-san/simple-todo-server/model"
	"github.com/kil-san/simple-todo-server/service"
)

type TodoHandler struct {
	db          *sql.DB
	repoFactory factory.RepoFactory
}

func NewTodoHandler(db *sql.DB, repoFactory factory.RepoFactory) TodoHandler {
	return TodoHandler{
		db:          db,
		repoFactory: repoFactory,
	}
}

func (h TodoHandler) Create(ctx context.Context, todo model.Todo) (error, int) {
	repo := h.repoFactory.CreateSqlRepo(h.db)
	svc := service.NewTodoService(repo)

	err := svc.CreateTodo(ctx, todo)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusCreated
}

func (h TodoHandler) Get(ctx context.Context, id string) ([]byte, error, int) {
	var response []byte
	repo := h.repoFactory.CreateSqlRepo(h.db)
	svc := service.NewTodoService(repo)

	todo, err := svc.GetTodo(ctx, id)
	if err != nil {
		return response, err, http.StatusNotFound
	}

	response, err = json.Marshal(todo)
	if err != nil {
		return response, err, http.StatusInternalServerError
	}

	return response, err, http.StatusOK
}

func (h TodoHandler) Update(ctx context.Context, id string, todo model.Todo) (error, int) {
	repo := h.repoFactory.CreateSqlRepo(h.db)
	svc := service.NewTodoService(repo)

	err := svc.UpdateTodo(ctx, id, todo)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func (h TodoHandler) Delete(ctx context.Context, id string) (error, int) {
	repo := h.repoFactory.CreateSqlRepo(h.db)
	svc := service.NewTodoService(repo)

	err := svc.DeleteTodo(ctx, id)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	return err, http.StatusOK
}
