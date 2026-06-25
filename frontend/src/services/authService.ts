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

interface MeIdentityResponse {
  user_id: string;
  email: string;
  role: string;
}

type MeResponse =
  | ApiDataResponse<User | { user?: User } | MeIdentityResponse>
  | {
      user?: User;
    };

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

type ApiErrorContext = "auth" | "profile";

export async function loginUser(input: LoginInput): Promise<AuthResult> {
  const response = await api.post<ApiDataResponse<AuthPayload>>(
    "/v1/auth/login",
    input,
  );
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

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null;
}

function isUser(value: unknown): value is User {
  return (
    isRecord(value) &&
    typeof value.id === "string" &&
    typeof value.email === "string"
  );
}

function isWrappedUser(value: unknown): value is { user: User } {
  return isRecord(value) && isUser(value.user);
}

function isMeIdentityResponse(value: unknown): value is MeIdentityResponse {
  return (
    isRecord(value) &&
    typeof value.user_id === "string" &&
    typeof value.email === "string" &&
    typeof value.role === "string"
  );
}

function normalizeUser(user: User, fallbackUser: User | null = null): User {
  return {
    ...fallbackUser,
    ...user,
    name: user.name || fallbackUser?.name || user.email,
    is_active: user.is_active ?? fallbackUser?.is_active ?? true,
  };
}

function normalizeMeIdentity(
  identity: MeIdentityResponse,
  fallbackUser: User | null = null,
): User {
  return {
    id: identity.user_id || fallbackUser?.id || "",
    name: fallbackUser?.name || identity.email,
    email: identity.email || fallbackUser?.email || "",
    role: identity.role || fallbackUser?.role,
    is_active: fallbackUser?.is_active ?? true,
    created_at: fallbackUser?.created_at,
  };
}

export function unwrapCurrentUser(
  payload: MeResponse,
  fallbackUser: User | null = null,
): User {
  if ("data" in payload) {
    const data = payload.data;

    if (isWrappedUser(data)) {
      return normalizeUser(data.user, fallbackUser);
    }

    if (isMeIdentityResponse(data)) {
      return normalizeMeIdentity(data, fallbackUser);
    }

    if (isUser(data)) {
      return normalizeUser(data, fallbackUser);
    }
  }

  if ("user" in payload && isUser(payload.user)) {
    return normalizeUser(payload.user, fallbackUser);
  }

  throw new Error("Format respon profil dari server tidak sesuai.");
}

export async function getCurrentUser(
  fallbackUser: User | null = null,
): Promise<User> {
  const response = await api.get<MeResponse>("/v1/me");

  return unwrapCurrentUser(response.data, fallbackUser);
}

export function getApiErrorMessage(
  error: unknown,
  fallback: string,
  context: ApiErrorContext = "auth",
): string {
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
      return context === "profile"
        ? "Sesi kamu sudah berakhir, silakan login ulang."
        : "Email atau kata sandi tidak sesuai.";
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
