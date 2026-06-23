import { Navigate, Outlet } from "react-router-dom";

import AuthLoading from "./auth/AuthLoading";
import { useAuth } from "../hooks/useAuth";

function ProtectedRoute() {
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    return <AuthLoading />;
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return <Outlet />;
}

export default ProtectedRoute;
