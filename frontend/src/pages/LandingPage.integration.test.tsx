import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { beforeEach, describe, expect, it, vi } from "vitest";

import App from "../App";
import { AuthContext } from "../hooks/useAuth";
import { listProducts } from "../services/productService";
import { productFixtures, productListResponse } from "../test-utils/fixtures";
import { createAuthValue } from "../test-utils/renderWithProviders";
import { ToastProvider } from "../context/ToastContext";

vi.mock("../services/productService", async () => {
  const actual = await vi.importActual<
    typeof import("../services/productService")
  >("../services/productService");

  return {
    ...actual,
    listProducts: vi.fn(),
  };
});

const mockedListProducts = vi.mocked(listProducts);

function renderAppAtRoot() {
  return render(
    <ToastProvider>
      <AuthContext.Provider
        value={createAuthValue({
          user: null,
          token: null,
          isAuthenticated: false,
          isLoading: false,
        })}
      >
        <MemoryRouter initialEntries={["/"]}>
          <App />
        </MemoryRouter>
      </AuthContext.Provider>
    </ToastProvider>,
  );
}

describe("LandingPage", () => {
  beforeEach(() => {
    mockedListProducts.mockReset();
  });

  it("renders the root landing page instead of redirecting to dashboard", async () => {
    mockedListProducts.mockResolvedValue(productListResponse(productFixtures));

    renderAppAtRoot();

    expect(
      screen.getByRole("navigation", { name: /landing navigation/i }),
    ).toBeInTheDocument();

    expect(
      screen.getByRole("link", { name: /bjir e-commerce/i }),
    ).toHaveAttribute("href", "/");

    expect(
      screen.getByRole("heading", {
        name: /commerce with elbows/i,
      }),
    ).toBeInTheDocument();

    expect(screen.queryByText(/customer dashboard/i)).not.toBeInTheDocument();
  });

  it("renders landing shell sections and CTA links", async () => {
    mockedListProducts.mockResolvedValue(productListResponse(productFixtures));

    renderAppAtRoot();

    expect(
      screen.getByRole("link", { name: /mulai belanja/i }),
    ).toHaveAttribute("href", "/products");

    expect(
      screen.getAllByRole("link", { name: /^daftar$/i })[0],
    ).toHaveAttribute("href", "/register");

    expect(
      screen.getByRole("heading", { name: /fresh from the shelf/i }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole("heading", {
        name: /built like a serious commerce console/i,
      }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole("heading", { name: /need the shelf opened/i }),
    ).toBeInTheDocument();
    expect(screen.getByText(/sharp catalog/i)).toBeInTheDocument();
  });

  it("loads featured products through the mocked product service", async () => {
    mockedListProducts.mockResolvedValue(productListResponse(productFixtures));

    renderAppAtRoot();

    expect(await screen.findByText("Brutal Keyboard")).toBeInTheDocument();
    expect(screen.getByText("Sharp Mouse")).toBeInTheDocument();

    expect(mockedListProducts).toHaveBeenCalledWith({
      page: 1,
      limit: 4,
      sort_by: "created_at",
      sort_order: "desc",
    });
  });

  it("shows a safe featured product fallback when loading fails", async () => {
    mockedListProducts.mockRejectedValue(new Error("Network down"));

    renderAppAtRoot();

    expect(
      await screen.findByText("Featured products could not be loaded."),
    ).toBeInTheDocument();
  });
});
