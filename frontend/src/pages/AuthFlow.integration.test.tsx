import { screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import Header from "../components/layout/Header";
import Login from "./Login";
import Register from "./Register";
import { loginUser, registerUser } from "../services/authService";
import {
  adminUser,
  customerUser,
  renderWithProviders,
} from "../test-utils/renderWithProviders";
import type { AuthContextType } from "../hooks/useAuth";

vi.mock("../services/authService", async () => {
  const actual = await vi.importActual<
    typeof import("../services/authService")
  >("../services/authService");

  return {
    ...actual,
    loginUser: vi.fn(),
    registerUser: vi.fn(),
  };
});

const mockedLoginUser = vi.mocked(loginUser);
const mockedRegisterUser = vi.mocked(registerUser);

function renderAuthRoutes(
  initialPath: string,
  options: {
    login?: AuthContextType["login"];
    logout?: AuthContextType["logout"];
  } = {},
) {
  const login = options.login ?? vi.fn<AuthContextType["login"]>();
  const logout = options.logout ?? vi.fn<AuthContextType["logout"]>();
  return renderWithProviders(
    [
      {
        path: "/login",
        element: <Login />,
      },
      {
        path: "/register",
        element: <Register />,
      },
      {
        path: "/dashboard",
        element: <p>Customer dashboard</p>,
      },
      {
        path: "/admin/dashboard",
        element: <p>Admin dashboard</p>,
      },
      {
        path: "/logout-shell",
        element: (
          <Header
            isCollapsed={false}
            cartCount={2}
            onToggleCollapse={() => undefined}
            onOpenMobile={() => undefined}
          />
        ),
      },
    ],
    {
      initialEntries: [initialPath],
      auth: {
        user: customerUser,
        token: null,
        isAuthenticated: false,
        isLoading: false,
        login,
        logout,
      },
    },
  );
}

describe("auth flow", () => {
  beforeEach(() => {
    mockedLoginUser.mockReset();
    mockedRegisterUser.mockReset();
    localStorage.clear();
    sessionStorage.clear();
  });

  it("redirects customer login to customer dashboard", async () => {
    const user = userEvent.setup();
    const login = vi.fn<AuthContextType["login"]>();

    mockedLoginUser.mockResolvedValue({
      accessToken: "customer-token",
      user: customerUser,
    });

    renderAuthRoutes("/login", { login });

    await user.type(screen.getByLabelText(/alamat email/i), customerUser.email);
    await user.type(screen.getByLabelText(/^kata sandi$/i), "password123");
    await user.click(screen.getByRole("button", { name: /^masuk$/i }));

    expect(await screen.findByText("Customer dashboard")).toBeInTheDocument();
    expect(login).toHaveBeenCalledWith("customer-token", customerUser, true);
  });

  it("redirects admin login to admin dashboard", async () => {
    const user = userEvent.setup();
    const login = vi.fn<AuthContextType["login"]>();

    mockedLoginUser.mockResolvedValue({
      accessToken: "admin-token",
      user: adminUser,
    });

    renderAuthRoutes("/login", { login });

    await user.type(screen.getByLabelText(/alamat email/i), adminUser.email);
    await user.type(screen.getByLabelText(/^kata sandi$/i), "password123");
    await user.click(screen.getByRole("button", { name: /^masuk$/i }));

    expect(await screen.findByText("Admin dashboard")).toBeInTheDocument();
    expect(login).toHaveBeenCalledWith("admin-token", adminUser, true);
  });

  it("shows login error for invalid credentials", async () => {
    const user = userEvent.setup();

    mockedLoginUser.mockRejectedValue(
      new Error("Email atau kata sandi tidak sesuai."),
    );

    renderAuthRoutes("/login");

    await user.type(screen.getByLabelText(/alamat email/i), "wrong@test.dev");
    await user.type(screen.getByLabelText(/^kata sandi$/i), "password123");
    await user.click(screen.getByRole("button", { name: /^masuk$/i }));

    expect(
      await screen.findByText("Email atau kata sandi tidak sesuai."),
    ).toBeInTheDocument();
  });

  it("registers valid user and redirects to dashboard", async () => {
    const user = userEvent.setup();
    const login = vi.fn<AuthContextType["login"]>();

    mockedRegisterUser.mockResolvedValue({
      accessToken: "new-customer-token",
      user: customerUser,
    });

    renderAuthRoutes("/register", { login });

    await user.type(screen.getByLabelText(/^nama$/i), customerUser.name);
    await user.type(screen.getByLabelText(/^email$/i), customerUser.email);
    await user.type(screen.getByLabelText(/^password$/i), "Password123!");
    await user.type(
      screen.getByLabelText(/konfirmasi password/i),
      "Password123!",
    );
    await user.click(screen.getByRole("button", { name: /daftar sekarang/i }));

    expect(await screen.findByText("Customer dashboard")).toBeInTheDocument();
    expect(login).toHaveBeenCalledWith(
      "new-customer-token",
      customerUser,
      true,
    );
  });

  it("shows register error for existing email", async () => {
    const user = userEvent.setup();

    mockedRegisterUser.mockRejectedValue(new Error("Email sudah terdaftar."));

    renderAuthRoutes("/register");

    await user.type(screen.getByLabelText(/^nama$/i), customerUser.name);
    await user.type(screen.getByLabelText(/^email$/i), customerUser.email);
    await user.type(screen.getByLabelText(/^password$/i), "Password123!");
    await user.type(
      screen.getByLabelText(/konfirmasi password/i),
      "Password123!",
    );
    await user.click(screen.getByRole("button", { name: /daftar sekarang/i }));

    expect(
      await screen.findByText("Email sudah terdaftar."),
    ).toBeInTheDocument();
  });

  it("logs out from header and navigates to login", async () => {
    const user = userEvent.setup();
    const logout = vi.fn<AuthContextType["logout"]>();

    renderWithProviders(
      [
        {
          path: "/logout-shell",
          element: (
            <Header
              isCollapsed={false}
              cartCount={2}
              onToggleCollapse={() => undefined}
              onOpenMobile={() => undefined}
            />
          ),
        },
        {
          path: "/login",
          element: <p>Login route after logout</p>,
        },
      ],
      {
        initialEntries: ["/logout-shell"],
        auth: {
          user: customerUser,
          token: "customer-token",
          isAuthenticated: true,
          isLoading: false,
          logout,
        },
      },
    );

    await user.click(screen.getByRole("button", { name: /halo/i }));
    await user.click(screen.getByRole("menuitem", { name: /logout/i }));

    await waitFor(() => {
      expect(logout).toHaveBeenCalled();
    });
    expect(
      await screen.findByText("Login route after logout"),
    ).toBeInTheDocument();
  });
});
