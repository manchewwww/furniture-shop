import React, {
  createContext,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import { me } from "../api/auth";
import { setAuthToken } from "../api/client";

type User = any;

type AuthContextType = {
  isAuthenticated: boolean;
  user?: User | null;
  refresh: () => Promise<void>;
  logout: () => void;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [user, setUser] = useState<User | null>(null);
  const [ready, setReady] = useState(false);

  const refresh = async () => {
    try {
      const u = await me();
      setUser(u);
    } catch {
      setUser(null);
    } finally {
      setReady(true);
    }
  };

  const logout = () => {
    setAuthToken(undefined);
    setUser(null);
  };

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (token) {
      refresh();
    } else {
      setReady(true);
    }
  }, []);

  const value = useMemo<AuthContextType>(
    () => ({
      isAuthenticated: !!user,
      user,
      refresh,
      logout,
    }),
    [user]
  );

  if (!ready) return null;

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
};
