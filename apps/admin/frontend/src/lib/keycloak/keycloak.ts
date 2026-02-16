import Keycloak, { KeycloakInitOptions } from "keycloak-js";

const g = globalThis as unknown as {
  __kc?: Keycloak;
  __kcInitPromise?: Promise<boolean>;
};

g.__kc ??= new Keycloak({
  url: import.meta.env.VITE_KEYCLOAK_URL,
  realm: "admin",
  clientId: "portcall-admin",
});

export const kc: Keycloak = g.__kc;

export function initKeycloak(opts?: KeycloakInitOptions) {
  if (g.__kcInitPromise) return g.__kcInitPromise as Promise<boolean>;

  const defaults: KeycloakInitOptions = {
    pkceMethod: "S256",
    checkLoginIframe: false,
    redirectUri: window.location.origin + "/",
  };

  console.log("Keycloak config:", {
    url: import.meta.env.VITE_KEYCLOAK_URL,
    realm: "admin",
    clientId: "portcall-admin",
    redirectUri: defaults.redirectUri,
  });

  // Add timeout to prevent hanging
  const timeoutPromise = new Promise<boolean>((_, reject) => {
    setTimeout(
      () => reject(new Error("Keycloak initialization timeout after 10s")),
      10000,
    );
  });

  g.__kcInitPromise = Promise.race([
    kc.init({ ...defaults, ...opts }),
    timeoutPromise,
  ])
    .then((authenticated) => {
      console.log("Keycloak init result:", authenticated);
      return authenticated;
    })
    .catch((err: unknown) => {
      console.error("Keycloak init error:", err);
      g.__kcInitPromise = undefined;
      throw err;
    });

  return g.__kcInitPromise as Promise<boolean>;
}

export function startTokenRefresh() {
  const interval = setInterval(async () => {
    try {
      await kc.updateToken(60);
    } catch {
      console.warn("Keycloak token refresh failed");
    }
  }, 30_000);

  return interval;
}
