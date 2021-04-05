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
	err := s.repo.Create(ctx, todo)
	if err != nil {
		return err
	}
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

func (s *TodoService) UpdateTodo(ctx context.Context, id string, todo model.Todo) error {
	err := s.repo.Update(ctx, id, todo)
	if err != nil {
		return err
	}
	return nil
}

func (s *TodoService) DeleteTodo(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
