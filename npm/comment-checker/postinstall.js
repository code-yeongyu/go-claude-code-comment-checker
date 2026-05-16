#!/usr/bin/env node

const {
  chmodSync,
  createWriteStream,
  existsSync,
  mkdirSync,
  rmSync,
  unlinkSync,
} = require("node:fs");
const http = require("node:http");
const https = require("node:https");
const { join } = require("node:path");
const { spawnSync } = require("node:child_process");
const { pipeline } = require("node:stream/promises");

const REPO = "code-yeongyu/go-claude-code-comment-checker";
const MAX_REDIRECTS = 5;

const PLATFORM_MAP = {
  "darwin-arm64": { os: "darwin", arch: "arm64", ext: "tar.gz" },
  "darwin-x64": { os: "darwin", arch: "amd64", ext: "tar.gz" },
  "linux-arm64": { os: "linux", arch: "arm64", ext: "tar.gz" },
  "linux-x64": { os: "linux", arch: "amd64", ext: "tar.gz" },
  "win32-x64": { os: "windows", arch: "amd64", ext: "zip" },
};

function getPlatformKey(platform = process.platform, arch = process.arch) {
  return `${platform}-${arch === "x64" ? "x64" : arch}`;
}

function getPackageVersion() {
  return require("./package.json").version;
}

function requestUrl(url) {
  return new Promise((resolve, reject) => {
    const client = url.startsWith("https:") ? https : http;
    const request = client.get(url, resolve);
    request.on("error", reject);
  });
}

async function downloadFile(url, dest, options = {}) {
  const request = options.request || requestUrl;
  await downloadFileWithRedirects(url, dest, request, 0);
}

async function downloadFileWithRedirects(url, dest, request, redirectCount) {
  if (redirectCount > MAX_REDIRECTS) {
    throw new Error(`Too many redirects while downloading ${url}`);
  }

  const response = await request(url);
  const statusCode = response.statusCode || 0;

  if ([301, 302, 303, 307, 308].includes(statusCode)) {
    const location = response.headers.location;
    response.resume();
    if (!location) {
      throw new Error(`Redirect without location while downloading ${url}`);
    }
    const redirectUrl = new URL(location, url).toString();
    await downloadFileWithRedirects(redirectUrl, dest, request, redirectCount + 1);
    return;
  }

  if (statusCode !== 200) {
    response.resume();
    rmSync(dest, { force: true });
    throw new Error(`Failed to download ${url}: HTTP ${statusCode}`);
  }

  try {
    await pipeline(response, createWriteStream(dest));
  } catch (error) {
    rmSync(dest, { force: true });
    throw error;
  }
}

function runCommand(command, args, options = { stdio: "inherit" }) {
  const result = spawnSync(command, args, options);

  if (result.error) {
    throw result.error;
  }

  if (result.status !== 0) {
    throw new Error(`Command failed with exit ${result.status}: ${command} ${args.join(" ")}`);
  }
}

function extractTarGz(archivePath, destDir, options = {}) {
  const runner = options.runCommand || runCommand;
  runner("tar", ["-xzf", archivePath, "-C", destDir], { stdio: "inherit" });
}

function extractZip(archivePath, destDir, options = {}) {
  const runner = options.runCommand || runCommand;
  const platform = options.platform || process.platform;

  if (platform === "win32") {
    runner(
      "powershell",
      [
        "-NoProfile",
        "-NonInteractive",
        "-ExecutionPolicy",
        "Bypass",
        "-Command",
        "Expand-Archive",
        "-LiteralPath",
        archivePath,
        "-DestinationPath",
        destDir,
        "-Force",
      ],
      { stdio: "inherit" }
    );
    return;
  }

  runner("unzip", ["-o", archivePath, "-d", destDir], { stdio: "inherit" });
}

async function install(options = {}) {
  const platform = options.platform || process.platform;
  const arch = options.arch || process.arch;
  const platformKey = getPlatformKey(platform, arch);
  const platformInfo = PLATFORM_MAP[platformKey];
  const log = options.log || console.log;
  const warn = options.warn || console.warn;

  if (!platformInfo) {
    warn(`[comment-checker] Warning: Unsupported platform ${platform}-${arch}`);
    warn(`[comment-checker] Supported platforms: ${Object.keys(PLATFORM_MAP).join(", ")}`);
    return;
  }

  const baseDir = options.baseDir || __dirname;
  const binDir = join(baseDir, "bin");
  const binaryName = platform === "win32" ? "comment-checker.exe" : "comment-checker";
  const bundledBinaryPath = join(baseDir, "vendor", platformKey, binaryName);
  const binaryPath = join(binDir, binaryName);

  if (existsSync(bundledBinaryPath)) {
    log(`[comment-checker] Bundled ${platformKey} binary already exists at ${bundledBinaryPath}`);
    return;
  }

  if (existsSync(binaryPath)) {
    log(`[comment-checker] Binary already exists at ${binaryPath}`);
    return;
  }

  const version = (options.getPackageVersion || getPackageVersion)();
  const { os, arch: releaseArch, ext } = platformInfo;
  const assetName = `comment-checker_v${version}_${os}_${releaseArch}.${ext}`;
  const downloadUrl = `https://github.com/${REPO}/releases/download/v${version}/${assetName}`;

  log(`[comment-checker] Downloading ${platformKey} binary from GitHub Releases...`);
  log(`[comment-checker] URL: ${downloadUrl}`);

  const mkdir = options.mkdir || mkdirSync;
  mkdir(binDir, { recursive: true });

  const archivePath = join(binDir, assetName);
  const downloader = options.downloadFile || downloadFile;

  await downloader(downloadUrl, archivePath);

  if (ext === "tar.gz") {
    (options.extractTarGz || extractTarGz)(archivePath, binDir);
  } else if (ext === "zip") {
    (options.extractZip || extractZip)(archivePath, binDir, { platform });
  }

  if (existsSync(archivePath)) {
    unlinkSync(archivePath);
  }

  if (platform !== "win32" && existsSync(binaryPath)) {
    chmodSync(binaryPath, 0o755);
  }

  log(`[comment-checker] Successfully installed ${platformKey} binary`);
}

async function main(options = {}) {
  const stderr = options.stderr || console;
  const exit = options.exit || process.exit;

  try {
    await install(options);
  } catch (error) {
    stderr.error(`[comment-checker] Failed to install binary: ${error.message}`);
    stderr.error("[comment-checker] You may need to download it manually from:");
    stderr.error(
      `[comment-checker] https://github.com/${REPO}/releases/tag/v${
        (options.getPackageVersion || getPackageVersion)()
      }`
    );
    exit(1);
  }
}

if (require.main === module) {
  main();
}

module.exports = {
  downloadFile,
  extractTarGz,
  extractZip,
  getPackageVersion,
  getPlatformKey,
  install,
  main,
  runCommand,
};
