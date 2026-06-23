import axios from "axios";

import api from "./api";
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

export async function payOrder(input: PayOrderInput): Promise<PaymentResult> {
  const response = await api.post<ApiDataResponse<PaymentResult>>(
    "/v1/payments/pay",
    input,
  );

  return response.data.data;
}