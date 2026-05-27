#!/usr/bin/env bash
set -euo pipefail

attempts="${1:-5}"
delay_seconds="${2:-3}"

export GOPROXY="${GOPROXY:-https://proxy.golang.org,direct}"
export GOSUMDB="${GOSUMDB:-sum.golang.org}"

for attempt in $(seq 1 "${attempts}"); do
  echo "go mod download attempt ${attempt}/${attempts}"
  if go mod download; then
    exit 0
  fi

  if [ "${attempt}" -eq "${attempts}" ]; then
    echo "go mod download failed after ${attempts} attempts"
    exit 1
  fi

  sleep "${delay_seconds}"
done
