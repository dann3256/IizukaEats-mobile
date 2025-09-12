// internal/repository.go
package repository

import (
	"context"
	"github.com/dann3256/IizukaEats-mobile/backend/internal/infrastructure/db/sqlc" // あなたのプロジェクトのパスに合わせてください
	"github.com/dann3256/IizukaEats-mobile/backend/internal/user/domain" // あなたのプロジェクトのパスに合わせてください
)

type Repository interface {
	CreateUser(ctx context.Context, user *domain.User) ( *domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type repositoryImpl struct {
	q *sqlc.Queries
}

func NewRepository(q *sqlc.Queries) Repository {
	return &repositoryImpl{q: q}
}

// ==================================================メソッド実装===============================================

func (r *repositoryImpl) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	// domain.User -> sqlc.CreateUserParams への変換
    params := sqlc.CreateUserParams{
        Name:         user.Name,
        Email:        user.Email,
        PasswordHash: user.PasswordHash,
    }

	createdSQLCUser, err := r.q.CreateUser(ctx, params)
    if err != nil {
        return nil, err
    }

	// sqlc.User -> domain.User への変換
    return &domain.User{
        ID:           createdSQLCUser.ID,
        Name:         createdSQLCUser.Name,
        Email:        createdSQLCUser.Email,
        PasswordHash: createdSQLCUser.PasswordHash,
    }, nil
	
}

func (r *repositoryImpl) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	getbyemail, err := r.q.GetUserByEmail(ctx,email)
    if err != nil {
        return nil, err
    }
    // sqlc.GetUserByEmailRow -> domain.User への変換   
    return &domain.User{
        ID:           getbyemail.ID,
        Name:         getbyemail.Name,
        Email:        getbyemail.Email,
        PasswordHash: getbyemail.PasswordHash,
    }, nil
}