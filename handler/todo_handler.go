package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kil-san/simple-todo-server/factory"
	"github.com/kil-san/simple-todo-server/middleware"
	"github.com/kil-san/simple-todo-server/model"
	"github.com/kil-san/simple-todo-server/service"
)

type TodoHandler interface {
	Routes() chi.Router
	Get(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type todoHandler struct {
	db          *sql.DB
	repoFactory factory.RepoFactory
}

func NewTodoHandler(db *sql.DB, repoFactory factory.RepoFactory) TodoHandler {
	return &todoHandler{
		db:          db,
		repoFactory: repoFactory,
	}
}

func (c *todoHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.CORS())

	r.Get("/{todoID}", c.Get)
	r.Post("/", c.Post)
	r.Put("/{todoID}", c.Update)
	r.Delete("/{todoID}", c.Delete)

	return r
}

func (c *todoHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	todoID := chi.URLParam(r, "todoID")

	_, err := strconv.ParseInt(todoID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var response []byte
	repo := c.repoFactory.CreateSqlRepo(c.db)
	svc := service.NewTodoService(repo)
	todo, err := svc.GetTodo(ctx, todoID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	response, err = json.Marshal(todo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(response)
}

func (c *todoHandler) Post(w http.ResponseWriter, r *http.Request) {
	var todo model.Todo
	rawBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(rawBody))
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	err = json.Unmarshal(rawBody, &todo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	repo := c.repoFactory.CreateSqlRepo(c.db)
	svc := service.NewTodoService(repo)
	err = svc.CreateTodo(ctx, todo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *todoHandler) Update(w http.ResponseWriter, r *http.Request) {
	var todo model.Todo
	rawBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(rawBody))
	todoID := chi.URLParam(r, "todoID")
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	err = json.Unmarshal(rawBody, &todo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	repo := c.repoFactory.CreateSqlRepo(c.db)
	svc := service.NewTodoService(repo)
	err = svc.UpdateTodo(ctx, todoID, todo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *todoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	todoID := chi.URLParam(r, "todoID")

	_, err := strconv.ParseInt(todoID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	repo := c.repoFactory.CreateSqlRepo(c.db)
	svc := service.NewTodoService(repo)
	err = svc.DeleteTodo(ctx, todoID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}
