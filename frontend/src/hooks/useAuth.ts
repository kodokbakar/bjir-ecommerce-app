import { createContext, useContext } from "react";

export interface User {
  id: string;
  name: string;
  email: string;
  role?: string;
  is_active?: boolean;
  created_at?: string;
}

export interface AuthContextType {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (token: string, userData: User, remember?: boolean) => void;
  logout: () => void;
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);

  if (context === undefined) {
    throw new Error("useAuth harus digunakan di dalam AuthProvider");
  }

  return context;
};
