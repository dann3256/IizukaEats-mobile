# ----------------- ステージ1: ビルド環境 (builder) -----------------
# Goの環境とalpneのユーザーランドを組み合わせたイメージを使用
# ビルド専用のステージ（一時的な作業場）として定義
FROM golang:1.25-alpine AS builder

# 作業ディレクトリをシンプルに設定
WORKDIR /app

# 依存関係のファイルのみをコピー
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# アプリケーションのソースコードを全てコピー
COPY backend/ ./

# アプリケーションをビルド
# GoプログラムからC言語のコードを呼び出すCGoを無効化し、Cライブラリに依存しないバイナリファイルを生成
# Go言語のクロスコンパイル(今いるOS以外のOC向けのバイナリファイルを作れる)を有効化し、Linux OS向けのバイナリを生成
# 成果物を現在の作業ディレクトリ(/app)に 'server' という名前で出力
RUN CGO_ENABLED=0 GOOS=linux go build -o ./server ./backend/cmd/server/main.go


# ----------------- ステージ2: 実行環境 (final) -----------------
FROM alpine:latest

# 作業ディレクトリを作成
WORKDIR /app

# ビルドステージの作業ディレクトリから、生成された実行可能ファイルのみをコピー
# コピー元とコピー先のパスが明確になる
COPY --from=builder /app/server .

# マイグレーションファイルをコピー
COPY backend/db/migrations ./db/migrations

# コンテナ起動時に実行するコマンド
# WORKDIRが/appだから./server で実行できる
CMD ["./server"]


