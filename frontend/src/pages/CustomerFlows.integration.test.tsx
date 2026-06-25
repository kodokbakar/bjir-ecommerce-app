import { screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import Cart from "./Cart";
import Checkout from "./Checkout";
import ProductDetail from "./ProductDetail";
import Products from "./Products";
import {
  addCartItem,
  checkoutCart,
  getCart,
  removeCartItem,
  updateCartItem,
} from "../services/cartService";
import {
  getProductBySlug,
  listCategories,
  listProducts,
} from "../services/productService";
import {
  cartFixture,
  categoryFixtures,
  orderFixture,
  productFixtures,
  productListResponse,
  updatedCartItemFixture,
} from "../test-utils/fixtures";
import { renderWithProviders } from "../test-utils/renderWithProviders";

vi.mock("../services/productService", async () => {
  const actual = await vi.importActual<
    typeof import("../services/productService")
  >("../services/productService");

  return {
    ...actual,
    listCategories: vi.fn(),
    listProducts: vi.fn(),
    getProductBySlug: vi.fn(),
  };
});

vi.mock("../services/cartService", async () => {
  const actual = await vi.importActual<
    typeof import("../services/cartService")
  >("../services/cartService");

  return {
    ...actual,
    getCart: vi.fn(),
    addCartItem: vi.fn(),
    updateCartItem: vi.fn(),
    removeCartItem: vi.fn(),
    checkoutCart: vi.fn(),
  };
});

const mockedListCategories = vi.mocked(listCategories);
const mockedListProducts = vi.mocked(listProducts);
const mockedGetProductBySlug = vi.mocked(getProductBySlug);
const mockedGetCart = vi.mocked(getCart);
const mockedAddCartItem = vi.mocked(addCartItem);
const mockedUpdateCartItem = vi.mocked(updateCartItem);
const mockedRemoveCartItem = vi.mocked(removeCartItem);
const mockedCheckoutCart = vi.mocked(checkoutCart);

describe("customer browsing, cart, and checkout flow", () => {
  beforeEach(() => {
    mockedListCategories.mockReset();
    mockedListProducts.mockReset();
    mockedGetProductBySlug.mockReset();
    mockedGetCart.mockReset();
    mockedAddCartItem.mockReset();
    mockedUpdateCartItem.mockReset();
    mockedRemoveCartItem.mockReset();
    mockedCheckoutCart.mockReset();
  });

  it("shows product list with pagination data", async () => {
    mockedListCategories.mockResolvedValue(categoryFixtures);
    mockedListProducts.mockResolvedValue(productListResponse(productFixtures));

    renderWithProviders(
      [
        {
          path: "/products",
          element: <Products />,
        },
      ],
      {
        initialEntries: ["/products"],
      },
    );

    expect(await screen.findByText("Brutal Keyboard")).toBeInTheDocument();
    expect(screen.getByText("Sharp Mouse")).toBeInTheDocument();
    expect(screen.getByText(/2 products found/i)).toBeInTheDocument();
  });

  it("loads product list using search query params", async () => {
    mockedListCategories.mockResolvedValue(categoryFixtures);
    mockedListProducts.mockResolvedValue(
      productListResponse([productFixtures[0]]),
    );

    renderWithProviders(
      [
        {
          path: "/products",
          element: <Products />,
        },
      ],
      {
        initialEntries: ["/products?search=keyboard"],
      },
    );

    expect(await screen.findByText("Brutal Keyboard")).toBeInTheDocument();
    expect(mockedListProducts).toHaveBeenCalledWith(
      expect.objectContaining({
        search: "keyboard",
      }),
    );
  });

  it("shows product detail and adds product to cart", async () => {
    const user = userEvent.setup();

    mockedGetProductBySlug.mockResolvedValue(productFixtures[0]);
    mockedAddCartItem.mockResolvedValue(cartFixture.items[0]);

    renderWithProviders(
      [
        {
          path: "/products/:slug",
          element: <ProductDetail />,
        },
      ],
      {
        initialEntries: ["/products/brutal-keyboard"],
      },
    );

    expect(
      await screen.findByRole("heading", { name: "Brutal Keyboard" }),
    ).toBeInTheDocument();

    await user.click(screen.getByRole("button", { name: /add to cart/i }));

    expect(mockedAddCartItem).toHaveBeenCalledWith("prod-keyboard", 1);
    expect(
      await screen.findByText("Product added to cart."),
    ).toBeInTheDocument();
  });

  it("updates cart quantity, recalculates total, and removes the item", async () => {
    const user = userEvent.setup();

    mockedGetCart.mockResolvedValue(cartFixture);
    mockedUpdateCartItem.mockResolvedValue(updatedCartItemFixture);
    mockedRemoveCartItem.mockResolvedValue(undefined);

    vi.spyOn(window, "confirm").mockReturnValue(true);

    renderWithProviders(
      [
        {
          path: "/cart",
          element: <Cart />,
        },
      ],
      {
        initialEntries: ["/cart"],
      },
    );

    expect(await screen.findByText("Brutal Keyboard")).toBeInTheDocument();
    expect(screen.getAllByText(/Rp\s*750\.000/).length).toBeGreaterThan(0);

    await user.click(
      screen.getByRole("button", {
        name: /increase brutal keyboard quantity/i,
      }),
    );

    expect(mockedUpdateCartItem).toHaveBeenCalledWith("cart-item-keyboard", 2);
    await waitFor(() => {
      expect(screen.getAllByText(/Rp\s*1\.500\.000/).length).toBeGreaterThan(0);
    });

    await user.click(screen.getByRole("button", { name: /remove/i }));

    expect(mockedRemoveCartItem).toHaveBeenCalledWith("cart-item-keyboard");
    expect(await screen.findByText("Cart is empty.")).toBeInTheDocument();
  });

  it("shows empty cart state", async () => {
    mockedGetCart.mockResolvedValue({
      items: [],
      total_price: 0,
    });

    renderWithProviders(
      [
        {
          path: "/cart",
          element: <Cart />,
        },
      ],
      {
        initialEntries: ["/cart"],
      },
    );

    expect(await screen.findByText("Cart is empty.")).toBeInTheDocument();
  });

  it("shows checkout shipping form and validates missing shipping address", async () => {
    const user = userEvent.setup();

    mockedGetCart.mockResolvedValue(cartFixture);

    renderWithProviders(
      [
        {
          path: "/checkout",
          element: <Checkout />,
        },
      ],
      {
        initialEntries: ["/checkout"],
      },
    );

    expect(
      await screen.findByText("Where should this go?"),
    ).toBeInTheDocument();

    await user.click(screen.getByRole("button", { name: /place order/i }));

    expect(
      await screen.findByText("Shipping address wajib diisi."),
    ).toBeInTheDocument();
    expect(mockedCheckoutCart).not.toHaveBeenCalled();
  });

  it("places order with valid shipping address and redirects to dashboard", async () => {
    const user = userEvent.setup();

    mockedGetCart.mockResolvedValue(cartFixture);
    mockedCheckoutCart.mockResolvedValue(orderFixture);

    renderWithProviders(
      [
        {
          path: "/checkout",
          element: <Checkout />,
        },
        {
          path: "/dashboard",
          element: <p>Dashboard after checkout</p>,
        },
      ],
      {
        initialEntries: ["/checkout"],
      },
    );

    expect(
      await screen.findByText("Where should this go?"),
    ).toBeInTheDocument();

    await user.type(
      screen.getByLabelText(/shipping address/i),
      "Jl. Testing No. 1",
    );
    await user.click(screen.getByRole("button", { name: /place order/i }));

    await waitFor(() => {
      expect(mockedCheckoutCart).toHaveBeenCalledWith({
        shipping_address: "Jl. Testing No. 1",
        notes: "",
      });
    });

    expect(
      await screen.findByText("Dashboard after checkout"),
    ).toBeInTheDocument();
  });

  it("shows empty cart checkout state", async () => {
    mockedGetCart.mockResolvedValue({
      items: [],
      total_price: 0,
    });

    renderWithProviders(
      [
        {
          path: "/checkout",
          element: <Checkout />,
        },
      ],
      {
        initialEntries: ["/checkout"],
      },
    );

    expect(await screen.findByText("Keranjang kosong.")).toBeInTheDocument();
  });
});
