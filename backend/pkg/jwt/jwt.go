package jwt

import (
	"fmt"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

// Manager はJWTの生成と検証を管理します。
type Manager struct {
	secretKey string
}

// NewManager は新しいJWTマネージャーを作成します。
// アプリケーション起動時に一度だけ呼び出し、DI（依存性注入）で使い回すことを想定しています。
func NewManager(secretKey string) (*Manager, error) {
	if secretKey == "" {
		return nil, fmt.Errorf("jwt secret key cannot be empty")
	}
	return &Manager{secretKey: secretKey}, nil
}

// Generate は与えられたClaimsからJWT文字列を生成する汎用的なメソッドです。
func (m *Manager) Generate(claims jwt.Claims) (string, error) {
	if claims == nil {
		return "", fmt.Errorf("claims cannot be empty")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

// Validate は与えられたJWT文字列を検証し、パースされたトークンを返します。
func (m *Manager) Validate(tokenStr string) (*jwt.Token, error) {
	if tokenStr == "" {
		return nil, fmt.Errorf("token string cannot be empty")
	}

	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// 署名方式がHMACであることを検証 (アルゴリズム偽装攻撃対策)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.secretKey), nil
	})
}


// --- ここからアプリケーション固有のロジック ---

const (
	accessTokenDuration  = 15 * time.Minute
	refreshTokenDuration = 7 * 24 * time.Hour
)

// AccessTokenClaims はアクセストークンに含まれる情報です。
type AccessTokenClaims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateTokensForUser はユーザー情報からアクセストークンとリフレッシュトークンを生成
func (m *Manager) GenerateTokensForUser(userName string, email string, PasswordHash string) (accessToken string, refreshToken string, err error) {
	

	// アクセストークン
	accessClaims := AccessTokenClaims{
		Name:  userName,
		Email: email, 
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   email,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken, err = m.Generate(&accessClaims)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// リフレッシュトークン
	refreshClaims := jwt.RegisteredClaims{
		
		Subject:   email,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenDuration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	refreshToken, err = m.Generate(&refreshClaims)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}
