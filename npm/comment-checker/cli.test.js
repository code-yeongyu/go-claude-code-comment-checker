const assert = require("node:assert/strict");
const {
  chmodSync,
  copyFileSync,
  mkdirSync,
  mkdtempSync,
  rmSync,
  writeFileSync,
} = require("node:fs");
const { tmpdir } = require("node:os");
const { join } = require("node:path");
const { spawnSync } = require("node:child_process");
const { test } = require("node:test");

const checker = require("./index");

function tempDir() {
  return mkdtempSync(join(tmpdir(), "comment-checker-cli-"));
}

function currentPlatformPackageName() {
  return checker.PLATFORM_PACKAGES[checker.getPlatformKey()];
}

test("cli wrapper propagates exit code and output", () => {
  // given
  const root = tempDir();
  try {
    copyFileSync(join(__dirname, "index.js"), join(root, "index.js"));
    copyFileSync(join(__dirname, "cli.js"), join(root, "cli.js"));

    const packageName = currentPlatformPackageName();
    const packageDir = join(root, "node_modules", ...packageName.split("/"));
    const binDir = join(packageDir, "bin");
    const binaryName = process.platform === "win32" ? "comment-checker.exe" : "comment-checker";
    const binaryPath = join(binDir, binaryName);
    mkdirSync(binDir, { recursive: true });
    writeFileSync(join(packageDir, "package.json"), JSON.stringify({ name: packageName }));
    writeFileSync(
      binaryPath,
      "#!/usr/bin/env node\nconsole.log('native stdout ' + process.argv.slice(2).join(' '));\nconsole.error('native stderr');\nprocess.exit(2);\n"
    );
    chmodSync(binaryPath, 0o755);

    // when
    const result = spawnSync(process.execPath, [join(root, "cli.js"), "--probe"], {
      encoding: "utf8",
    });

    // then
    assert.equal(result.status, 2);
    assert.match(result.stdout, /native stdout --probe/);
    assert.match(result.stderr, /native stderr/);
  } finally {
    rmSync(root, { recursive: true, force: true });
  }
});
