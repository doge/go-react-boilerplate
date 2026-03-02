import { createContext, useContext, useEffect, useMemo, useState } from "react";
import { getSession, login as loginRequest, logout as logoutRequest, refreshAccessToken } from "../lib/api";

type User = {
  uid: string;
  email: string;
};

type AuthContextValue = {
  user: User | null;
  isLoading: boolean;
  login: (username: string, password: string) => Promise<User>;
  logout: () => Promise<void>;
  refreshUser: () => Promise<User | null>;
};

const AuthContext = createContext<AuthContextValue | null>(null);
let bootstrapAuthPromise: Promise<User | null> | null = null;

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    if (!bootstrapAuthPromise) {
      bootstrapAuthPromise = refreshUser();
    }

    void bootstrapAuthPromise.finally(() => setIsLoading(false));
  }, []);

  async function refreshUser(): Promise<User | null> {
    const refreshed = await refreshAccessToken();
    if (!refreshed) {
      setUser(null);
      return null;
    }

    try {
      const session = await getSession();
      setUser(session);
      return session;
    } catch {
      setUser(null);
      return null;
    }
  }

  async function login(username: string, password: string): Promise<User> {
    const loginPayload = await loginRequest(username, password);
    const session = await getSession();
    const nextUser = { uid: session.uid, email: session.email || loginPayload.email };
    setUser(nextUser);
    return nextUser;
  }

  async function logout(): Promise<void> {
    await logoutRequest();
    setUser(null);
  }

  const value = useMemo(
    () => ({
      user,
      isLoading,
      login,
      logout,
      refreshUser,
    }),
    [user, isLoading],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used inside AuthProvider");
  }
  return context;
}
