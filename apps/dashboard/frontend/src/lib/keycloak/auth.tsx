import type Keycloak from "keycloak-js";
import React, {
  createContext,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import { getKeycloak, initKeycloak, startTokenRefresh } from "./keycloak";

type AuthCtx = {
  isLoaded: boolean;
  isAuthenticated: boolean;
  getToken: () => Promise<string | undefined>;
  token?: string;
  user?: { sub?: string; email?: string; name?: string };
  login: () => void;
  logout: () => void;
};

const Ctx = createContext<AuthCtx | null>(null);

type RuntimeWindow = Window & { __ENV__?: { E2E_MODE?: boolean } };

function isE2EMode() {
  return Boolean((window as RuntimeWindow).__ENV__?.E2E_MODE);
}

function E2EAuthProvider({ children }: { children: React.ReactNode }) {
  const value = useMemo<AuthCtx>(
    () => ({
      isLoaded: true,
      isAuthenticated: true,
      token: "e2e-fake-token",
      user: { sub: "e2e", email: "admin@portcall.internal", name: "E2E" },
      login: () => {},
      logout: () => {},
      getToken: async () => "e2e-fake-token",
    }),
    [],
  );
  return <Ctx.Provider value={value}>{children}</Ctx.Provider>;
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  if (isE2EMode()) {
    return <E2EAuthProvider>{children}</E2EAuthProvider>;
  }

  return <RealAuthProvider>{children}</RealAuthProvider>;
}

function RealAuthProvider({ children }: { children: React.ReactNode }) {
  const [isLoaded, setLoaded] = useState(false);
  const [isAuthenticated, setAuth] = useState(false);
  const [token, setToken] = useState<string | undefined>();
  const [user, setUser] = useState<AuthCtx["user"]>();
  const [kc, setKc] = useState<Keycloak | null>(null);

  useEffect(() => {
    let alive = true;
    let interval: ReturnType<typeof setInterval> | undefined;
    (async () => {
      const ok = await initKeycloak();
      if (!alive) return;
      const keycloak = await getKeycloak();
      setKc(keycloak);
      setAuth(ok);
      setToken(keycloak.token);
      setUser(parseJWT(keycloak.token));
      setLoaded(true);
      interval = await startTokenRefresh(() =>
        keycloak.login({ redirectUri: window.location.href }),
      );
    })();
    return () => {
      alive = false;
      if (interval) clearInterval(interval);
    };
  }, []);

  useEffect(() => {
    if (!kc) return;
    kc.onAuthSuccess = () => setAuth(true);
    kc.onAuthLogout = () => {
      setAuth(false);
      setToken(undefined);
      setUser(undefined);
    };
    kc.onTokenExpired = async () => {
      try {
        const ok = await kc.updateToken(60);
        if (ok) {
          setToken(kc.token);
          setUser(parseJWT(kc.token));
        }
      } catch {
        kc.login({ redirectUri: window.location.href });
      }
    };
  }, [kc]);

  const value = useMemo<AuthCtx>(
    () => ({
      isLoaded,
      isAuthenticated,
      user,
      login: () => kc?.login({ redirectUri: window.location.href }),
      logout: () => kc?.logout({ redirectUri: window.location.origin }),
      getToken: async () => {
        if (!isLoaded || !isAuthenticated || !kc) return;
        await kc.updateToken(60);
        return kc.token;
      },
      token,
    }),
    [isLoaded, isAuthenticated, token, user, kc],
  );

  return <Ctx.Provider value={value}>{children}</Ctx.Provider>;
}

export function useAuth() {
  const ctx = useContext(Ctx);
  if (!ctx) throw new Error("useAuth must be used within <AuthProvider>");
  return ctx;
}
export function SignedIn({ children }: React.PropsWithChildren) {
  const { isLoaded, isAuthenticated } = useAuth();
  if (!isLoaded) return null;
  return isAuthenticated ? <>{children}</> : null;
}
export function SignedOut({ children }: React.PropsWithChildren) {
  const { isLoaded, isAuthenticated } = useAuth();
  if (!isLoaded) return null;
  return !isAuthenticated ? <>{children}</> : null;
}
export function RedirectToSignIn() {
  const { isLoaded, isAuthenticated, login } = useAuth();
  useEffect(() => {
    if (isLoaded && !isAuthenticated) login();
  }, [isLoaded, isAuthenticated, login]);
  return (
    <div className="auth-transition-screen" role="status" aria-live="polite">
      <div className="auth-transition-card">Redirecting to sign in...</div>
    </div>
  );
}

function parseJWT(token?: string): AuthCtx["user"] | undefined {
  if (!token) return undefined;
  const [, payload] = token.split(".");
  if (!payload) return undefined;
  try {
    const parsed = JSON.parse(
      atob(payload.replace(/-/g, "+").replace(/_/g, "/")),
    ) as Record<string, unknown>;
    return {
      sub: typeof parsed.sub === "string" ? parsed.sub : undefined,
      email: typeof parsed.email === "string" ? parsed.email : undefined,
      name: typeof parsed.name === "string" ? parsed.name : undefined,
    };
  } catch {
    return undefined;
  }
}
