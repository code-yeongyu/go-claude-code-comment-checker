#!/usr/bin/env bash
# Usage: verify-npm-packages.sh <version>
set -euo pipefail

VERSION="$1"
ROOT=npm/comment-checker
KEYS=(darwin-arm64 darwin-x64 linux-arm64 linux-x64 win32-x64)

root_pack=$(cd "$ROOT" && npm pack --dry-run --json)
echo "$root_pack" | jq -e '
  .[0].files
  | map(.path)
  | sort == ["cli.js", "index.js", "package.json", "postinstall.js"]
'
expected_deps=$(printf '%s\n' "${KEYS[@]}" | jq -R --arg v "$VERSION" '{("@code-yeongyu/comment-checker-" + .): $v}' | jq -s 'add')
jq -e --argjson deps "$expected_deps" --arg v "$VERSION" \
  '.version == $v and .optionalDependencies == $deps' "$ROOT/package.json"

for key in "${KEYS[@]}"; do
  package="npm/comment-checker-$key"
  binary=comment-checker
  mode=493
  if [ "$key" = win32-x64 ]; then binary=comment-checker.exe; fi
  pack=$(cd "$package" && npm pack --dry-run --json)
  echo "$pack" | jq -e --arg bin "bin/$binary" --arg v "$VERSION" --argjson mode "$mode" '
    .[0].version == $v
    and (.[0].files | map(.path) | sort == [$bin, "package.json"])
    and (.[0].files | any(.path == $bin and .size > 1000000 and (.mode == $mode or ($bin | endswith(".exe")))))
  '
done
