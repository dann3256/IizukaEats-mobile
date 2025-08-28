// internal/usecase.go
package usecase

import (
	"context"
	"github.com/dann3256/IizukaEats-mobile/backend/internal/user/repository"
	"github.com/dann3256/IizukaEats-mobile/backend/internal/infrastructure/db/sqlc" 
	"github.com/dann3256/IizukaEats-mobile/backend/internal/transport/http/ogen"

	"golang.org/x/crypto/bcrypt"
)

type Usecase interface {
	RegisterUser(ctx context.Context, req *openapi.RegisterUserReq) (sqlc.User, error)
}

type usecaseImpl struct {
	repo repository.Repository // さっき作ったRepositoryインターフェース
}

func NewUsecase(repo repository.Repository) Usecase {
	return &usecaseImpl{repo: repo}
}

func (u *usecaseImpl) RegisterUser(ctx context.Context, req *openapi.RegisterUserReq) (sqlc.User, error) {
	
	// 1. パスワードをハッシュ化する
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.PasswordHash.Value), bcrypt.DefaultCost)
	if err != nil {
		return sqlc.User{}, err
	}

	// 2. Repositoryを呼び出してDBに保存する
	params := sqlc.CreateUserParams{
		Name:         string(req.Name.Value),
		Email:        string(req.Email.Value),
		PasswordHash: string(hashedPassword),
	}
	createdUser, err := u.repo.CreateUser(ctx, params)
	if err != nil {
		return sqlc.User{}, err
	}

	return createdUser, nil
}