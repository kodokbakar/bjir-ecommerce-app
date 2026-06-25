import { Navigate, Outlet } from "react-router-dom";

import AuthLoading from "./auth/AuthLoading";
import { useAuth } from "../hooks/useAuth";
import { isAdminRole } from "../utils/authRouting";

function AdminProtectedRoute() {
  const { isAuthenticated, isLoading, user } = useAuth();

  if (isLoading) {
    return <AuthLoading />;
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  if (!isAdminRole(user?.role)) {
    return <Navigate to="/dashboard" replace />;
  }

  return <Outlet />;
}

export default AdminProtectedRoute;
