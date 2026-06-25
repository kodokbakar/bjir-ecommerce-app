import { useState, type CSSProperties } from "react";
import { ImageOff } from "lucide-react";

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
  sizes?: string;
  srcSet?: string;
}

function ProductImagePlaceholder({ alt }: { alt: string }) {
  return (
    <span
      className="product-image-placeholder"
      role="img"
      aria-label={`No image for ${alt}`}
    >
      <ImageOff className="h-11 w-11" aria-hidden="true" />
      <span>Image not available</span>
    </span>
  );
}

function ProductImage({
  src,
  alt,
  className = "",
  fallbackSrc = "",
  width = 640,
  height = 480,
  loading = "lazy",
  sizes = "(max-width: 720px) 100vw, (max-width: 1180px) 50vw, 25vw",
  srcSet,
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

  const style = {
    aspectRatio: `${width} / ${height}`,
  } satisfies CSSProperties;

  return (
    <span className={rootClassName} style={style}>
      {status === "loading" && (
        <span className="product-image-skeleton" aria-hidden="true" />
      )}

      {activeSrc && status !== "error" ? (
        <img
          className="product-image-element"
          src={activeSrc}
          srcSet={srcSet}
          sizes={sizes}
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
