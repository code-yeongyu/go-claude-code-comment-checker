#!/usr/bin/env node

const { existsSync, mkdirSync, createWriteStream, chmodSync, unlinkSync } = require("fs");
const { join } = require("path");
const { execSync } = require("child_process");
const https = require("https");
const http = require("http");

const REPO = "code-yeongyu/go-claude-code-comment-checker";

const PLATFORM_MAP = {
  "darwin-arm64": { os: "darwin", arch: "arm64", ext: "tar.gz" },
  "darwin-x64": { os: "darwin", arch: "amd64", ext: "tar.gz" },
  "linux-arm64": { os: "linux", arch: "arm64", ext: "tar.gz" },
  "linux-x64": { os: "linux", arch: "amd64", ext: "tar.gz" },
  "win32-x64": { os: "windows", arch: "amd64", ext: "zip" },
};

function getPlatformKey() {
  const platform = process.platform;
  const arch = process.arch;
  return `${platform}-${arch === "x64" ? "x64" : arch}`;
}

function getPackageVersion() {
  const pkg = require("./package.json");
  return pkg.version;
}

function downloadFile(url, dest) {
  return new Promise((resolve, reject) => {
    const file = createWriteStream(dest);
    const protocol = url.startsWith("https") ? https : http;

    const request = protocol.get(url, (response) => {
      if (response.statusCode === 302 || response.statusCode === 301) {
        // Follow redirect
        downloadFile(response.headers.location, dest).then(resolve).catch(reject);
        return;
      }

      if (response.statusCode !== 200) {
        reject(new Error(`Failed to download: ${response.statusCode}`));
        return;
      }

      response.pipe(file);
      file.on("finish", () => {
        file.close();
        resolve();
      });
    });

    request.on("error", (err) => {
      unlinkSync(dest);
      reject(err);
    });

    file.on("error", (err) => {
      unlinkSync(dest);
      reject(err);
    });
  });
}

async function extractTarGz(archivePath, destDir) {
  execSync(`tar -xzf "${archivePath}" -C "${destDir}"`, { stdio: "pipe" });
}

async function extractZip(archivePath, destDir) {
  if (process.platform === "win32") {
    execSync(`powershell -command "Expand-Archive -Path '${archivePath}' -DestinationPath '${destDir}' -Force"`, { stdio: "pipe" });
  } else {
    execSync(`unzip -o "${archivePath}" -d "${destDir}"`, { stdio: "pipe" });
  }
}

async function main() {
  const platformKey = getPlatformKey();
  const platformInfo = PLATFORM_MAP[platformKey];

  if (!platformInfo) {
    console.warn(
      `[comment-checker] Warning: Unsupported platform ${process.platform}-${process.arch}`
    );
    console.warn(`[comment-checker] Supported platforms: ${Object.keys(PLATFORM_MAP).join(", ")}`);
    process.exit(0);
  }

  const binDir = join(__dirname, "bin");
  const binaryName = process.platform === "win32" ? "comment-checker.exe" : "comment-checker";
  const binaryPath = join(binDir, binaryName);

  // Skip if binary already exists
  if (existsSync(binaryPath)) {
    console.log(`[comment-checker] Binary already exists at ${binaryPath}`);
    return;
  }

  const version = getPackageVersion();
  const { os, arch, ext } = platformInfo;
  const assetName = `comment-checker_v${version}_${os}_${arch}.${ext}`;
  const downloadUrl = `https://github.com/${REPO}/releases/download/v${version}/${assetName}`;

  console.log(`[comment-checker] Downloading ${platformKey} binary from GitHub Releases...`);
  console.log(`[comment-checker] URL: ${downloadUrl}`);

  if (!existsSync(binDir)) {
    mkdirSync(binDir, { recursive: true });
  }

  const archivePath = join(binDir, assetName);

  try {
    await downloadFile(downloadUrl, archivePath);

    if (ext === "tar.gz") {
      await extractTarGz(archivePath, binDir);
    } else if (ext === "zip") {
      await extractZip(archivePath, binDir);
    }

    // Clean up archive
    if (existsSync(archivePath)) {
      unlinkSync(archivePath);
    }

    // Ensure executable (not needed on Windows)
    if (process.platform !== "win32" && existsSync(binaryPath)) {
      chmodSync(binaryPath, 0o755);
    }

    console.log(`[comment-checker] Successfully installed ${platformKey} binary`);
  } catch (err) {
    console.error(`[comment-checker] Failed to download binary: ${err.message}`);
    console.error(`[comment-checker] You may need to download it manually from:`);
    console.error(`[comment-checker] https://github.com/${REPO}/releases/tag/v${version}`);
    process.exit(0); // Don't fail installation
  }
}

main();
