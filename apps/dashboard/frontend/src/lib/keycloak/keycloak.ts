// src/keycloak.ts
import Keycloak, { KeycloakInitOptions } from "keycloak-js";

// Reuse across HMR/StrictMode
const g = globalThis as any;

g.__kc ??= new Keycloak({
  url: "http://localhost:8080", // TODO: move to env
  realm: "dev",
  clientId: "portcall",
});

export const kc: Keycloak = g.__kc;

// Cache the single init() promise so subsequent calls reuse it
export function initKeycloak(opts?: KeycloakInitOptions) {
  if (g.__kcInitPromise) return g.__kcInitPromise as Promise<boolean>;

  const defaults: KeycloakInitOptions = {
    onLoad: "check-sso", // don't auto-redirect; lets <SignedOut/> decide
    pkceMethod: "S256",
    checkLoginIframe: false,
  };

  g.__kcInitPromise = kc.init({ ...defaults, ...opts }).catch((err: any) => {
    // reset cache if init failed so you can retry
    g.__kcInitPromise = null;
    throw err;
  });

  return g.__kcInitPromise as Promise<boolean>;
}

// Optional: start one refresh interval (also memoized)
export function startTokenRefresh() {
  console.log("Starting Keycloak token refresh...");

  const interval = setInterval(async () => {
    try {
      console.log("Refreshing Keycloak token...");
      // refresh if token expires in < 60s
      await kc.updateToken(60);
    } catch {
      // if refresh fails, let your app trigger login when needed
      console.warn(
        "Keycloak token refresh failed, user may need to log in again."
      );
    }
  }, 30_000);

  return interval;
}
