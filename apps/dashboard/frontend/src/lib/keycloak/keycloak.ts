// src/keycloak.ts
import Keycloak, { KeycloakInitOptions } from "keycloak-js";

// Reuse across HMR/StrictMode
type KeycloakGlobals = typeof globalThis & {
  __kc?: Keycloak;
  __kcInitPromise?: Promise<boolean> | null;
};
const g = globalThis as KeycloakGlobals;

// Initialize Keycloak instance lazily after config is loaded
let keycloakInstance: Keycloak | null = null;

export async function getKeycloak(): Promise<Keycloak> {
  if (keycloakInstance) {
    return keycloakInstance;
  }

  // Fetch runtime config from backend
  const response = await fetch("/api/config");
  const config = await response.json();
  const keycloakUrl =
    config.data?.keycloak_url ||
    import.meta.env.VITE_KEYCLOAK_URL ||
    "http://localhost:8090";

  keycloakInstance = new Keycloak({
    url: keycloakUrl,
    realm: "dev",
    clientId: "portcall",
  });

  return keycloakInstance;
}

// Cache the single init() promise so subsequent calls reuse it
export async function initKeycloak(opts?: KeycloakInitOptions) {
  if (g.__kcInitPromise) return g.__kcInitPromise as Promise<boolean>;

  const defaults: KeycloakInitOptions = {
    onLoad: "check-sso", // don't auto-redirect; lets <SignedOut/> decide
    pkceMethod: "S256",
    checkLoginIframe: false,
  };

  const kc = await getKeycloak();
  g.__kc = kc;

  g.__kcInitPromise = kc.init({ ...defaults, ...opts }).catch((err: unknown) => {
    // reset cache if init failed so you can retry
    g.__kcInitPromise = null;
    throw err;
  });

  return g.__kcInitPromise as Promise<boolean>;
}

// Optional: start one refresh interval (also memoized)
export async function startTokenRefresh(onFail?: () => void) {
  console.log("Starting Keycloak token refresh...");
  const kc = await getKeycloak();

  const interval = setInterval(async () => {
    try {
      console.log("Refreshing Keycloak token...");
      // refresh if token expires in < 60s
      await kc.updateToken(60);
    } catch {
      // if refresh fails, let your app trigger login when needed
      console.warn(
        "Keycloak token refresh failed, user may need to log in again.",
      );
      onFail?.();
    }
  }, 30_000);

  return interval;
}
