package handler

import (
    "context"

    "github.com/dann3256/IizukaEats-mobile/backend/internal/transport/http/ogen"
    "github.com/dann3256/IizukaEats-mobile/backend/internal/user/usecase"
    
)

type APIHandler struct {
    uc usecase.Usecase
}

func NewAPIHandler(uc usecase.Usecase) openapi.Handler {
    return &APIHandler{uc: uc}
}

//  /register のリクエストを処理する
func (h *APIHandler) RegisterUser(ctx context.Context, req *openapi.RegisterUserReq) (openapi.RegisterUserRes, error) {
    // Usecaseを呼び出す
    domainUser, err := h.uc.RegisterUser(ctx, req)
    if err != nil {
        return nil, err
    }

     // domain.User -> openapi.User (レスポンス用) への変換
   response := &openapi.UserResponse{
    ID:    openapi.ID(domainUser.ID),
    Name:  openapi.Name(domainUser.Name),
    Email: openapi.Email(domainUser.Email),
    // PasswordHash はレスポンスに含めるべきではない
}

    return response, nil
}

//  /login のリクエストを処理する
func (h *APIHandler) Login(ctx context.Context, req openapi.OptLoginReq) (openapi.LoginRes, error) {
    // リクエストが設定されているか確認
    
    // 必要に応じてデータを取得してレスポンスを構築
    loginResponse, err := h.uc.Login(ctx, &req.Value)
	if err != nil {
		// Usecaseから返されたエラーを返す (例: 401 Unauthorized)
		return nil,err
	}

	// 3. Usecaseからのレスポンスを返す
	return &loginResponse, nil
}



// GetUser は /me のリクエストを処理する（未実装）
func (h *APIHandler) GetUser(ctx context.Context) (openapi.GetUserRes, error) {
    // 必要に応じてデータを取得してレスポンスを構築
    return &openapi.UserResponse{
        Name:         openapi.Name("John Doe"),
        Email:        openapi.Email("john.doe@example.com"),
        
    }, nil
}

