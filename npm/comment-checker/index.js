const { existsSync } = require("fs");
const { join } = require("path");

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

function getBinaryPath() {
  const platformKey = getPlatformKey();

  if (!platformKey) {
    throw new Error(
      `Unsupported platform: ${process.platform}-${process.arch}. ` +
        `Supported: darwin-arm64, darwin-x64, linux-arm64, linux-x64`
    );
  }

  const packageName = PLATFORM_PACKAGES[platformKey];

  // Try to find the platform-specific package
  try {
    const packagePath = require.resolve(`${packageName}/package.json`);
    const binDir = join(packagePath, "..", "bin");
    const binaryPath = join(binDir, "comment-checker");

    if (existsSync(binaryPath)) {
      return binaryPath;
    }
  } catch {
    // Package not found
  }

  // Fallback: check if binary exists in this package's bin directory
  const localBinaryPath = join(__dirname, "bin", "comment-checker");
  if (existsSync(localBinaryPath)) {
    return localBinaryPath;
  }

  throw new Error(
    `comment-checker binary not found. ` +
      `Platform package ${packageName} may not be installed. ` +
      `Try reinstalling: npm install @anthropic-claude/comment-checker`
  );
}

module.exports = {
  getBinaryPath,
  getPlatformKey,
  PLATFORM_PACKAGES,
};
