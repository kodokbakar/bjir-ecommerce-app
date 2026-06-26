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
  {
    id: "prod-monitor",
    name: "Grid Monitor",
    slug: "grid-monitor",
    description: "A bright monitor for commerce dashboards.",
    price: 1850000,
    stock: 0,
    category_id: "cat-keyboards",
    category: categoryFixtures[0],
    image_url: "/uploads/products/monitor.webp",
    is_active: true,
    created_at: "2026-06-25T03:00:00Z",
    updated_at: "2026-06-25T03:00:00Z",
  },
  {
    id: "prod-headset",
    name: "Orange Headset",
    slug: "orange-headset",
    description: "A loud headset for support calls.",
    price: 480000,
    stock: 12,
    category_id: "cat-keyboards",
    category: categoryFixtures[0],
    image_url: "/uploads/products/headset.webp",
    is_active: true,
    created_at: "2026-06-25T04:00:00Z",
    updated_at: "2026-06-25T04:00:00Z",
  },
  {
    id: "prod-deskmat",
    name: "Ink Deskmat",
    slug: "ink-deskmat",
    description: "Deskmat with brutal grid lines.",
    price: 210000,
    stock: 4,
    category_id: "cat-keyboards",
    category: categoryFixtures[0],
    image_url: "/uploads/products/deskmat.webp",
    is_active: true,
    created_at: "2026-06-25T05:00:00Z",
    updated_at: "2026-06-25T05:00:00Z",
  },
  {
    id: "prod-speaker",
    name: "Shelf Speaker",
    slug: "shelf-speaker",
    description: "Compact speaker for storefront noise.",
    price: 390000,
    stock: 7,
    category_id: "cat-keyboards",
    category: categoryFixtures[0],
    image_url: "/uploads/products/speaker.webp",
    is_active: true,
    created_at: "2026-06-25T06:00:00Z",
    updated_at: "2026-06-25T06:00:00Z",
  },
  {
    id: "prod-camera",
    name: "Catalog Camera",
    slug: "catalog-camera",
    description: "Camera for product photos.",
    price: 2450000,
    stock: 1,
    category_id: "cat-keyboards",
    category: categoryFixtures[0],
    image_url: "/uploads/products/camera.webp",
    is_active: true,
    created_at: "2026-06-25T07:00:00Z",
    updated_at: "2026-06-25T07:00:00Z",
  },
  {
    id: "prod-stand",
    name: "Counter Stand",
    slug: "counter-stand",
    description: "A hard-edged laptop stand.",
    price: 260000,
    stock: 15,
    category_id: "cat-keyboards",
    category: categoryFixtures[0],
    image_url: "/uploads/products/stand.webp",
    is_active: true,
    created_at: "2026-06-25T08:00:00Z",
    updated_at: "2026-06-25T08:00:00Z",
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
