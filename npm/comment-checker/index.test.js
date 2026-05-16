const assert = require("node:assert/strict");
const {
  chmodSync,
  mkdirSync,
  mkdtempSync,
  realpathSync,
  rmSync,
  writeFileSync,
} = require("node:fs");
const { tmpdir } = require("node:os");
const { join } = require("node:path");
const { test } = require("node:test");

const checker = require("./index");

function tempDir() {
  return mkdtempSync(join(tmpdir(), "comment-checker-index-"));
}

function writeFakePlatformPackage(root, packageName, binaryName) {
  const packageDir = join(root, "node_modules", ...packageName.split("/"));
  const binDir = join(packageDir, "bin");
  mkdirSync(binDir, { recursive: true });
  writeFileSync(join(packageDir, "package.json"), JSON.stringify({ name: packageName }));
  const binaryPath = join(binDir, binaryName);
  writeFileSync(binaryPath, "#!/usr/bin/env node\nprocess.exit(0);\n");
  chmodSync(binaryPath, 0o755);
  return binaryPath;
}

test("maps supported platforms to code-yeongyu packages", () => {
  // given
  const packages = checker.PLATFORM_PACKAGES;

  // when
  const resolvedPackages = packages;

  // then
  assert.equal(resolvedPackages["darwin-arm64"], "@code-yeongyu/comment-checker-darwin-arm64");
  assert.equal(resolvedPackages["darwin-x64"], "@code-yeongyu/comment-checker-darwin-x64");
  assert.equal(resolvedPackages["linux-arm64"], "@code-yeongyu/comment-checker-linux-arm64");
  assert.equal(resolvedPackages["linux-x64"], "@code-yeongyu/comment-checker-linux-x64");
  assert.equal(resolvedPackages["win32-x64"], "@code-yeongyu/comment-checker-win32-x64");
});

test("resolves platform package without postinstall", () => {
  // given
  const root = tempDir();
  try {
    const binaryPath = writeFakePlatformPackage(
      root,
      "@code-yeongyu/comment-checker-linux-x64",
      "comment-checker"
    );

    // when
    const result = checker.getBinaryPath({
      arch: "x64",
      platform: "linux",
      resolveFrom: root,
    });

    // then
    assert.equal(result, realpathSync(binaryPath));
  } finally {
    rmSync(root, { recursive: true, force: true });
  }
});

test("uses exe name for Windows package resolution", () => {
  // given
  const root = tempDir();
  try {
    const binaryPath = writeFakePlatformPackage(
      root,
      "@code-yeongyu/comment-checker-win32-x64",
      "comment-checker.exe"
    );

    // when
    const result = checker.getBinaryPath({
      arch: "x64",
      platform: "win32",
      resolveFrom: root,
    });

    // then
    assert.equal(result, realpathSync(binaryPath));
  } finally {
    rmSync(root, { recursive: true, force: true });
  }
});
