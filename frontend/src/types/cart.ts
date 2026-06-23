import type { Product } from "./product";

export interface CartItem {
  id: string;
  user_id: string;
  product_id: string;
  product?: Product | null;
  quantity: number;
  subtotal: number;
  created_at?: string;
  updated_at?: string;
}

export interface Cart {
  items: CartItem[];
  total_price: number;
}

export interface AddCartItemInput {
  product_id: string;
  quantity: number;
}

export interface UpdateCartItemInput {
  quantity: number;
}