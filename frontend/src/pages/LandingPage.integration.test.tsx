import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter } from "react-router-dom";
import { beforeEach, describe, expect, it, vi } from "vitest";

import App from "../App";
import { AuthContext, type AuthContextType } from "../hooks/useAuth";
import { listProducts } from "../services/productService";
import { productFixtures, productListResponse } from "../test-utils/fixtures";
import {
  adminUser,
  createAuthValue,
  customerUser,
} from "../test-utils/renderWithProviders";
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

function renderAppAtRoot(authOverrides: Partial<AuthContextType> = {}) {
  return render(
    <ToastProvider>
      <AuthContext.Provider
        value={createAuthValue({
          user: null,
          token: null,
          isAuthenticated: false,
          isLoading: false,
          ...authOverrides,
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
      screen.getByRole("link", { name: /daftar sekarang/i }),
    ).toHaveAttribute("href", "/register");

    expect(
      screen.getByRole("heading", { name: /produk unggulan terbaru/i }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole("heading", {
        name: /belanja cepat, aman, tanpa drama/i,
      }),
    ).toBeInTheDocument();

    expect(screen.getByText(/etalase digital yang rapi/i)).toBeInTheDocument();
    expect(
      screen.getByRole("heading", { name: /pengiriman cepat/i }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole("heading", { name: /pembayaran aman/i }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole("heading", { name: /produk original/i }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole("heading", { name: /garansi uang kembali/i }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole("heading", { name: /need the shelf opened/i }),
    ).toBeInTheDocument();
    expect(screen.getByText(/belanja cepat, stok jelas/i)).toBeInTheDocument();

    expect(
      screen.getByRole("navigation", { name: /footer quick links/i }),
    ).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /^produk$/i })).toHaveAttribute(
      "href",
      "/products",
    );
    expect(screen.getByRole("link", { name: /pesanan saya/i })).toHaveAttribute(
      "href",
      "/orders",
    );
    expect(screen.getByRole("link", { name: /keranjang/i })).toHaveAttribute(
      "href",
      "/cart",
    );

    expect(screen.getByRole("link", { name: /github/i })).toHaveAttribute(
      "href",
      "https://github.com/kodokbakar/bjir-ecommerce-app",
    );
    expect(screen.getByRole("link", { name: /email/i })).toHaveAttribute(
      "href",
      "mailto:support@bjir-commerce.test",
    );
    expect(
      screen.getByText(`© ${new Date().getFullYear()} Bjir E-Commerce`),
    ).toBeInTheDocument();
  });

  it("toggles the mobile landing navigation menu", async () => {
    const user = userEvent.setup();
    mockedListProducts.mockResolvedValue(productListResponse(productFixtures));

    renderAppAtRoot();

    const menuButton = screen.getByRole("button", {
      name: /buka menu navigasi landing/i,
    });

    expect(menuButton).toHaveAttribute("aria-expanded", "false");

    await user.click(menuButton);

    expect(
      screen.getByRole("button", {
        name: /tutup menu navigasi landing/i,
      }),
    ).toHaveAttribute("aria-expanded", "true");

    expect(
      screen.getByRole("navigation", { name: /mobile landing navigation/i }),
    ).toBeInTheDocument();

    await user.click(
      screen.getByRole("button", {
        name: /tutup menu navigasi landing/i,
      }),
    );

    expect(
      screen.getByRole("button", {
        name: /buka menu navigasi landing/i,
      }),
    ).toHaveAttribute("aria-expanded", "false");
  });

  it("shows hero loading placeholders while auth state is loading", () => {
    mockedListProducts.mockResolvedValue(productListResponse(productFixtures));

    renderAppAtRoot({
      isLoading: true,
    });

    expect(screen.getByLabelText(/loading hero actions/i)).toBeInTheDocument();
    expect(
      screen.queryByRole("link", { name: /daftar sekarang/i }),
    ).not.toBeInTheDocument();
  });

  it("points the hero dashboard CTA to the customer dashboard when logged in", () => {
    mockedListProducts.mockResolvedValue(productListResponse(productFixtures));

    renderAppAtRoot({
      user: customerUser,
      token: "customer-token",
      isAuthenticated: true,
    });

    expect(
      screen.getByRole("link", { name: /buka dashboard/i }),
    ).toHaveAttribute("href", "/dashboard");

    expect(
      screen.getByRole("link", { name: /lanjut belanja/i }),
    ).toHaveAttribute("href", "/products");
  });

  it("points the hero dashboard CTA to the admin dashboard when logged in as admin", () => {
    mockedListProducts.mockResolvedValue(productListResponse(productFixtures));

    renderAppAtRoot({
      user: adminUser,
      token: "admin-token",
      isAuthenticated: true,
    });

    expect(
      screen.getByRole("link", { name: /buka dashboard/i }),
    ).toHaveAttribute("href", "/admin/dashboard");
  });

  it("loads featured products through the mocked product service", async () => {
    mockedListProducts.mockResolvedValue(productListResponse(productFixtures));

    renderAppAtRoot();

    expect(await screen.findByText("Brutal Keyboard")).toBeInTheDocument();
    expect(screen.getByText("Sharp Mouse")).toBeInTheDocument();

    expect(
      screen.getByRole("link", { name: /lihat produk brutal keyboard/i }),
    ).toHaveAttribute("href", "/products/brutal-keyboard");
    expect(screen.getByText("Grid Monitor")).toBeInTheDocument();
    expect(screen.getAllByText(/in stock/i).length).toBeGreaterThan(0);
    expect(screen.getAllByText(/low stock/i).length).toBeGreaterThan(0);
    expect(screen.getByText(/out of stock/i)).toBeInTheDocument();

    expect(mockedListProducts).toHaveBeenCalledWith({
      page: 1,
      limit: 8,
      sort_by: "created_at",
      sort_order: "desc",
    });
  });

  it("shows a safe featured product fallback when loading fails", async () => {
    mockedListProducts.mockRejectedValue(new Error("Network down"));

    renderAppAtRoot();

    expect(
      await screen.findByText("Produk belum tersedia"),
    ).toBeInTheDocument();
  });

  it("shows an empty featured products fallback when API returns no products", async () => {
    mockedListProducts.mockResolvedValue(productListResponse([]));

    renderAppAtRoot();

    expect(
      await screen.findByText("Produk belum tersedia"),
    ).toBeInTheDocument();
  });
});
