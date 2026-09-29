#!/usr/bin/env bash
# Usage: prepare-npm-packages.sh <version> <tag> <artifacts-dir>
set -euo pipefail

VERSION="$1"
TAG="$2"
ARTIFACTS="$3"
ROOT=npm/comment-checker
STAGE=$(mktemp -d)
trap 'rm -rf "$STAGE"' EXIT

# npm platform key : release archive suffix : binary name
PLATFORMS=(
  "darwin-arm64:darwin_arm64.tar.gz:comment-checker"
  "darwin-x64:darwin_amd64.tar.gz:comment-checker"
  "linux-arm64:linux_arm64.tar.gz:comment-checker"
  "linux-x64:linux_amd64.tar.gz:comment-checker"
  "win32-x64:windows_amd64.zip:comment-checker.exe"
)

for entry in "${PLATFORMS[@]}"; do
  IFS=: read -r key suffix binary <<<"$entry"
  package="npm/comment-checker-$key"
  archive="$ARTIFACTS/comment-checker_${TAG}_${suffix}"
  mkdir -p "$STAGE/$key" "$package/bin"
  case "$archive" in
    *.zip) unzip -o -q "$archive" -d "$STAGE/$key" ;;
    *) tar -xzf "$archive" -C "$STAGE/$key" ;;
  esac
  cp "$STAGE/$key/$binary" "$package/bin/$binary"
  chmod 755 "$package/bin/$binary"
  jq --arg v "$VERSION" '.version = $v' "$package/package.json" >"$package/package.json.tmp"
  mv "$package/package.json.tmp" "$package/package.json"
done

jq --arg v "$VERSION" '.version = $v | .optionalDependencies |= with_entries(.value = $v)' \
  "$ROOT/package.json" >"$ROOT/package.json.tmp"
mv "$ROOT/package.json.tmp" "$ROOT/package.json"
