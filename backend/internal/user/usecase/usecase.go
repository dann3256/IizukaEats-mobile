// internal/usecase.go
package usecase

import (
	"context"
	"github.com/dann3256/IizukaEats-mobile/backend/internal/user/repository"
	"github.com/dann3256/IizukaEats-mobile/backend/internal/user/domain"
	"github.com/dann3256/IizukaEats-mobile/backend/internal/transport/http/ogen"
	"github.com/dann3256/IizukaEats-mobile/backend/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Usecase interface {
	RegisterUser(ctx context.Context, req *openapi.RegisterUserReq) (*domain.User, error)
	Login(ctx context.Context, req *openapi.LoginReq) (openapi.LoginResponse, error)
}

type usecaseImpl struct {
	repo repository.Repository // さっき作ったRepositoryインターフェース
	jwtManager *jwt.Manager // JWTマネージャー
}

func NewUsecase(repo repository.Repository, jwtManager *jwt.Manager) Usecase {
	return &usecaseImpl{
		repo:       repo,
		jwtManager: jwtManager,
	}
}



// ==================================================メソッド実装===============================================
func (u *usecaseImpl) RegisterUser(ctx context.Context, req *openapi.RegisterUserReq) (*domain.User, error) {
	
	// 1. パスワードをハッシュ化する
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.PasswordHash.Value), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// openapi.RegisterUserReq -> domain.User への変換
    params := &domain.User{
        Name:         string(req.Name.Value),
        Email:        string(req.Email.Value),
        PasswordHash: string(hashedPassword),
    }
	createdUser, err := u.repo.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}



func (u *usecaseImpl) Login(ctx context.Context, req *openapi.LoginReq) (openapi.LoginResponse, error) {
	// 1. メールアドレスでユーザーを取得
	user, err := u.repo.GetUserByEmail(ctx, string(req.Email.Value))
	if err != nil {
		return openapi.LoginResponse{}, err
	}	
	// 2. パスワードを検証
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.PasswordHash.Value))
	if err != nil {
		return openapi.LoginResponse{}, err
	}
	// u.jwtManagerのメソッドを呼び出す
	accessToken, refreshToken, err := u.jwtManager.GenerateTokensForUser(user.Name, user.Email, user.PasswordHash)
	if err != nil {
		return openapi.LoginResponse{}, err
	}

	// 4. レスポンスにトークンをセット
	return openapi.LoginResponse{
		AccessToken:  openapi.NewOptString(accessToken),
		RefreshToken: openapi.NewOptString(refreshToken),
	}, nil
}




// ==============================================================================================================