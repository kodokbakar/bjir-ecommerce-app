import api, { getResponseErrorMessage } from "./api";
import type { PayOrderInput, PaymentResult } from "../types/payment";

interface ApiDataResponse<TData> {
  success?: boolean;
  message?: string;
  error?: string;
  details?: string;
  data: TData;
}

export function getPaymentErrorMessage(
  error: unknown,
  fallback: string,
): string {
  return getResponseErrorMessage(error, fallback);
}

export async function payOrder(input: PayOrderInput): Promise<PaymentResult> {
  const response = await api.post<ApiDataResponse<PaymentResult>>(
    "/v1/payments/pay",
    input,
  );

  return response.data.data;
}