import React, {
    createContext,
    useContext,
    useState,
    useEffect,
    type ReactNode,
} from "react";
// import api from "../services/api";


interface User {
    id: string;
    name: string;
    email: string;

}

interface AuthContextType {
    user: User | null;
    token: string | null;
    isAuthenticated: boolean;
    isLoading: boolean;
    login: (token: string, userData: User) => void;
    logout: () => void;
}


const AuthContext = createContext<AuthContextType | undefined>(undefined);


export const AuthProvider = ({ children }: { children: ReactNode }) => {
    const [user, setUser] = useState<User | null>(null);
    const [token, setToken] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState<boolean>(true);


    useEffect(() => {
    const initializeAuth = async () => {
        const storedToken = localStorage.getItem("token");
        const storedUser = localStorage.getItem("user");

        if (storedToken && storedUser) {
        try {
            setToken(storedToken);
            setUser(JSON.parse(storedUser));

          // await api.get('/auth/me');
            } catch (error) {
                console.error("Gagal memulihkan sesi:", error);
                localStorage.removeItem("token");
                localStorage.removeItem("user");
            }
        }
    setIsLoading(false);
    };

    initializeAuth();
    }, []);


    const login = (newToken: string, userData: User) => {
    localStorage.setItem("token", newToken);
    localStorage.setItem("user", JSON.stringify(userData));
    setToken(newToken);
    setUser(userData);
    };


    const logout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    setToken(null);
    setUser(null);
    };

    const isAuthenticated = !!token;

    return (
    <AuthContext.Provider
        value={{ user, token, isAuthenticated, isLoading, login, logout }}
    >
        {children}
    </AuthContext.Provider>
    );
};

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (context === undefined) {
    throw new Error("useAuth harus digunakan di dalam AuthProvider");
    }
    return context;
};
