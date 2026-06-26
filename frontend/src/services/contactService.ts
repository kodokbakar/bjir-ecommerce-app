import api, { getResponseErrorMessage } from "./api";

export interface ContactInput {
  name: string;
  email: string;
  message: string;
}

export async function sendContactMessage(input: ContactInput): Promise<void> {
  await api.post("/v1/contact", input);
}

export function getContactErrorMessage(
  error: unknown,
  fallback: string,
): string {
  return getResponseErrorMessage(error, fallback);
}
