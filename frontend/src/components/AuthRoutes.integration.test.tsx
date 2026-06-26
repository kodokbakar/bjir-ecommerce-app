import { screen } from "@testing-library/react";
import { Navigate, Outlet } from "react-router-dom";
import { describe, expect, it, vi } from "vitest";

import AdminProtectedRoute from "./AdminProtectedRoute";
import ProtectedRoute from "./ProtectedRoute";
import {
  adminUser,
  customerUser,
  renderWithProviders,
} from "../test-utils/renderWithProviders";

describe("auth route guards", () => {
  it("redirects unauthenticated dashboard access to login", async () => {
    renderWithProviders(
      [
        {
          element: <ProtectedRoute />,
          children: [
            {
              path: "/dashboard",
              element: <p>Private dashboard</p>,
            },
          ],
        },
        {
          path: "/login",
          element: <p>Login page</p>,
        },
      ],
      {
        initialEntries: ["/dashboard"],
        auth: {
          user: null,
          token: null,
          isAuthenticated: false,
          isLoading: false,
          login: vi.fn(),
          logout: vi.fn(),
        },
      },
    );

    expect(await screen.findByText("Login page")).toBeInTheDocument();
  });

  it("redirects customer away from admin routes", async () => {
    renderWithProviders(
      [
        {
          element: <AdminProtectedRoute />,
          children: [
            {
              path: "/admin/dashboard",
              element: <p>Admin dashboard</p>,
            },
          ],
        },
        {
          path: "/dashboard",
          element: <p>Customer dashboard</p>,
        },
      ],
      {
        initialEntries: ["/admin/dashboard"],
        auth: {
          user: customerUser,
          token: "customer-token",
          isAuthenticated: true,
          isLoading: false,
        },
      },
    );

    expect(await screen.findByText("Customer dashboard")).toBeInTheDocument();
  });

  it("allows admin routes for admin user", async () => {
    renderWithProviders(
      [
        {
          element: <AdminProtectedRoute />,
          children: [
            {
              path: "/admin/dashboard",
              element: <p>Admin dashboard</p>,
            },
          ],
        },
      ],
      {
        initialEntries: ["/admin/dashboard"],
        auth: {
          user: adminUser,
          token: "admin-token",
          isAuthenticated: true,
          isLoading: false,
        },
      },
    );

    expect(await screen.findByText("Admin dashboard")).toBeInTheDocument();
  });

  it("redirects root customer to customer dashboard and admin to admin dashboard", async () => {
    const customerResult = renderWithProviders(
      [
        {
          path: "/dashboard",
          element: <p>Customer dashboard</p>,
        },
        {
          path: "/admin/dashboard",
          element: <p>Admin dashboard</p>,
        },
      ],
      {
        initialEntries: ["/dashboard"],
        auth: {
          user: customerUser,
          token: "customer-token",
          isAuthenticated: true,
          isLoading: false,
        },
      },
    );

    expect(await screen.findByText("Customer dashboard")).toBeInTheDocument();
    customerResult.unmount();

    renderWithProviders(
      [
        {
          path: "/dashboard",
          element: <p>Customer dashboard</p>,
        },
        {
          path: "/admin/dashboard",
          element: <p>Admin dashboard</p>,
        },
      ],
      {
        initialEntries: ["/admin/dashboard"],
        auth: {
          user: adminUser,
          token: "admin-token",
          isAuthenticated: true,
          isLoading: false,
        },
      },
    );

    expect(await screen.findByText("Admin dashboard")).toBeInTheDocument();
  });

  it("keeps admin user allowed inside customer protected routes", async () => {
    renderWithProviders(
      [
        {
          element: <ProtectedRoute />,
          children: [
            {
              path: "/dashboard",
              element: <Outlet />,
              children: [
                {
                  index: true,
                  element: <p>Customer route for authenticated user</p>,
                },
              ],
            },
          ],
        },
        {
          path: "/login",
          element: <Navigate to="/login" />,
        },
      ],
      {
        initialEntries: ["/dashboard"],
        auth: {
          user: adminUser,
          token: "admin-token",
          isAuthenticated: true,
          isLoading: false,
        },
      },
    );

    expect(
      await screen.findByText("Customer route for authenticated user"),
    ).toBeInTheDocument();
  });
});
