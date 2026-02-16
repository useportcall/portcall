import React, {
  createContext,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import { kc, initKeycloak, startTokenRefresh } from "./keycloak";

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

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [isLoaded, setLoaded] = useState(false);
  const [isAuthenticated, setAuth] = useState(false);
  const [token, setToken] = useState<string | undefined>();
  const [user, setUser] = useState<AuthCtx["user"]>();

  useEffect(() => {
    let alive = true;
    let interval: NodeJS.Timeout | undefined;

    (async () => {
      try {
        console.log("Initializing Keycloak...");
        const ok = await initKeycloak();
        console.log("Keycloak initialized:", ok);
        if (!alive) return;
        setAuth(ok);
        setToken(kc.token);
        setUser(parseJWT(kc.token));
        setLoaded(true);
        interval = startTokenRefresh();
      } catch (error) {
        console.error("Keycloak initialization failed:", error);
        if (!alive) return;
        setLoaded(true); // Set loaded even on error so user sees something
      }
    })();
    return () => {
      alive = false;
      if (interval) clearInterval(interval);
    };
  }, []);

  useEffect(() => {
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
        // ignore
      }
    };
  }, []);

  const value = useMemo<AuthCtx>(
    () => ({
      isLoaded,
      isAuthenticated,
      user,
      login: () => kc.login({ redirectUri: window.location.href }),
      logout: () => kc.logout({ redirectUri: window.location.origin }),
      getToken: async () => {
        if (!isLoaded || !isAuthenticated) return;
        await kc.updateToken(60);
        return kc.token;
      },
      token,
    }),
    [isLoaded, isAuthenticated, token, user],
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
  return null;
}

function parseJWT(token?: string) {
  if (!token) return;
  const [, payload] = token.split(".");
  try {
    const json = JSON.parse(
      atob(payload.replace(/-/g, "+").replace(/_/g, "/")),
    );
    return { sub: json.sub, email: json.email, name: json.name };
  } catch {
    // ignore
  }
}
