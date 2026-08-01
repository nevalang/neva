#!/usr/bin/env bash

set -euo pipefail

repo_root=$(cd "$(dirname "$0")/.." && pwd)
temp_dir=$(mktemp -d)
trap 'rm -rf "$temp_dir"' EXIT

mock_bin="$temp_dir/mock-bin"
mkdir -p "$mock_bin"

printf '%s\n' '#!/usr/bin/env bash' 'set -euo pipefail' 'if [[ "${MOCK_CURL_FAIL:-}" == 1 ]]; then exit 23; fi' 'if [[ "$*" == *"api.github.com/repos/nevalang/neva/releases/latest"* ]]; then echo '\''{"tag_name":"v9.9.9"}'\''; exit 0; fi' 'while [[ $# -gt 0 ]]; do if [[ "$1" == "-o" ]]; then printf "%s\\n" "#!/usr/bin/env sh" "exit 0" > "$2"; exit 0; fi; shift; done' 'exit 1' > "$mock_bin/curl"
chmod +x "$mock_bin/curl"

install_dir="$temp_dir/install"
env HOME="$temp_dir/home" PATH="$mock_bin:$PATH" NEVA_INSTALL_DIR="$install_dir" bash "$repo_root/scripts/install.sh"
test -x "$install_dir/neva"

if env HOME="$temp_dir/fail-home" PATH="$mock_bin:$PATH" MOCK_CURL_FAIL=1 NEVA_INSTALL_DIR="$temp_dir/fail-install" bash "$repo_root/scripts/install.sh"; then
	echo 'installer unexpectedly succeeded after a download failure' >&2
	exit 1
fi
test ! -e "$temp_dir/fail-install/neva"
