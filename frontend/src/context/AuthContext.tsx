import { useEffect, useState, type ReactNode } from "react";

import { AuthContext, type User } from "../hooks/useAuth";
import api, { clearAuthStorage, readStoredToken } from "../services/api";
import { unwrapCurrentUser } from "../services/authService";

function readStoredUser(): User | null {
  const rawUser =
    localStorage.getItem("user") || sessionStorage.getItem("user");

  if (!rawUser) {
    return null;
  }

  try {
    return JSON.parse(rawUser) as User;
  } catch {
    return null;
  }
}

function persistAuth(token: string, user: User, remember: boolean): void {
  const targetStorage = remember ? localStorage : sessionStorage;
  const staleStorage = remember ? sessionStorage : localStorage;

  staleStorage.removeItem("token");
  staleStorage.removeItem("user");

  targetStorage.setItem("token", token);
  targetStorage.setItem("user", JSON.stringify(user));
}

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(true);

  useEffect(() => {
    let isMounted = true;

    async function initializeAuth() {
      const storedToken = readStoredToken();
      const storedUser = readStoredUser();

      if (!storedToken) {
        if (isMounted) {
          setIsLoading(false);
        }

        return;
      }

      try {
        api.defaults.headers.common.Authorization = `Bearer ${storedToken}`;

        const response = await api.get("/v1/me");
        const userData = unwrapCurrentUser(response.data, storedUser);

        if (isMounted) {
          setToken(storedToken);
          setUser(userData);
        }
      } catch (error) {
        console.error("Sesi kedaluwarsa atau token tidak valid:", error);

        clearAuthStorage();
        delete api.defaults.headers.common.Authorization;

        if (isMounted) {
          setToken(null);
          setUser(null);
        }
      } finally {
        if (isMounted) {
          setIsLoading(false);
        }
      }
    }

    initializeAuth();

    return () => {
      isMounted = false;
    };
  }, []);

  const login = (newToken: string, userData: User, remember = true) => {
    persistAuth(newToken, userData, remember);
    api.defaults.headers.common.Authorization = `Bearer ${newToken}`;
    setToken(newToken);
    setUser(userData);
  };

  const logout = () => {
    clearAuthStorage();
    delete api.defaults.headers.common.Authorization;
    setToken(null);
    setUser(null);
  };

  const isAuthenticated = Boolean(token);

  return (
    <AuthContext.Provider
      value={{ user, token, isAuthenticated, isLoading, login, logout }}
    >
      {children}
    </AuthContext.Provider>
  );
};
