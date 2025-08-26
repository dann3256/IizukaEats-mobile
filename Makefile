.PYHONY: install sqlc openapi

install:
	@bash scripts/install-tools.sh

sqlc:
	@bash scripts/generate-sqlc.sh

openapi:
	@bash scripts/generate-openapi.sh