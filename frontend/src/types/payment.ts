export type PaymentMethod = "bank_transfer" | "credit_card" | "ewallet";

export interface PayOrderInput {
  order_id: string;
  method: PaymentMethod;
}

export interface PaymentResult {
  id: string;
  order_id: string;
  provider: string;
  payment_method: PaymentMethod;
  transaction_id: string;
  amount: number;
  status: string;
  paid_at?: string;
}