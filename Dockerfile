# ----------------- ステージ1: ビルド環境 -----------------
# Goの公式イメージをビルド用のベースイメージとして使用
FROM golang:1.23.0-alpine AS builder

# 作業ディレクトリを作成
WORKDIR /go/src/github.com/dann3256/IizukaEats-mobile/backend

# まず依存関係のファイルのみをコピーし、キャッシュを有効活用する
COPY backend/go.mod  backend/go.sum ./
RUN go mod download

# アプリケーションのソースコードを全てコピー
COPY backend/ ./
RUN go mod tidy

# アプリケーションをビルド
# CGO_ENABLED=0: C言語のライブラリに依存しない静的バイナリを生成
# GOOS=linux: Linux環境向けの実行ファイルを生成
# -o /server: ビルド成果物を /server という名前で出力
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server/main.go


# ----------------- ステージ2: 実行環境 -----------------
# 軽量なAlpine Linuxを最終的な実行環境のベースイメージとして使用
FROM alpine:latest

# 作業ディレクトリを作成
WORKDIR /app

# ビルドステージから生成された実行可能ファイルのみをコピー
COPY --from=builder /server .

# (オプション) データベースのマイグレーションファイルもコピーしておく
# コンテナに入ってマイグレーションを実行したい場合に便利
COPY backend/db/migrations ./db/migrations

# コンテナがリッスンするポートを指定（Goアプリの実装に合わせて変更してください）
EXPOSE 8080

# コンテナ起動時に実行するコマンド
CMD ["/app/server"]