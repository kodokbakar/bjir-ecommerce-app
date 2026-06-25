import type { User } from "../hooks/useAuth";

export const ADMIN_DASHBOARD_PATH = "/admin/dashboard";
export const CUSTOMER_DASHBOARD_PATH = "/dashboard";

export function isAdminRole(role?: string | null): boolean {
  const normalizedRole = role?.trim().toLowerCase();

  return normalizedRole === "admin" || normalizedRole === "superadmin";
}

export function getDashboardPath(user?: Pick<User, "role"> | null): string {
  return isAdminRole(user?.role)
    ? ADMIN_DASHBOARD_PATH
    : CUSTOMER_DASHBOARD_PATH;
}
