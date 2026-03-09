#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="${1:-.}"
OUT_DIR="${2:-$REPO_ROOT/.agent/reports}"

mkdir -p "$OUT_DIR"

TS="$(date +"%Y%m%d-%H%M%S")"
LINT_LOG="$OUT_DIR/golangci-lint-$TS.log"
TEST_LOG="$OUT_DIR/go-test-$TS.log"
VULN_JSON="$OUT_DIR/govulncheck-$TS.json"
SUMMARY_JSON="$OUT_DIR/govuln-summary-$TS.json"
SUMMARY_TXT="$OUT_DIR/vuln-summary-$TS.txt"

ensure_tool() {
  local binary="$1"
  local install_path="$2"

  if command -v "$binary" >/dev/null 2>&1; then
    return
  fi

  echo "installing $binary via $install_path"
  go install "$install_path"
}

ensure_tool "govulncheck" "golang.org/x/vuln/cmd/govulncheck@latest"
ensure_tool "golangci-lint" "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest"

if [[ -x "$(go env GOPATH)/bin/golangci-lint" ]]; then
  LINT_BIN="$(go env GOPATH)/bin/golangci-lint"
else
  LINT_BIN="$(command -v golangci-lint)"
fi

if [[ -x "$(go env GOPATH)/bin/govulncheck" ]]; then
  VULN_BIN="$(go env GOPATH)/bin/govulncheck"
else
  VULN_BIN="$(command -v govulncheck)"
fi

if ! command -v jq >/dev/null 2>&1; then
  echo "jq is required to summarize govulncheck JSON output"
  exit 2
fi

pushd "$REPO_ROOT" >/dev/null

set +e
"$LINT_BIN" run --timeout=4m ./... >"$LINT_LOG" 2>&1
LINT_STATUS=$?

go test -count=1 -p 1 ./... >"$TEST_LOG" 2>&1
TEST_STATUS=$?

"$VULN_BIN" -json ./... >"$VULN_JSON"
VULN_STATUS=$?
set -e

jq -s '
  def finding_class($f):
    if ([$f.trace[]?.module] | index("github.com/nevalang/neva")) != null then "called"
    elif ([$f.trace[]? | select(has("package"))] | length) > 0 then "imported"
    else "required"
    end;
  (map(select(has("osv")) | .osv) | INDEX(.id)) as $osvmap |
  (map(select(has("finding")) | .finding) | group_by(.osv) |
    map({
      id: .[0].osv,
      class: (if (map(finding_class(.)) | index("called")) != null then "called" elif (map(finding_class(.)) | index("imported")) != null then "imported" else "required" end),
      fixed_version: .[0].fixed_version,
      found_modules: ([.[].trace[0].module] | unique),
      summary: ($osvmap[.[0].osv].summary // "")
    })
  )
' "$VULN_JSON" >"$SUMMARY_JSON"

CALLED_COUNT="$(jq '[.[] | select(.class == "called")] | length' "$SUMMARY_JSON")"
IMPORTED_COUNT="$(jq '[.[] | select(.class == "imported")] | length' "$SUMMARY_JSON")"
REQUIRED_COUNT="$(jq '[.[] | select(.class == "required")] | length' "$SUMMARY_JSON")"
TOTAL_COUNT="$(jq 'length' "$SUMMARY_JSON")"

{
  echo "vuln-scan timestamp: $TS"
  echo "lint status: $LINT_STATUS"
  echo "go test status: $TEST_STATUS"
  echo "govulncheck status: $VULN_STATUS"
  echo "vulnerabilities: total=$TOTAL_COUNT called=$CALLED_COUNT imported=$IMPORTED_COUNT required=$REQUIRED_COUNT"
  echo
  echo "reachable vulnerabilities (called):"
  jq -r '.[] | select(.class == "called") | "- \(.id) | fixed: \(.fixed_version) | module(s): \(.found_modules | join(",")) | \(.summary)"' "$SUMMARY_JSON"
} >"$SUMMARY_TXT"

popd >/dev/null

echo "logs:"
echo "  $LINT_LOG"
echo "  $TEST_LOG"
echo "  $VULN_JSON"
echo "  $SUMMARY_JSON"
echo "  $SUMMARY_TXT"

if [[ "$LINT_STATUS" -ne 0 || "$TEST_STATUS" -ne 0 || "$VULN_STATUS" -ne 0 || "$CALLED_COUNT" -gt 0 ]]; then
  exit 1
fi
