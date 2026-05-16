const { existsSync } = require("node:fs");
const { join } = require("node:path");

const SUPPORTED_PLATFORMS = [
  "darwin-arm64",
  "darwin-x64",
  "linux-arm64",
  "linux-x64",
  "win32-x64",
];

function getPlatformKey(platform = process.platform, arch = process.arch) {
  return `${platform}-${arch === "x64" ? "x64" : arch}`;
}

function getBinaryName(platform = process.platform) {
  return platform === "win32" ? "comment-checker.exe" : "comment-checker";
}

function getBinaryPath(options = {}) {
  const platform = options.platform || process.platform;
  const arch = options.arch || process.arch;
  const platformKey = getPlatformKey(platform, arch);

  if (!SUPPORTED_PLATFORMS.includes(platformKey)) {
    throw new Error(
      `Unsupported platform: ${platform}-${arch}. ` +
        `Supported: ${SUPPORTED_PLATFORMS.join(", ")}`
    );
  }

  const binaryName = getBinaryName(platform);
  const baseDir = options.baseDir || __dirname;
  const bundledBinaryPath = join(baseDir, "vendor", platformKey, binaryName);

  if (existsSync(bundledBinaryPath)) {
    return bundledBinaryPath;
  }

  const localBinaryPath = join(baseDir, "bin", binaryName);
  if (existsSync(localBinaryPath)) {
    return localBinaryPath;
  }

  throw new Error(
    `comment-checker binary not found. ` +
      `Expected bundled binary at ${bundledBinaryPath} or postinstall binary at ${localBinaryPath}. ` +
      `Try reinstalling: npm install @code-yeongyu/comment-checker`
  );
}

module.exports = {
  getBinaryName,
  getBinaryPath,
  getPlatformKey,
  SUPPORTED_PLATFORMS,
};
