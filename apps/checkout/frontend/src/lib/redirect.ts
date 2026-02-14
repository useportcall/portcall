/**
 * Validates and performs a safe redirect. Only allows relative paths
 * or same-origin URLs to prevent open-redirect attacks.
 */
export function safeRedirect(url: string | undefined, fallback = "/success") {
  const target = url || fallback;

  // Allow relative paths starting with /
  if (target.startsWith("/") && !target.startsWith("//")) {
    window.location.href = target;
    return;
  }

  // Allow same-origin absolute URLs only
  try {
    const parsed = new URL(target);
    if (parsed.origin === window.location.origin) {
      window.location.href = target;
      return;
    }
  } catch {
    // Malformed URL — fall through to fallback
  }

  window.location.href = fallback;
}
