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
