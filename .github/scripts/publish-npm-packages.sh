#!/usr/bin/env bash
# Usage: publish-npm-packages.sh <version>
# Platform packages go first so the root package never points at a missing optional dependency.
set -euo pipefail

VERSION="$1"
PACKAGES=(
  npm/comment-checker-darwin-arm64
  npm/comment-checker-darwin-x64
  npm/comment-checker-linux-arm64
  npm/comment-checker-linux-x64
  npm/comment-checker-win32-x64
  npm/comment-checker
)

for package in "${PACKAGES[@]}"; do
  name=$(jq -r .name "$package/package.json")
  if [ "$(npm view "$name@$VERSION" version 2>/dev/null)" = "$VERSION" ]; then
    echo "$name@$VERSION is already published; skipping"
    continue
  fi
  (cd "$package" && npm publish --access public)
done
