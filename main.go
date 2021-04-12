package main

import (
	"fmt"
	"net/http"

	"github.com/kil-san/simple-todo-server/connection"
	"github.com/kil-san/simple-todo-server/factory"
	"github.com/kil-san/simple-todo-server/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "unknwon.dev/clog/v2"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	db, err := connection.NewSqliteConnection("sqlite.db")
	if err != nil {
		panic("could not open connection to db")
	}
	defer db.Close()
	repoFactory := factory.NewRepoFactory()
	todoHandler := handler.NewTodoHandler(db, repoFactory)

	err = log.NewConsole()
	if err != nil {
		panic("unable to create new logger: " + err.Error())
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Mount("/todos", todoHandler.Routes())

	fmt.Printf("Server running at localhost:8000\n")
	http.ListenAndServe(":8000", r)
}
