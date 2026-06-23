import { getApiOrigin } from "../services/api";

export function getImageUrl(imageUrl?: string | null, fallbackSrc = ""): string {
  const trimmed = imageUrl?.trim() ?? "";

  if (!trimmed) {
    return fallbackSrc.trim();
  }

  if (/^https?:\/\//i.test(trimmed)) {
    return trimmed;
  }

  const imagePath = trimmed.replace(/^\/+/, "");
  const uploadPath = imagePath.startsWith("uploads/")
    ? `/${imagePath}`
    : `/uploads/${imagePath}`;

  return `${getApiOrigin()}${uploadPath}`;
}
