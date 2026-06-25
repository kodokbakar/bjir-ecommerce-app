import { Navigate } from "react-router-dom";

import AuthLoading from "./auth/AuthLoading";
import { useAuth } from "../hooks/useAuth";
import { getDashboardPath } from "../utils/authRouting";

function DashboardRedirect() {
  const { user, isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    return <AuthLoading />;
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return <Navigate to={getDashboardPath(user)} replace />;
}

export default DashboardRedirect;
