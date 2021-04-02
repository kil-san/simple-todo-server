package repo

import (
	"context"

	"github.com/kil-san/simple-todo-server/model"
)

type Repo interface {
	Create(ctx context.Context, data model.Todo) error
	Get(ctx context.Context, id string) (model.Todo, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, data model.Todo) error
}
