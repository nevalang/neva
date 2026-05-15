#!/usr/bin/env sh
set -eu

: "${GOLANGCI_LINT_VERSION:=v2.5.0}"

files=$(git diff --cached --name-only -- '*.go')
if [ -z "$files" ]; then
  echo "no staged Go files"
  exit 0
fi

targets=$(printf '%s\n' "$files" | while IFS= read -r file; do
  dir=$(dirname "$file")
  if [ "$dir" = "." ]; then
    echo "./..."
  else
    echo "./$dir/..."
  fi
done | sort -u)

# shellcheck disable=SC2086
go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@"$GOLANGCI_LINT_VERSION" run --fix $targets
