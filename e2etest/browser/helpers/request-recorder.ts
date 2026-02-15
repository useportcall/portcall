import { Page } from "@playwright/test";

export function recordDashboardMutations(page: Page) {
  const calls: string[] = [];

  page.on("requestfinished", (request) => {
    const method = request.method();
    if (method !== "POST" && method !== "DELETE") return;
    const url = request.url();
    if (!url.includes("/api/apps/")) return;
    calls.push(`${method} ${url}`);
  });

  return {
    calls,
    hasPath(path: string) {
      return calls.some((entry) => entry.includes(path));
    },
  };
}
