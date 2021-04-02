package service

import (
	"context"

	"github.com/kil-san/simple-todo-server/model"
	"github.com/kil-san/simple-todo-server/repo"
)

type TodoService struct {
	repo repo.Repo
}

func NewTodoService(repo repo.Repo) *TodoService {
	return &TodoService{
		repo: repo,
	}
}

func (s *TodoService) CreateTodo(ctx context.Context, todo model.Todo) error {
	s.repo.Create(ctx, todo)
	return nil
}

func (s *TodoService) GetTodo(ctx context.Context, id string) (model.Todo, error) {
	var todo model.Todo
	todo, err := s.repo.Get(ctx, id)
	if err != nil {
		return todo, err
	}
	return todo, nil
}
