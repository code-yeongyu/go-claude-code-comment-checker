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

function writeFakeBinary(root, relativeDir, binaryName) {
  const binDir = join(root, relativeDir);
  mkdirSync(binDir, { recursive: true });
  const binaryPath = join(binDir, binaryName);
  writeFileSync(binaryPath, "#!/usr/bin/env node\nprocess.exit(0);\n");
  chmodSync(binaryPath, 0o755);
  return binaryPath;
}

function writeFakePlatformPackage(root, platformKey, binaryName) {
  const packageDir = join(root, "node_modules", "@code-yeongyu", `comment-checker-${platformKey}`);
  mkdirSync(packageDir, { recursive: true });
  writeFileSync(join(packageDir, "package.json"), JSON.stringify({ name: `@code-yeongyu/comment-checker-${platformKey}` }));
  return writeFakeBinary(packageDir, "bin", binaryName);
}

function rootPackageDir(root) {
  const baseDir = join(root, "node_modules", "@code-yeongyu", "comment-checker");
  mkdirSync(baseDir, { recursive: true });
  return baseDir;
}

test("lists supported platforms", () => {
  // given
  const platforms = checker.SUPPORTED_PLATFORMS;

  // when
  const resolvedPlatforms = platforms;

  // then
  assert.deepEqual(resolvedPlatforms, [
    "darwin-arm64",
    "darwin-x64",
    "linux-arm64",
    "linux-x64",
    "win32-x64",
  ]);
});

test("resolves bundled vendor binary without postinstall", () => {
  // given
  const root = tempDir();
  try {
    const binaryPath = writeFakeBinary(root, join("vendor", "linux-x64"), "comment-checker");

    // when
    const result = checker.getBinaryPath({
      arch: "x64",
      baseDir: root,
      platform: "linux",
    });

    // then
    assert.equal(realpathSync(result), realpathSync(binaryPath));
  } finally {
    rmSync(root, { recursive: true, force: true });
  }
});

test("prefers the installed platform package over a bundled vendor binary", () => {
  // given
  const root = tempDir();
  try {
    const baseDir = rootPackageDir(root);
    writeFakeBinary(baseDir, join("vendor", "linux-arm64"), "comment-checker");
    const binaryPath = writeFakePlatformPackage(root, "linux-arm64", "comment-checker");

    // when
    const result = checker.getBinaryPath({ arch: "arm64", baseDir, platform: "linux" });

    // then
    assert.equal(realpathSync(result), realpathSync(binaryPath));
  } finally {
    rmSync(root, { recursive: true, force: true });
  }
});

test("resolves the Windows platform package exe", () => {
  // given
  const root = tempDir();
  try {
    const baseDir = rootPackageDir(root);
    const binaryPath = writeFakePlatformPackage(root, "win32-x64", "comment-checker.exe");

    // when
    const result = checker.getBinaryPath({ arch: "x64", baseDir, platform: "win32" });

    // then
    assert.equal(realpathSync(result), realpathSync(binaryPath));
  } finally {
    rmSync(root, { recursive: true, force: true });
  }
});

test("falls back to the postinstall binary when no platform package is installed", () => {
  // given
  const root = tempDir();
  try {
    const baseDir = rootPackageDir(root);
    writeFakePlatformPackage(root, "darwin-x64", "comment-checker");
    const binaryPath = writeFakeBinary(baseDir, "bin", "comment-checker");

    // when
    const result = checker.getBinaryPath({ arch: "arm64", baseDir, platform: "darwin" });

    // then
    assert.equal(realpathSync(result), realpathSync(binaryPath));
  } finally {
    rmSync(root, { recursive: true, force: true });
  }
});

test("reports every layout it looked for when no binary is installed", () => {
  // given
  const root = tempDir();
  try {
    const baseDir = rootPackageDir(root);

    // when / then
    assert.throws(
      () => checker.getBinaryPath({ arch: "x64", baseDir, platform: "linux" }),
      /@code-yeongyu\/comment-checker-linux-x64 package.*vendor.*bin/s
    );
  } finally {
    rmSync(root, { recursive: true, force: true });
  }
});

test("uses exe name for bundled Windows binary", () => {
  // given
  const root = tempDir();
  try {
    const binaryPath = writeFakeBinary(root, join("vendor", "win32-x64"), "comment-checker.exe");

    // when
    const result = checker.getBinaryPath({
      arch: "x64",
      baseDir: root,
      platform: "win32",
    });

    // then
    assert.equal(realpathSync(result), realpathSync(binaryPath));
  } finally {
    rmSync(root, { recursive: true, force: true });
  }
});
