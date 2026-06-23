import axios from "axios";

import api from "./api";
import type {
  AddCartItemInput,
  Cart,
  CartItem,
  UpdateCartItemInput,
} from "../types/cart";
import type { CheckoutInput, Order } from "../types/order";

type ApiDataResponse<T> = {
  success?: boolean;
  message?: string;
  error?: string;
  details?: string;
  data: T;
};

type CartChangeListener = () => void;

const cartChangeListeners = new Set<CartChangeListener>();

function isApiDataResponse<T>(value: unknown): value is ApiDataResponse<T> {
  return typeof value === "object" && value !== null && "data" in value;
}

function unwrapData<T>(value: T | ApiDataResponse<T>): T {
  return isApiDataResponse<T>(value) ? value.data : value;
}

function requireValue(value: string, fieldName: string): string {
  const trimmed = value.trim();

  if (!trimmed) {
    throw new Error(`${fieldName} is required`);
  }

  return trimmed;
}

function requirePositiveQuantity(quantity: number): number {
  if (!Number.isInteger(quantity) || quantity < 1) {
    throw new Error("Quantity must be at least 1");
  }

  return quantity;
}

export function subscribeToCartChanges(
  listener: CartChangeListener,
): () => void {
  cartChangeListeners.add(listener);

  return () => {
    cartChangeListeners.delete(listener);
  };
}

export function notifyCartChanged(): void {
  cartChangeListeners.forEach((listener) => {
    listener();
  });
}

export function getCartItemPrice(item: CartItem): number {
  if (typeof item.product?.price === "number") {
    return item.product.price;
  }

  return item.quantity > 0 ? item.subtotal / item.quantity : 0;
}

export function getCartItemSubtotal(item: CartItem): number {
  return getCartItemPrice(item) * item.quantity;
}

export function getCartItemCount(cart: Cart): number {
  return cart.items.length;
}

export function normalizeCart(cart: Cart): Cart {
  const items = cart.items.map((item) => ({
    ...item,
    subtotal: getCartItemSubtotal(item),
  }));

  return {
    items,
    total_price: items.reduce(
      (total, item) => total + getCartItemSubtotal(item),
      0,
    ),
  };
}

export function getCartErrorMessage(error: unknown, fallback: string): string {
  if (axios.isAxiosError(error)) {
    const responseData = error.response?.data as
      | {
          message?: string;
          error?: string;
          details?: string;
        }
      | undefined;

    return (
      responseData?.details ||
      responseData?.message ||
      responseData?.error ||
      fallback
    );
  }

  if (error instanceof Error && error.message) {
    return error.message;
  }

  return fallback;
}

export async function getCart(): Promise<Cart> {
  const response = await api.get<Cart | ApiDataResponse<Cart>>("/v1/cart");

  return unwrapData(response.data);
}

export async function getCartCount(): Promise<number> {
  const cart = await getCart();

  return getCartItemCount(normalizeCart(cart));
}

export async function addCartItem(
  productID: string,
  quantity = 1,
): Promise<CartItem> {
  const payload: AddCartItemInput = {
    product_id: requireValue(productID, "Product ID"),
    quantity: requirePositiveQuantity(quantity),
  };

  const response = await api.post<CartItem | ApiDataResponse<CartItem>>(
    "/v1/cart/items",
    payload,
  );

  const item = unwrapData(response.data);
  notifyCartChanged();

  return item;
}

export async function updateCartItem(
  itemID: string,
  quantity: number,
): Promise<CartItem> {
  const payload: UpdateCartItemInput = {
    quantity: requirePositiveQuantity(quantity),
  };

  const response = await api.put<CartItem | ApiDataResponse<CartItem>>(
    `/v1/cart/items/${encodeURIComponent(requireValue(itemID, "Cart item ID"))}`,
    payload,
  );

  const item = unwrapData(response.data);
  notifyCartChanged();

  return item;
}

export async function removeCartItem(itemID: string): Promise<void> {
  await api.delete(
    `/v1/cart/items/${encodeURIComponent(requireValue(itemID, "Cart item ID"))}`,
  );

  notifyCartChanged();
}

export async function checkoutCart(input: CheckoutInput = {}): Promise<Order> {
  const shippingAddress = input.shipping_address?.trim();
  const notes = input.notes?.trim();

  const payload: CheckoutInput = {
    ...(shippingAddress ? { shipping_address: shippingAddress } : {}),
    ...(notes ? { notes } : {}),
  };

  const response = await api.post<Order | ApiDataResponse<Order>>(
    "/v1/orders/checkout",
    payload,
  );

  const order = unwrapData(response.data);
  notifyCartChanged();

  return order;
}
export const cartService = {
  getCart,
  getCartCount,
  addCartItem,
  updateCartItem,
  removeCartItem,
  checkoutCart,
  getCartItemPrice,
  getCartItemSubtotal,
  getCartItemCount,
  normalizeCart,
  subscribeToCartChanges,
  notifyCartChanged,
  getCartErrorMessage,
};
