// Launches the Go browser-harness test as a subprocess.
// Reads a single JSON line from stdout with server URLs.
// Stores the config in a temp file for tests to read.
import { spawn, ChildProcess } from "child_process";
import * as fs from "fs";
import * as path from "path";

const CONFIG_PATH = path.join(__dirname, ".harness-config.json");
const PID_PATH = path.join(__dirname, ".harness-pid");

export default async function globalSetup() {
  const root = path.resolve(__dirname, "../../..");

  const child: ChildProcess = spawn(
    "go",
    [
      "test",
      "-v",
      "-run",
      "TestBrowserHarness",
      "-count=1",
      "-timeout=300s",
      "./e2etest/",
    ],
    {
      cwd: root,
      stdio: ["ignore", "pipe", "pipe"],
      env: { ...process.env, BROWSER_HARNESS: "1", GOCACHE: "/tmp/go-build-cache" },
    },
  );

  // Collect stderr for debugging
  let stderr = "";
  child.stderr?.on("data", (d: Buffer) => {
    stderr += d.toString();
  });

  const config = await new Promise<string>((resolve, reject) => {
    let buf = "";
    child.stdout?.on("data", (d: Buffer) => {
      buf += d.toString();
      // Find the JSON config line (starts with '{')
      for (const line of buf.split("\n")) {
        const trimmed = line.trim();
        if (trimmed.startsWith("{") && trimmed.includes("dashboard_url")) {
          resolve(trimmed);
          return;
        }
      }
    });
    child.on("error", reject);
    child.on("exit", (code: number) => {
      if (!buf.includes("dashboard_url")) {
        reject(new Error(`Harness exited (${code}): ${stderr.slice(-500)}`));
      }
    });
    setTimeout(() => reject(new Error("Harness startup timed out")), 120_000);
  });

  // Validate JSON
  JSON.parse(config);
  fs.writeFileSync(CONFIG_PATH, config);
  fs.writeFileSync(PID_PATH, String(child.pid));

  // Keep reference so Node doesn't kill it
  (globalThis as any).__harnessChild = child;
}
