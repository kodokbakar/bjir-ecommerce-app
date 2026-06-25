import { render } from "@testing-library/react";
import {
  createMemoryRouter,
  RouterProvider,
  type InitialEntry,
  type RouteObject,
} from "react-router-dom";
import type { ReactElement } from "react";

import { AuthContext, type AuthContextType, type User } from "../hooks/useAuth";
import { ToastProvider } from "../context/ToastContext";
import ToastViewport from "../components/ToastViewport";

export const customerUser: User = {
  id: "user-customer",
  name: "Bintang Customer",
  email: "customer@example.test",
  role: "customer",
};

export const adminUser: User = {
  id: "user-admin",
  name: "Bintang Admin",
  email: "admin@example.test",
  role: "admin",
};

export function createAuthValue(
  overrides: Partial<AuthContextType> = {},
): AuthContextType {
  const user = overrides.user ?? customerUser;
  const token = overrides.token ?? "test-token";

  return {
    user,
    token,
    isAuthenticated: true,
    isLoading: false,
    login: () => undefined,
    logout: () => undefined,
    ...overrides,
  };
}

interface RenderWithProvidersOptions {
  initialEntries?: InitialEntry[];
  auth?: Partial<AuthContextType>;
}

export function renderWithProviders(
  routes: RouteObject[],
  options: RenderWithProvidersOptions = {},
) {
  const router = createMemoryRouter(routes, {
    initialEntries: options.initialEntries ?? ["/"],
  });

  const authValue = createAuthValue(options.auth);

  const result = render(
    <ToastProvider>
      <AuthContext.Provider value={authValue}>
        <RouterProvider router={router} />
        <ToastViewport />
      </AuthContext.Provider>
    </ToastProvider>,
  );

  return {
    router,
    authValue,
    ...result,
  };
}

export function renderPage(
  element: ReactElement,
  options: RenderWithProvidersOptions = {},
) {
  return renderWithProviders(
    [
      {
        path: "*",
        element,
      },
    ],
    options,
  );
}
