// internal/repository.go
package repository

import (
	"context"
	"github.com/dann3256/IizukaEats-mobile/backend/internal/infrastructure/db/sqlc" // あなたのプロジェクトのパスに合わせてください
)

type Repository interface {
	CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error)
}

type repositoryImpl struct {
	q *sqlc.Queries
}

func NewRepository(q *sqlc.Queries) Repository {
	return &repositoryImpl{q: q}
}

// CreateUser はsqlcのCreateUserを呼び出すだけ
func (r *repositoryImpl) CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error) {
	return r.q.CreateUser(ctx, arg)
}