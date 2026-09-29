const assert = require("node:assert/strict");
const { readFileSync } = require("node:fs");
const { join } = require("node:path");
const { test } = require("node:test");

const { SUPPORTED_PLATFORMS, getPlatformPackageName } = require("./index");

function readManifest(directory) {
  return JSON.parse(readFileSync(join(__dirname, "..", directory, "package.json"), "utf8"));
}

test("root package depends on exactly one same-version platform package per supported platform", () => {
  // given
  const root = readManifest("comment-checker");

  // when
  const expected = Object.fromEntries(SUPPORTED_PLATFORMS.map((key) => [getPlatformPackageName(key), root.version]));

  // then
  assert.deepEqual(root.optionalDependencies, expected);
  assert.equal(root.files.includes("vendor"), false);
});

for (const key of SUPPORTED_PLATFORMS) {
  test(`${key} package installs only on its own os and cpu and ships only bin/`, () => {
    // given
    const manifest = readManifest(`comment-checker-${key}`);
    const [os, cpu] = key.split("-");

    // when / then
    assert.equal(manifest.name, getPlatformPackageName(key));
    assert.deepEqual(manifest.os, [os]);
    assert.deepEqual(manifest.cpu, [cpu]);
    assert.deepEqual(manifest.files, ["bin"]);
    assert.equal(manifest.version, readManifest("comment-checker").version);
  });
}
