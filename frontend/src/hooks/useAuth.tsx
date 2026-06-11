import { createContext, useContext } from "react";

// Pindahkan interface ke sini agar bersih
export interface User {
    id: string;
    name: string;
    email: string;
}

export interface AuthContextType {
    user: User | null;
    token: string | null;
    isAuthenticated: boolean;
    isLoading: boolean;
    login: (token: string, userData: User) => void;
    logout: () => void;
}


export const AuthContext = createContext<AuthContextType | undefined>(undefined);

// Ekspor Hook
export const useAuth = () => {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error("useAuth harus digunakan di dalam AuthProvider");
    }
    return context;
};