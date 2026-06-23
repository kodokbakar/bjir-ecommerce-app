import axios from "axios";

import api from "./api";
import type { User } from "../hooks/useAuth";

interface ApiDataResponse<T> {
  success?: boolean;
  message?: string;
  data: T;
}

interface AuthPayload {
  access_token: string;
  token_type?: string;
  expires_in?: number;
  user: User;
}

export interface LoginInput {
  email: string;
  password: string;
}

export interface RegisterInput {
  name: string;
  email: string;
  password: string;
}

export interface AuthResult {
  accessToken: string;
  user: User;
}

export async function loginUser(input: LoginInput): Promise<AuthResult> {
  const response = await api.post<ApiDataResponse<AuthPayload>>("/v1/auth/login", input);
  const payload = response.data.data;

  if (!payload?.access_token || !payload.user) {
    throw new Error("Format respon login dari server tidak sesuai.");
  }

  return {
    accessToken: payload.access_token,
    user: payload.user,
  };
}

export async function registerUser(input: RegisterInput): Promise<void> {
  await api.post<ApiDataResponse<AuthPayload>>("/v1/auth/register", input);
}

export function getApiErrorMessage(error: unknown, fallback: string): string {
  if (axios.isAxiosError(error)) {
    const responseData = error.response?.data;

    if (
      responseData &&
      typeof responseData === "object" &&
      "message" in responseData &&
      typeof responseData.message === "string"
    ) {
      return responseData.message;
    }

    if (error.response?.status === 401) {
      return "Email atau kata sandi tidak sesuai.";
    }

    if (error.response?.status === 409) {
      return "Email sudah terdaftar.";
    }

    if (error.response?.status === 403) {
      return "Akun ini sedang tidak aktif.";
    }
  }

  if (error instanceof Error && error.message) {
    return error.message;
  }

  return fallback;
}
