const assert = require("node:assert/strict");
const { mkdtempSync, readFileSync, rmSync } = require("node:fs");
const { tmpdir } = require("node:os");
const { join } = require("node:path");
const { Readable } = require("node:stream");
const { test } = require("node:test");

const postinstall = require("./postinstall");

function tempDir() {
  return mkdtempSync(join(tmpdir(), "comment-checker-postinstall-"));
}

test("download waits for archive bytes before resolving", async () => {
  // given
  const root = tempDir();
  try {
    const archivePath = join(root, "archive.bin");
    let requestedUrl = "";
    const request = async (url) => {
      requestedUrl = url;
      const response = Readable.from(["first-", "second"]);
      response.statusCode = 200;
      response.headers = {};
      return response;
    };

    // when
    await postinstall.downloadFile("https://example.test/archive.bin", archivePath, { request });

    // then
    assert.equal(requestedUrl, "https://example.test/archive.bin");
    assert.equal(readFileSync(archivePath, "utf8"), "first-second");
  } finally {
    rmSync(root, { recursive: true, force: true });
  }
});

test("extractZip failure surfaces inherited stdio command failure", () => {
  // given
  const calls = [];
  const runCommand = (command, args, options) => {
    calls.push({ command, args, options });
    throw new Error("powershell failed");
  };

  // when
  assert.throws(
    () =>
      postinstall.extractZip("C:\\tmp\\archive.zip", "C:\\tmp\\bin", {
        platform: "win32",
        runCommand,
      }),
    /powershell failed/
  );

  // then
  assert.equal(calls.length, 1);
  assert.equal(calls[0].command, "powershell");
  assert.deepEqual(calls[0].options, { stdio: "inherit" });
  assert.ok(calls[0].args.includes("Expand-Archive"));
});

test("main exits non-zero when installation fails", async () => {
  // given
  const root = tempDir();
  const exits = [];
  const errors = [];

  try {
    // when
    await postinstall.main({
      arch: "x64",
      baseDir: root,
      downloadFile: async () => {
        throw new Error("download failed");
      },
      exit: (code) => {
        exits.push(code);
      },
      extractTarGz: () => {},
      extractZip: () => {},
      getPackageVersion: () => "0.7.0",
      log: () => {},
      mkdir: () => {},
      platform: "linux",
      stderr: {
        error: (message) => {
          errors.push(message);
        },
      },
    });
  } finally {
    rmSync(root, { recursive: true, force: true });
  }

  // then
  assert.deepEqual(exits, [1]);
  assert.match(errors.join("\n"), /Failed to install binary/);
  assert.match(errors.join("\n"), /download failed/);
});
