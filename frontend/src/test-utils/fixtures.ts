import type { Cart } from "../types/cart";
import type { Order, OrderListResponse } from "../types/order";
import type { Category, Product, ProductListResponse } from "../types/product";
import type { AdminDashboardStats } from "../services/dashboardService";

export const categoryFixtures: Category[] = [
  {
    id: "cat-keyboards",
    name: "Keyboards",
    slug: "keyboards",
    description: "Mechanical boards.",
  },
];

export const productFixtures: Product[] = [
  {
    id: "prod-keyboard",
    name: "Brutal Keyboard",
    slug: "brutal-keyboard",
    description: "A loud mechanical keyboard for serious typing.",
    price: 750000,
    stock: 8,
    category_id: "cat-keyboards",
    category: categoryFixtures[0],
    image_url: "/uploads/products/keyboard.webp",
    is_active: true,
    created_at: "2026-06-25T01:00:00Z",
    updated_at: "2026-06-25T01:00:00Z",
  },
  {
    id: "prod-mouse",
    name: "Sharp Mouse",
    slug: "sharp-mouse",
    description: "Fast mouse.",
    price: 320000,
    stock: 2,
    category_id: "cat-keyboards",
    category: categoryFixtures[0],
    image_url: "/uploads/products/mouse.webp",
    is_active: true,
    created_at: "2026-06-25T02:00:00Z",
    updated_at: "2026-06-25T02:00:00Z",
  },
];

export function productListResponse(
  products = productFixtures,
): ProductListResponse {
  return {
    data: products,
    meta: {
      page: 1,
      limit: 12,
      total: products.length,
      total_pages: 1,
      sort_by: "",
      sort_order: "",
      category_id: "",
      category: "",
      search: "",
    },
  };
}

export const cartFixture: Cart = {
  items: [
    {
      id: "cart-item-keyboard",
      product_id: "prod-keyboard",
      quantity: 1,
      subtotal: 750000,
      product: productFixtures[0],
    },
  ],
  total_price: 750000,
} as Cart;

export const updatedCartItemFixture = {
  ...cartFixture.items[0],
  quantity: 2,
  subtotal: 1500000,
};

export const orderFixture: Order = {
  id: "order-1",
  order_number: "ORD-0001",
  user_id: "user-customer",
  user_name: "Bintang Customer",
  user_email: "customer@example.test",
  total_amount: 750000,
  status: "pending",
  shipping_address: "Jl. Testing No. 1",
  notes: "",
  items: [
    {
      id: "order-item-1",
      product_id: "prod-keyboard",
      product_name: "Brutal Keyboard",
      quantity: 1,
      price: 750000,
      subtotal: 750000,
    },
  ],
  created_at: "2026-06-25T03:00:00Z",
  updated_at: "2026-06-25T03:00:00Z",
} as Order;

export function orderListResponse(orders = [orderFixture]): OrderListResponse {
  return {
    data: orders,
    meta: {
      page: 1,
      limit: 10,
      total: orders.length,
      total_pages: 1,
    },
  };
}

export const adminStatsFixture: AdminDashboardStats = {
  totalOrders: 12,
  totalRevenue: 5000000,
  pendingOrders: 3,
  completedToday: 2,
  revenueToday: 1250000,
  totalProducts: 9,
  totalCategories: 4,
};
