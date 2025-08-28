package handler

import (
    "context"
    "errors"

    "github.com/dann3256/IizukaEats-mobile/backend/internal/transport/http/ogen"
    "github.com/dann3256/IizukaEats-mobile/backend/internal/user/usecase"
)

type APIHandler struct {
    uc usecase.Usecase
}

func NewAPIHandler(uc usecase.Usecase) openapi.Handler {
    return &APIHandler{uc: uc}
}

// RegisterUser は /register のリクエストを処理する
func (h *APIHandler) RegisterUser(ctx context.Context, req openapi.OptRegisterUserReq) (openapi.RegisterUserRes, error) {
    // リクエストが設定されているか確認
    if !req.Set {
        return nil, errors.New("request is not set")
    }

    // Usecaseを呼び出す
    createdUser, err := h.uc.RegisterUser(ctx, &req.Value)
    if err != nil {
        return nil, err
    }

    // 成功レスポンスを組み立てる
    return &openapi.UserResponse{
        Name:         openapi.NewOptName(openapi.Name(createdUser.Name)),
        Email:        openapi.NewOptEmail(openapi.Email(createdUser.Email)),
        PasswordHash: openapi.NewOptPasswordHash(openapi.PasswordHash(createdUser.PasswordHash)),
    }, nil
}

// GetUser は /me のリクエストを処理する（未実装）
func (h *APIHandler) GetUser(ctx context.Context) (openapi.GetUserRes, error) {
    // 必要に応じてデータを取得してレスポンスを構築
    return &openapi.UserResponse{
        Name:         openapi.NewOptName("John Doe"),
        Email:        openapi.NewOptEmail("john.doe@example.com"),
        PasswordHash: openapi.NewOptPasswordHash("hashed_password"),
    }, nil
}

// Login は /login のリクエストを処理する（未実装）
func (h *APIHandler) Login(ctx context.Context, req openapi.OptLoginReq) (openapi.LoginRes, error) {
    // 必要に応じてデータを取得してレスポンスを構築
    return &openapi.LoginResponse{
        AccessToken:  openapi.NewOptString("dummy_access_token"),
		RefreshToken: openapi.NewOptString("dummy_refresh_token"),
    }, nil
}