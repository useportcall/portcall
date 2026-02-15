// Shuts down the Go browser harness subprocess.
import * as fs from "fs";
import * as path from "path";

const PID_PATH = path.join(__dirname, ".harness-pid");

export default async function globalTeardown() {
  const child = (globalThis as any).__harnessChild;
  if (child && !child.killed) {
    child.kill("SIGTERM");
    await new Promise<void>((r) => child.on("exit", r));
  }

  // Also try PID file fallback
  if (fs.existsSync(PID_PATH)) {
    const pid = parseInt(fs.readFileSync(PID_PATH, "utf8"), 10);
    try {
      process.kill(pid, "SIGTERM");
    } catch {
      /* already exited */
    }
  }
}
