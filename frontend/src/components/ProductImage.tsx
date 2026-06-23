import { useState } from "react";

import { getImageUrl } from "../utils/image";

type ProductImageStatus = "loading" | "loaded" | "error";

interface ProductImageProps {
  src?: string | null;
  alt: string;
  className?: string;
  fallbackSrc?: string | null;
  width?: number;
  height?: number;
  loading?: "eager" | "lazy";
}

function ProductImagePlaceholder({ alt }: { alt: string }) {
  return (
    <span className="product-image-placeholder" role="img" aria-label={`No image for ${alt}`}>
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.4">
        <path d="M4 7h16v12H4z" />
        <path d="M8 7a4 4 0 0 1 8 0" />
        <path d="M8 13h8" />
      </svg>
      <span>Image not available</span>
    </span>
  );
}

function ProductImage({
  src,
  alt,
  className = "",
  fallbackSrc = "",
  width,
  height,
  loading = "lazy",
}: ProductImageProps) {
  const initialSrc = getImageUrl(src) || getImageUrl(fallbackSrc);
  const [activeSrc, setActiveSrc] = useState(initialSrc);
  const [status, setStatus] = useState<ProductImageStatus>(
    initialSrc ? "loading" : "error",
  );

  function handleLoad() {
    setStatus("loaded");
  }

  function handleError() {
    const fallbackUrl = getImageUrl(fallbackSrc);

    if (fallbackUrl && fallbackUrl !== activeSrc) {
      setActiveSrc(fallbackUrl);
      setStatus("loading");
      return;
    }

    setStatus("error");
  }

  const rootClassName = [
    "product-image",
    className,
    status === "loaded" ? "is-loaded" : "",
    status === "error" ? "is-error" : "",
  ]
    .filter(Boolean)
    .join(" ");

  return (
    <span className={rootClassName}>
      {status === "loading" && <span className="product-image-skeleton" aria-hidden="true" />}

      {activeSrc && status !== "error" ? (
        <img
          className="product-image-element"
          src={activeSrc}
          alt={alt}
          width={width}
          height={height}
          loading={loading}
          decoding="async"
          onLoad={handleLoad}
          onError={handleError}
        />
      ) : (
        <ProductImagePlaceholder alt={alt} />
      )}
    </span>
  );
}

export default ProductImage;
