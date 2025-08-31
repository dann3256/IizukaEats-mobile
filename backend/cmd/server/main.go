package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"github.com/dann3256/IizukaEats-mobile/backend/internal/infrastructure/db/sqlc"
	"github.com/dann3256/IizukaEats-mobile/backend/internal/transport/http/ogen"
	"github.com/dann3256/IizukaEats-mobile/backend/internal/user/handler"
	"github.com/dann3256/IizukaEats-mobile/backend/internal/user/repository"
	"github.com/dann3256/IizukaEats-mobile/backend/internal/user/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq" // PostgreSQLドライバ

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // PostgreSQL用のドライバ
	_ "github.com/golang-migrate/migrate/v4/source/file"       // ファイルからマイグレーションを読み込むためのドライバ
)

type SecurityHandler struct{}

// HandleBearerAuth は Bearer トークン認証を処理する
func (s *SecurityHandler) HandleBearerAuth(ctx context.Context, operationName openapi.OperationName, t openapi.BearerAuth) (context.Context, error) {
	// Bearer トークンの認証ロジックを実装
	if t.Token == "" {
		return ctx, errors.New("missing token")
	}

	// 仮の認証成功ロジック
	// 必要に応じてトークンを検証し、ユーザー情報をコンテキストに追加
	return ctx, nil
}

func main() {
	// データベース接続情報
	dsn := "postgres://user:password@db:5432/iizukaeats_db?sslmode=disable"

	// --- ▼▼▼ ここからマイグレーション処理を追加 ▼▼▼ ---
	log.Println("データベースマイグレーションを開始します...")
	m, err := migrate.New(
		"file://./db/migrations", // ← typoを修正し、Docker内のパスに合わせる
		dsn,
	)
	if err != nil {
		log.Fatalf("マイグレーションの初期化に失敗しました: %v", err)
	}

	// マイグレーションを適用
	if err := m.Up(); err != nil {
		// ErrNoChangeは「既にデータベースが最新で変更がない」という正常なエラーなので、それ以外の場合のみエラーとする
		if err != migrate.ErrNoChange {
			log.Fatalf("マイグレーションの適用に失敗しました: %v", err)
		}
		log.Println("データベースは既に最新です。")
	} else {
		log.Println("データベースマイグレーションが正常に完了しました。")
	}

	// データベース接続 (マイグレーションが完了した後に行う)
	dbpool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("DB接続失敗: %v", err)
	}
	defer dbpool.Close()

	// 各レイヤーを作成
	queries := sqlc.New(dbpool)
	repo := repository.NewRepository(queries)
	uc := usecase.NewUsecase(repo)
	h := handler.NewAPIHandler(uc)

	// SecurityHandler を作成
	secHandler := &SecurityHandler{}

	// ogenサーバーを作成
	srv, err := openapi.NewServer(h, secHandler)
	if err != nil {
		log.Fatalf("サーバー作成失敗: %v", err)
	}

	// HTTPサーバー起動
	port := 8080
	log.Printf("サーバー起動 http://localhost:%d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), srv); err != nil {
		log.Fatalf("サーバー起動失敗: %v", err)
	}
}