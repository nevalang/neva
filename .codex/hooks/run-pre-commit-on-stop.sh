#!/bin/sh
set -eu

repo_root=$(git rev-parse --show-toplevel)
cd "$repo_root"

changed_files=$( {
	git diff --name-only --diff-filter=ACMR
	git diff --cached --name-only --diff-filter=ACMR
	git ls-files --others --exclude-standard
} | sort -u )

if [ -z "$changed_files" ]; then
	exit 0
fi

printf '%s\n' "$changed_files" |
	tr '\n' '\0' |
	exec lefthook run pre-commit --files-from-stdin --no-stage-fixed
