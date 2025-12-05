#!/usr/bin/env node

const { existsSync, mkdirSync, symlinkSync, unlinkSync, chmodSync } = require("fs");
const { join, dirname } = require("path");

const PLATFORM_PACKAGES = {
  "darwin-arm64": "@anthropic-claude/comment-checker-darwin-arm64",
  "darwin-x64": "@anthropic-claude/comment-checker-darwin-x64",
  "linux-arm64": "@anthropic-claude/comment-checker-linux-arm64",
  "linux-x64": "@anthropic-claude/comment-checker-linux-x64",
};

function getPlatformKey() {
  const platform = process.platform;
  const arch = process.arch;

  if (platform === "darwin") {
    return arch === "arm64" ? "darwin-arm64" : "darwin-x64";
  }
  if (platform === "linux") {
    return arch === "arm64" ? "linux-arm64" : "linux-x64";
  }

  return null;
}

function main() {
  const platformKey = getPlatformKey();

  if (!platformKey) {
    console.warn(
      `[comment-checker] Warning: Unsupported platform ${process.platform}-${process.arch}`
    );
    process.exit(0);
  }

  const packageName = PLATFORM_PACKAGES[platformKey];
  let sourceBinaryPath;

  try {
    const packagePath = require.resolve(`${packageName}/package.json`);
    sourceBinaryPath = join(dirname(packagePath), "bin", "comment-checker");
  } catch {
    console.warn(
      `[comment-checker] Warning: Platform package ${packageName} not found. ` +
        `Binary may need to be downloaded manually.`
    );
    process.exit(0);
  }

  if (!existsSync(sourceBinaryPath)) {
    console.warn(
      `[comment-checker] Warning: Binary not found at ${sourceBinaryPath}`
    );
    process.exit(0);
  }

  // Create symlink in this package's bin directory
  const binDir = join(__dirname, "bin");
  const targetBinaryPath = join(binDir, "comment-checker");

  if (!existsSync(binDir)) {
    mkdirSync(binDir, { recursive: true });
  }

  // Remove existing symlink if present
  if (existsSync(targetBinaryPath)) {
    unlinkSync(targetBinaryPath);
  }

  // Create symlink to platform-specific binary
  symlinkSync(sourceBinaryPath, targetBinaryPath);

  // Ensure executable
  try {
    chmodSync(targetBinaryPath, 0o755);
  } catch {
    // Ignore chmod errors on symlinks
  }

  console.log(`[comment-checker] Linked ${platformKey} binary`);
}

main();
