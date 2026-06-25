import { useEffect, useMemo, useState, type FormEvent } from "react";
import { Link, useNavigate, useSearchParams } from "react-router-dom";
import {
  AlertTriangle,
  ArrowLeft,
  ImagePlus,
  Loader2,
  Save,
} from "lucide-react";

import ProductImage from "../../components/ProductImage";
import {
  createProduct,
  getProductById,
  getProductErrorMessage,
  listCategories,
  updateProduct,
  uploadProductImage,
} from "../../services/productService";
import type { Category, Product, ProductInput } from "../../types/product";
import { formatRupiah, getProductImage } from "../../utils/product";

interface ProductFormState {
  name: string;
  categoryID: string;
  price: string;
  stock: string;
  description: string;
}

interface ProductFormErrors {
  name?: string;
  categoryID?: string;
  price?: string;
  stock?: string;
  form?: string;
  image?: string;
}

const EMPTY_FORM: ProductFormState = {
  name: "",
  categoryID: "",
  price: "",
  stock: "0",
  description: "",
};

function parseNumberInput(value: string): number {
  const digits = value.replace(/[^\d]/g, "");

  if (!digits) {
    return 0;
  }

  return Number(digits);
}

function formatPriceInput(value: string): string {
  const parsed = parseNumberInput(value);

  if (parsed <= 0) {
    return "";
  }

  return formatRupiah(parsed);
}

function buildFormFromProduct(product: Product): ProductFormState {
  return {
    name: product.name,
    categoryID: product.category_id ?? product.category?.id ?? "",
    price: product.price > 0 ? formatRupiah(product.price) : "",
    stock: String(product.stock ?? 0),
    description: product.description ?? "",
  };
}

function validateForm(form: ProductFormState): ProductFormErrors {
  const errors: ProductFormErrors = {};
  const price = parseNumberInput(form.price);
  const stock = Number(form.stock);

  if (!form.name.trim()) {
    errors.name = "Nama produk wajib diisi.";
  } else if (form.name.trim().length > 150) {
    errors.name = "Nama produk maksimal 150 karakter.";
  }

  if (!form.categoryID.trim()) {
    errors.categoryID = "Kategori wajib dipilih.";
  }

  if (price <= 0) {
    errors.price = "Harga harus lebih dari 0.";
  }

  if (!Number.isInteger(stock) || stock < 0) {
    errors.stock = "Stok harus angka bulat minimal 0.";
  }

  return errors;
}

function hasErrors(errors: ProductFormErrors): boolean {
  return Object.values(errors).some(Boolean);
}

function ProductFormSkeleton() {
  return (
    <section className="admin-page" aria-label="Loading product form">
      <div className="admin-product-form-skeleton hero" />
      <div className="admin-product-form-skeleton panel" />
    </section>
  );
}

function ProductForm() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const productID = searchParams.get("id")?.trim() ?? "";
  const isEditMode = Boolean(productID);

  const [form, setForm] = useState<ProductFormState>(EMPTY_FORM);
  const [categories, setCategories] = useState<Category[]>([]);
  const [currentProduct, setCurrentProduct] = useState<Product | null>(null);
  const [selectedImage, setSelectedImage] = useState<File | null>(null);
  const [errors, setErrors] = useState<ProductFormErrors>({});
  const [isLoading, setIsLoading] = useState(isEditMode);
  const [isSaving, setIsSaving] = useState(false);

  const currentImagePath = useMemo(() => {
    if (!currentProduct) {
      return "";
    }

    return getProductImage(currentProduct);
  }, [currentProduct]);

  useEffect(() => {
    let isActive = true;

    async function loadFormData() {
      setIsLoading(true);
      setErrors({});

      try {
        const [categoryResult, productResult] = await Promise.all([
          listCategories(),
          productID ? getProductById(productID) : Promise.resolve(null),
        ]);

        if (isActive) {
          setCategories(categoryResult);
          setCurrentProduct(productResult);

          if (productResult) {
            setForm(buildFormFromProduct(productResult));
          } else {
            setForm(EMPTY_FORM);
          }
        }
      } catch (loadError) {
        if (isActive) {
          setErrors({
            form: getProductErrorMessage(
              loadError,
              "Product form data could not be loaded.",
            ),
          });
        }
      } finally {
        if (isActive) {
          setIsLoading(false);
        }
      }
    }

    loadFormData();

    return () => {
      isActive = false;
    };
  }, [productID]);

  function updateField(field: keyof ProductFormState, value: string) {
    setForm((currentForm) => ({
      ...currentForm,
      [field]: value,
    }));

    setErrors((currentErrors) => ({
      ...currentErrors,
      [field]: undefined,
      form: undefined,
    }));
  }

  function handlePriceChange(value: string) {
    updateField("price", formatPriceInput(value));
  }

  function handleImageChange(fileList: FileList | null) {
    const file = fileList?.[0] ?? null;

    setSelectedImage(file);
    setErrors((currentErrors) => ({
      ...currentErrors,
      image: undefined,
      form: undefined,
    }));
  }

  function buildPayload(): ProductInput {
    return {
      category_id: form.categoryID,
      name: form.name.trim(),
      description: form.description.trim(),
      price: parseNumberInput(form.price),
      stock: Number(form.stock),
      image_url: currentProduct?.image_url ?? "",
    };
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    const nextErrors = validateForm(form);

    if (hasErrors(nextErrors)) {
      setErrors(nextErrors);
      return;
    }

    setIsSaving(true);
    setErrors({});

    try {
      const payload = buildPayload();
      const savedProduct = isEditMode
        ? await updateProduct(productID, payload)
        : await createProduct(payload);

      if (selectedImage) {
        await uploadProductImage(savedProduct.id, selectedImage);
      }

      navigate("/admin/products", { replace: true });
    } catch (submitError) {
      setErrors({
        form: getProductErrorMessage(
          submitError,
          isEditMode
            ? "Product could not be updated. Check the fields and try again."
            : "Product could not be created. Check the fields and try again.",
        ),
      });
    } finally {
      setIsSaving(false);
    }
  }

  if (isLoading) {
    return <ProductFormSkeleton />;
  }

  return (
    <section className="admin-page" aria-labelledby="admin-product-form-title">
      <Link className="admin-product-form-back" to="/admin/products">
        <ArrowLeft className="h-4 w-4" aria-hidden="true" />
        Back to products
      </Link>

      <header className="admin-page-header">
        <span>{isEditMode ? "Edit Product" : "Create Product"}</span>
        <h1 id="admin-product-form-title">
          {isEditMode ? "Edit item." : "New item."}
        </h1>
        <p>
          Fill the product fields, keep the catalog data clean, and upload one
          simple product image.
        </p>
      </header>

      {errors.form && (
        <div className="admin-products-notice is-error" role="alert">
          <AlertTriangle className="h-5 w-5" aria-hidden="true" />
          <span>{errors.form}</span>
        </div>
      )}

      <form className="admin-product-form" onSubmit={handleSubmit}>
        <div className="admin-product-form-main">
          <div className="admin-product-form-card">
            <div className="admin-product-form-heading">
              <h2>Product data</h2>
              <p>Required fields are marked by validation messages.</p>
            </div>

            <div className="admin-product-field">
              <label htmlFor="product-name">Nama</label>
              <input
                id="product-name"
                type="text"
                maxLength={150}
                value={form.name}
                onChange={(event) => updateField("name", event.target.value)}
                aria-invalid={Boolean(errors.name)}
                aria-describedby={
                  errors.name ? "product-name-error" : undefined
                }
              />
              {errors.name && (
                <p
                  id="product-name-error"
                  className="admin-product-field-error"
                >
                  {errors.name}
                </p>
              )}
            </div>

            <div className="admin-product-field">
              <label htmlFor="product-category">Kategori</label>
              <select
                id="product-category"
                value={form.categoryID}
                onChange={(event) =>
                  updateField("categoryID", event.target.value)
                }
                aria-invalid={Boolean(errors.categoryID)}
                aria-describedby={
                  errors.categoryID ? "product-category-error" : undefined
                }
              >
                <option value="">Pilih kategori</option>
                {categories.map((category) => (
                  <option key={category.id} value={category.id}>
                    {category.name}
                  </option>
                ))}
              </select>
              {errors.categoryID && (
                <p
                  id="product-category-error"
                  className="admin-product-field-error"
                >
                  {errors.categoryID}
                </p>
              )}
            </div>

            <div className="admin-product-form-grid">
              <div className="admin-product-field">
                <label htmlFor="product-price">Harga</label>
                <input
                  id="product-price"
                  type="text"
                  inputMode="numeric"
                  value={form.price}
                  onChange={(event) => handlePriceChange(event.target.value)}
                  placeholder="Rp0"
                  aria-invalid={Boolean(errors.price)}
                  aria-describedby={
                    errors.price ? "product-price-error" : undefined
                  }
                />
                {errors.price && (
                  <p
                    id="product-price-error"
                    className="admin-product-field-error"
                  >
                    {errors.price}
                  </p>
                )}
              </div>

              <div className="admin-product-field">
                <label htmlFor="product-stock">Stok</label>
                <input
                  id="product-stock"
                  type="number"
                  min={0}
                  step={1}
                  value={form.stock}
                  onChange={(event) => updateField("stock", event.target.value)}
                  aria-invalid={Boolean(errors.stock)}
                  aria-describedby={
                    errors.stock ? "product-stock-error" : undefined
                  }
                />
                {errors.stock && (
                  <p
                    id="product-stock-error"
                    className="admin-product-field-error"
                  >
                    {errors.stock}
                  </p>
                )}
              </div>
            </div>

            <div className="admin-product-field">
              <label htmlFor="product-description">Deskripsi</label>
              <textarea
                id="product-description"
                rows={6}
                value={form.description}
                onChange={(event) =>
                  updateField("description", event.target.value)
                }
              />
            </div>
          </div>
        </div>

        <aside className="admin-product-form-side">
          <div className="admin-product-form-card">
            <div className="admin-product-form-heading">
              <h2>Image</h2>
              <p>Upload one product image. Gallery sorting can come later.</p>
            </div>

            <div className="admin-product-image-preview">
              {currentImagePath ? (
                <ProductImage
                  key={currentImagePath}
                  src={currentImagePath}
                  alt={form.name || "Product"}
                  width={760}
                  height={520}
                  sizes="(max-width: 1040px) 100vw, 380px"
                />
              ) : (
                <div className="admin-product-image-empty">
                  <ImagePlus className="h-10 w-10" aria-hidden="true" />
                  <span>No image selected</span>
                </div>
              )}
            </div>

            <div className="admin-product-field">
              <label htmlFor="product-image">Image upload</label>
              <input
                id="product-image"
                type="file"
                accept="image/*"
                onChange={(event) => handleImageChange(event.target.files)}
              />
              {selectedImage && (
                <p className="admin-product-field-hint">
                  Selected: {selectedImage.name}
                </p>
              )}
              {errors.image && (
                <p className="admin-product-field-error">{errors.image}</p>
              )}
            </div>
          </div>

          <div className="admin-product-form-actions">
            <Link to="/admin/products">Cancel</Link>
            <button type="submit" disabled={isSaving}>
              {isSaving ? (
                <>
                  <Loader2
                    className="h-4 w-4 animate-spin"
                    aria-hidden="true"
                  />
                  Saving...
                </>
              ) : (
                <>
                  <Save className="h-4 w-4" aria-hidden="true" />
                  Save Product
                </>
              )}
            </button>
          </div>
        </aside>
      </form>
    </section>
  );
}

export default ProductForm;
