set -euo pipefail

ROOT_DIR=$(git rev-parse --show-toplevel)
BIN_DIR="${ROOT_DIR}/backend/.bin"
OPENAPI_OUT_DIR="${ROOT_DIR}/backend/internal/transport/http/ogen"
BUNDLE_FILE="${ROOT_DIR}/backend/docs/api/bundle.yaml"

rm -f "$BUNDLE_FILE"
npx --prefix "$BIN_DIR" redocly bundle  backend/docs/api/openapi.yaml -o backend/docs/api/bundle.yaml 
"${BIN_DIR}/ogen" --target "${OPENAPI_OUT_DIR}" --clean --package openapi backend/docs/api/bundle.yaml