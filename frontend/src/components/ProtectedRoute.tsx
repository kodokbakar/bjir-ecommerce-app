import { Navigate, Outlet } from "react-router-dom";
import { useAuth } from "../hooks/useAuth";
import { C } from "../styles/tokens";

const ProtectedRoute = () => {
    const { isAuthenticated, isLoading } = useAuth();

    if (isLoading) {
        return (
            <div
                style={{
                    display: "flex",
                    minHeight: "100vh",
                    width: "100vw",
                    alignItems: "center",
                    justifyContent: "center",
                    background: C.secondary,
                }}
            >
                <div style={{ textAlign: "center" }}>
                    <p
                        style={{
                            margin: 0,
                            color: C.textMuted,
                            fontSize: 14,
                            fontWeight: 500,
                        }}
                    >
                        Memeriksa sesi...
                    </p>
                </div>
            </div>
        );
    }

    if (!isAuthenticated) {
        return <Navigate to="/login" replace />;
    }

    return <Outlet />;
};

export default ProtectedRoute;