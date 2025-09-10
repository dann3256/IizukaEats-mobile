.PHONY : install sqlc openapi tidy docker

install:
	@bash scripts/install-tools.sh

sqlc:
	@bash scripts/generate-sqlc.sh

openapi:
	@bash scripts/generate-openapi.sh

#必要なライブラリはこれでインストール,不要なライブラリは削除
tidy:
	@cd backend && go mod tidy

# イメージ作成
docker-build:
	@docker-compose build
# イメージを使ってコンテナを起動
docker-up:
	@docker-compose up -d
# コンテナ停止
docker-down:
	@docker-compose down

# 既存のコンテナを再利用せず、強制的に作り直して起動
docker-rebuild:
	@docker-compose up -d --force-recreate