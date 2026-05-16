#!/usr/bin/env node

const { spawn } = require("node:child_process");
const { getBinaryPath } = require("./index");

function run(argv = process.argv.slice(2), options = {}) {
  const exit = options.exit || process.exit;
  const stderr = options.stderr || process.stderr;
  const spawnCommand = options.spawn || spawn;
  let binaryPath;

  try {
    binaryPath = (options.getBinaryPath || getBinaryPath)();
  } catch (error) {
    stderr.write(`[comment-checker] ${error.message}\n`);
    exit(1);
    return;
  }

  const child = spawnCommand(binaryPath, argv, { stdio: "inherit" });

  child.on("error", (error) => {
    stderr.write(`[comment-checker] Failed to start binary: ${error.message}\n`);
    exit(1);
  });

  child.on("exit", (code, signal) => {
    if (signal) {
      process.kill(process.pid, signal);
      return;
    }

    exit(code === null ? 1 : code);
  });
}

if (require.main === module) {
  run();
}

module.exports = { run };
