import { Navigate, Outlet } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

const ProtectedRoute = () => {
    const { isAuthenticated, isLoading } = useAuth();


    if (isLoading) {
        return (
        <div className="flex h-screen w-screen items-center justify-center bg-gray-50">
            <div className="text-center">
            <p className="text-gray-600 font-medium">Memeriksa sesi...</p>
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
