const { existsSync } = require("node:fs");
const { createRequire } = require("node:module");
const { join } = require("node:path");

const PLATFORM_PACKAGES = {
  "darwin-arm64": "@code-yeongyu/comment-checker-darwin-arm64",
  "darwin-x64": "@code-yeongyu/comment-checker-darwin-x64",
  "linux-arm64": "@code-yeongyu/comment-checker-linux-arm64",
  "linux-x64": "@code-yeongyu/comment-checker-linux-x64",
  "win32-x64": "@code-yeongyu/comment-checker-win32-x64",
};

function getPlatformKey(platform = process.platform, arch = process.arch) {
  return `${platform}-${arch === "x64" ? "x64" : arch}`;
}

function getBinaryName(platform = process.platform) {
  return platform === "win32" ? "comment-checker.exe" : "comment-checker";
}

function getResolver(resolveFrom = __dirname) {
  return createRequire(join(resolveFrom, "package.json"));
}

function getBinaryPath(options = {}) {
  const platform = options.platform || process.platform;
  const arch = options.arch || process.arch;
  const platformKey = getPlatformKey(platform, arch);
  const packageName = PLATFORM_PACKAGES[platformKey];

  if (!packageName) {
    throw new Error(
      `Unsupported platform: ${platform}-${arch}. ` +
        `Supported: ${Object.keys(PLATFORM_PACKAGES).join(", ")}`
    );
  }

  const binaryName = getBinaryName(platform);
  const resolver = options.resolver || getResolver(options.resolveFrom);
  let packageResolveError = null;

  try {
    const packagePath = resolver.resolve(`${packageName}/package.json`);
    const binaryPath = join(packagePath, "..", "bin", binaryName);

    if (existsSync(binaryPath)) {
      return binaryPath;
    }
  } catch (error) {
    packageResolveError = error;
  }

  const localBinaryPath = join(__dirname, "bin", binaryName);
  if (existsSync(localBinaryPath)) {
    return localBinaryPath;
  }

  const resolveDetails = packageResolveError
    ? ` ${packageName} resolution failed: ${packageResolveError.message}.`
    : "";

  throw new Error(
    `comment-checker binary not found. ` +
      `Platform package ${packageName} may not be installed. ` +
      resolveDetails +
      `Try reinstalling: npm install @code-yeongyu/comment-checker`
  );
}

module.exports = {
  getBinaryName,
  getBinaryPath,
  getPlatformKey,
  PLATFORM_PACKAGES,
};
