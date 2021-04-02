package factory

import (
	"database/sql"

	"github.com/kil-san/simple-todo-server/repo"
)

type repoFactory struct{}

type RepoFactory interface {
	CreateSqlRepo(db *sql.DB) repo.SqlRepo
}

func NewRepoFactory() RepoFactory {
	return &repoFactory{}
}

func (f *repoFactory) CreateSqlRepo(db *sql.DB) repo.SqlRepo {
	return repo.NewSqlRepo(db)
}
