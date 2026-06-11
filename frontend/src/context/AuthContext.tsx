import { useState, useEffect, type ReactNode } from "react";
import { AuthContext, type User } from "../hooks/useAuth";
import api from "../services/api";

export const AuthProvider = ({ children }: { children: ReactNode }) => {
    const [user, setUser] = useState<User | null>(null);
    const [token, setToken] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState<boolean>(true);

    useEffect(() => {
        const initializeAuth = async () => {
            const storedToken = localStorage.getItem("token");
            const storedUser = localStorage.getItem("user");

            if (storedToken) {
                try {
                    api.defaults.headers.common["Authorization"] = `Bearer ${storedToken}`;
                    const response = await api.get("/v1/auth/me");

                    setToken(storedToken);

                    const userData =
                        response.data?.data?.user ||
                        response.data?.data ||
                        response.data?.user ||
                        (storedUser ? JSON.parse(storedUser) : null);

                    if (!userData) {
                        throw new Error("Data user tidak ditemukan pada response /me");
                    }

                    localStorage.setItem("user", JSON.stringify(userData));
                    setUser(userData);
                } catch (error) {
                    console.error("Sesi kedaluwarsa atau token tidak valid:", error);
                    localStorage.removeItem("token");
                    localStorage.removeItem("user");
                    delete api.defaults.headers.common["Authorization"];
                    setToken(null);
                    setUser(null);
                }
            }
            setIsLoading(false);
        };

        initializeAuth();
    }, []);

    const login = (newToken: string, userData: User) => {
        localStorage.setItem("token", newToken);
        localStorage.setItem("user", JSON.stringify(userData));
        api.defaults.headers.common["Authorization"] = `Bearer ${newToken}`;
        setToken(newToken);
        setUser(userData);
    };

    const logout = () => {
        localStorage.removeItem("token");
        localStorage.removeItem("user");
        delete api.defaults.headers.common["Authorization"];
        setToken(null);
        setUser(null);
    };

    const isAuthenticated = !!token;

    return (
        <AuthContext.Provider value={{ user, token, isAuthenticated, isLoading, login, logout }}>
            {children}
        </AuthContext.Provider>
    );
};