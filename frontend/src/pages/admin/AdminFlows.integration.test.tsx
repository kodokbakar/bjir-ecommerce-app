import { screen, waitFor, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import AdminDashboard from "./AdminDashboard";
import AdminOrders from "./AdminOrders";
import AdminProducts from "./AdminProducts";
import ProductForm from "./ProductForm";
import {
  createProduct,
  deleteProduct,
  getProductById,
  listCategories,
  listProducts,
  updateProduct,
  uploadProductImage,
} from "../../services/productService";
import {
  listAdminOrders,
  updateOrderStatus,
} from "../../services/orderService";
import { getAdminDashboardStats } from "../../services/dashboardService";
import {
  adminStatsFixture,
  categoryFixtures,
  orderFixture,
  orderListResponse,
  productFixtures,
  productListResponse,
} from "../../test-utils/fixtures";
import {
  adminUser,
  renderWithProviders,
} from "../../test-utils/renderWithProviders";

vi.mock("../../services/productService", async () => {
  const actual = await vi.importActual<
    typeof import("../../services/productService")
  >("../../services/productService");

  return {
    ...actual,
    listCategories: vi.fn(),
    listProducts: vi.fn(),
    getProductById: vi.fn(),
    createProduct: vi.fn(),
    updateProduct: vi.fn(),
    uploadProductImage: vi.fn(),
    deleteProduct: vi.fn(),
  };
});

vi.mock("../../services/orderService", async () => {
  const actual = await vi.importActual<
    typeof import("../../services/orderService")
  >("../../services/orderService");

  return {
    ...actual,
    listAdminOrders: vi.fn(),
    updateOrderStatus: vi.fn(),
  };
});

vi.mock("../../services/dashboardService", async () => {
  const actual = await vi.importActual<
    typeof import("../../services/dashboardService")
  >("../../services/dashboardService");

  return {
    ...actual,
    getAdminDashboardStats: vi.fn(),
  };
});

const mockedListCategories = vi.mocked(listCategories);
const mockedListProducts = vi.mocked(listProducts);
const mockedGetProductById = vi.mocked(getProductById);
const mockedCreateProduct = vi.mocked(createProduct);
const mockedUpdateProduct = vi.mocked(updateProduct);
const mockedUploadProductImage = vi.mocked(uploadProductImage);
const mockedDeleteProduct = vi.mocked(deleteProduct);
const mockedListAdminOrders = vi.mocked(listAdminOrders);
const mockedUpdateOrderStatus = vi.mocked(updateOrderStatus);
const mockedGetAdminDashboardStats = vi.mocked(getAdminDashboardStats);

function renderAdmin(
  routes: Parameters<typeof renderWithProviders>[0],
  path: string,
) {
  return renderWithProviders(routes, {
    initialEntries: [path],
    auth: {
      user: adminUser,
      token: "admin-token",
      isAuthenticated: true,
      isLoading: false,
    },
  });
}

describe("admin critical flows", () => {
  beforeEach(() => {
    mockedListCategories.mockReset();
    mockedListProducts.mockReset();
    mockedGetProductById.mockReset();
    mockedCreateProduct.mockReset();
    mockedUpdateProduct.mockReset();
    mockedUploadProductImage.mockReset();
    mockedDeleteProduct.mockReset();
    mockedListAdminOrders.mockReset();
    mockedUpdateOrderStatus.mockReset();
    mockedGetAdminDashboardStats.mockReset();
  });

  it("shows admin dashboard stat cards", async () => {
    mockedGetAdminDashboardStats.mockResolvedValue(adminStatsFixture);

    renderAdmin(
      [
        {
          path: "/admin/dashboard",
          element: <AdminDashboard />,
        },
      ],
      "/admin/dashboard",
    );

    expect(
      await screen.findByRole("heading", { name: /control room/i }),
    ).toBeInTheDocument();
    expect(screen.getByText("Total Orders")).toBeInTheDocument();
    expect(screen.getByText("12")).toBeInTheDocument();
    expect(screen.getByText("Products")).toBeInTheDocument();
    expect(screen.getByText("9")).toBeInTheDocument();
  });

  it("shows admin products and deletes a product after confirmation", async () => {
    const user = userEvent.setup();

    mockedListProducts.mockResolvedValue(
      productListResponse([productFixtures[0]]),
    );
    mockedDeleteProduct.mockResolvedValue(undefined);
    vi.spyOn(window, "confirm").mockReturnValue(true);

    renderAdmin(
      [
        {
          path: "/admin/products",
          element: <AdminProducts />,
        },
      ],
      "/admin/products",
    );

    expect(await screen.findByText("Brutal Keyboard")).toBeInTheDocument();

    await user.click(screen.getByRole("button", { name: /delete/i }));

    expect(mockedDeleteProduct).toHaveBeenCalledWith("prod-keyboard");
    expect(
      await screen.findByText("Brutal Keyboard deleted."),
    ).toBeInTheDocument();
    expect(await screen.findByText("No products found.")).toBeInTheDocument();
  });

  it("creates a product from admin form", async () => {
    const user = userEvent.setup();

    mockedListCategories.mockResolvedValue(categoryFixtures);
    mockedCreateProduct.mockResolvedValue(productFixtures[0]);
    mockedUploadProductImage.mockResolvedValue(undefined);

    renderAdmin(
      [
        {
          path: "/admin/products/new",
          element: <ProductForm />,
        },
        {
          path: "/admin/products",
          element: <p>Admin products route</p>,
        },
      ],
      "/admin/products/new",
    );

    expect(
      await screen.findByRole("heading", { name: /new item/i }),
    ).toBeInTheDocument();

    await user.type(screen.getByLabelText(/^nama$/i), "Brutal Keyboard");
    await user.selectOptions(
      screen.getByLabelText(/kategori/i),
      "cat-keyboards",
    );
    await user.type(screen.getByLabelText(/harga/i), "750000");
    await user.clear(screen.getByLabelText(/stok/i));
    await user.type(screen.getByLabelText(/stok/i), "8");
    await user.type(screen.getByLabelText(/deskripsi/i), "A loud keyboard.");
    await user.click(screen.getByRole("button", { name: /save product/i }));

    await waitFor(() => {
      expect(mockedCreateProduct).toHaveBeenCalledWith(
        expect.objectContaining({
          category_id: "cat-keyboards",
          name: "Brutal Keyboard",
          price: 750000,
          stock: 8,
          description: "A loud keyboard.",
        }),
      );
    });

    expect(await screen.findByText("Admin products route")).toBeInTheDocument();
  });

  it("loads edit product form with existing data", async () => {
    mockedListCategories.mockResolvedValue(categoryFixtures);
    mockedGetProductById.mockResolvedValue(productFixtures[0]);
    mockedUpdateProduct.mockResolvedValue(productFixtures[0]);

    renderAdmin(
      [
        {
          path: "/admin/products/edit",
          element: <ProductForm />,
        },
      ],
      "/admin/products/edit?id=prod-keyboard",
    );

    expect(
      await screen.findByRole("heading", { name: /edit item/i }),
    ).toBeInTheDocument();
    expect(screen.getByDisplayValue("Brutal Keyboard")).toBeInTheDocument();
    expect(screen.getByDisplayValue("8")).toBeInTheDocument();
  });

  it("filters admin orders by status and updates order status", async () => {
    const user = userEvent.setup();

    const paidOrder = {
      ...orderFixture,
      status: "paid",
    } satisfies typeof orderFixture;

    mockedListAdminOrders.mockResolvedValue(orderListResponse([orderFixture]));
    mockedUpdateOrderStatus.mockResolvedValue(paidOrder);

    renderAdmin(
      [
        {
          path: "/admin/orders",
          element: <AdminOrders />,
        },
      ],
      "/admin/orders?status=pending",
    );

    expect(await screen.findByText("ORD-0001")).toBeInTheDocument();

    expect(mockedListAdminOrders).toHaveBeenCalledWith(
      expect.objectContaining({
        status: "pending",
      }),
    );

    const row = screen.getByText("ORD-0001").closest("article");

    expect(row).not.toBeNull();

    await user.click(
      within(row as HTMLElement).getByRole("button", { name: /status/i }),
    );

    expect(await screen.findByRole("dialog")).toBeInTheDocument();

    await user.click(screen.getByRole("button", { name: /confirm/i }));

    expect(mockedUpdateOrderStatus).toHaveBeenCalledWith("order-1", "paid");
    expect(
      await screen.findByText("ORD-0001 changed from Pending to Paid."),
    ).toBeInTheDocument();
  });
});
