package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kil-san/simple-todo-server/handler"
	"github.com/kil-san/simple-todo-server/middleware"
	"github.com/kil-san/simple-todo-server/model"
)

type TodoController interface {
	Routes() chi.Router
	Get(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
}

type todoController struct {
	handler handler.TodoHandler
}

func NewTodoController(todoHandler handler.TodoHandler) TodoController {
	return &todoController{
		handler: todoHandler,
	}
}

func (c *todoController) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.CORS())

	r.Get("/{todoID}", c.Get)
	r.Post("/", c.Post)
	r.Put("/{todoID}", c.Update)
	r.Delete("/{todoID}", c.Delete)

	return r
}

func (c *todoController) Get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	todoID := chi.URLParam(r, "todoID")

	_, err := strconv.ParseInt(todoID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err, status := c.handler.Get(ctx, todoID)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(status)
		return
	}

	w.Write(response)
	w.WriteHeader(status)
}

func (c *todoController) Post(w http.ResponseWriter, r *http.Request) {
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

	err, status := c.handler.Create(ctx, todo)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(status)
		return
	}

	w.WriteHeader(status)
}

func (c *todoController) Update(w http.ResponseWriter, r *http.Request) {
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

	err, status := c.handler.Update(ctx, todoID, todo)
	if err != nil {
		w.WriteHeader(status)
		return
	}

	w.WriteHeader(status)
}

func (c *todoController) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	todoID := chi.URLParam(r, "todoID")

	_, err := strconv.ParseInt(todoID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err, status := c.handler.Delete(ctx, todoID)
	if err != nil {
		w.WriteHeader(status)
		return
	}

	w.WriteHeader(status)
}
