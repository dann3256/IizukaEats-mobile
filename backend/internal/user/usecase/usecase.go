// internal/usecase.go
package usecase

import (
	"context"
	"github.com/dann3256/IizukaEats-mobile/backend/internal/user/repository"
	"github.com/dann3256/IizukaEats-mobile/backend/internal/user/domain"
	"github.com/dann3256/IizukaEats-mobile/backend/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Usecase interface {
	RegisterUser(ctx context.Context, req *CreateUserInputDTO) (*CreateUserOutputDTO, error)
	Login(ctx context.Context, req *LoginInputDTO) (*LoginOutputDTO, error)
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
type CreateUserInputDTO struct {
    Name     string
    Email    string
    Password string // 生パスワードを受け取る
}
type CreateUserOutputDTO struct {
	ID    int32
	Name  string
	Email string
}
func (u *usecaseImpl) RegisterUser(ctx context.Context, req *CreateUserInputDTO) (*CreateUserOutputDTO, error) {
	
	// 1. パスワードをハッシュ化する
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	// CreateUserInputDTO -> domain.User への変換
    params := &domain.User{
        Name:         string(req.Name),
        Email:        string(req.Email),
        PasswordHash: string(hashedPassword),
    }

	domainUser, err := u.repo.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}
	// domain.User -> CreateUserOutputDTO への変換		
	createdUser := &CreateUserOutputDTO{
		ID:    domainUser.ID,
		Name:  domainUser.Name,
		Email: domainUser.Email,
	}
	return createdUser, nil
}



type LoginInputDTO struct {
	Email    string
	Password string
}
type LoginOutputDTO struct {
	AccessToken  string
	RefreshToken string
}
func (u *usecaseImpl) Login(ctx context.Context, req * LoginInputDTO) (*LoginOutputDTO, error) {
	// 1. メールアドレスでユーザーを取得
	user, err := u.repo.GetUserByEmail(ctx, string(req.Email))
	if err != nil {
		return nil, err
	}	
	// 2. パスワードを検証
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, err
	}
	// u.jwtManagerのメソッドを呼び出す
	accessToken, refreshToken, err := u.jwtManager.GenerateTokensForUser(user.Name, user.Email, user.PasswordHash)
	if err != nil {
		return nil, err
	}

	// 4. レスポンスにトークンをセット
	return &LoginOutputDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ==============================================================================================================