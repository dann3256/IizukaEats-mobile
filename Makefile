.PHONY : install sqlc openapi tidy docker

install:
	@bash scripts/install-tools.sh

sqlc:
	@bash scripts/generate-sqlc.sh

openapi:
	@bash scripts/generate-openapi.sh

#必要なライブラリはこれでインストール
#不要なライブラリはgo mod tidyで削除
tidy:
	@go mod tidy

docker-up:
	@docker-compose up -d

docker-down:
	@docker-compose down