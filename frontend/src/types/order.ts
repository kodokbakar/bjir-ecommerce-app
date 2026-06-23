export type OrderStatus = "pending" | "paid" | "shipped" | "delivered" | "cancelled";

export interface CheckoutInput {
  shipping_address?: string;
  notes?: string;
}

export interface OrderItem {
  id: string;
  order_id: string;
  product_id: string;
  product_name: string;
  quantity: number;
  price: number;
  subtotal: number;
  created_at?: string;
}

export interface Order {
  id: string;
  user_id: string;
  order_number: string;
  status: OrderStatus;
  total_amount: number;
  shipping_address?: string;
  notes?: string;
  items?: OrderItem[];
  created_at?: string;
  updated_at?: string;
}