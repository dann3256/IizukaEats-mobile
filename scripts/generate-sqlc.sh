set -euo pipefail

ROOT_DIR=$(git rev-parse --show-toplevel)
BIN_DIR="${ROOT_DIR}/backend/.bin"
SQLC_OUT_DIR="${ROOT_DIR}/backend/internal/db/sqlc"

rm -rf "${SQLC_OUT_DIR}"

(
  cd backend && "${BIN_DIR}/sqlc" generate
)