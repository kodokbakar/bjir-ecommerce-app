import { useEffect, useMemo, useState, type FormEvent } from "react";
import {
  AlertTriangle,
  Edit3,
  FolderTree,
  Plus,
  Save,
  Tags,
  Trash2,
  X,
  RefreshCcw,
} from "lucide-react";

import {
  createCategory,
  deleteCategory,
  getCategoryErrorMessage,
  listAdminCategories,
  updateCategory,
} from "../../services/categoryService";
import type {
  Category,
  CategoryInput,
  CategoryListMeta,
} from "../../types/product";
import { useToast } from "../../context/toast"

const ADMIN_CATEGORY_LIMIT = 10;

const EMPTY_META: CategoryListMeta = {
  page: 1,
  limit: ADMIN_CATEGORY_LIMIT,
  total: 0,
  total_pages: 0,
};

const EMPTY_FORM: CategoryFormState = {
  name: "",
  parentID: "",
  description: "",
  imageUrl: "",
};

interface CategoryFormState {
  name: string;
  parentID: string;
  description: string;
  imageUrl: string;
}

interface CategoryFormErrors {
  name?: string;
  form?: string;
}

type CategoryModalMode = "create" | "edit";

function getParentName(category: Category, categories: Category[]): string {
  if (!category.parent_id) {
    return "Root category";
  }

  return (
    category.parent?.name ||
    categories.find((item) => item.id === category.parent_id)?.name ||
    "Unknown parent"
  );
}

function buildCategoryForm(category: Category): CategoryFormState {
  return {
    name: category.name,
    parentID: category.parent_id ?? "",
    description: category.description ?? "",
    imageUrl: category.image_url ?? "",
  };
}

function buildCategoryInput(form: CategoryFormState): CategoryInput {
  return {
    parent_id: form.parentID || null,
    name: form.name.trim(),
    description: form.description.trim(),
    image_url: form.imageUrl.trim(),
  };
}

function validateCategoryForm(form: CategoryFormState): CategoryFormErrors {
  const errors: CategoryFormErrors = {};

  if (!form.name.trim()) {
    errors.name = "Category name is required.";
  }

  return errors;
}

function hasErrors(errors: CategoryFormErrors): boolean {
  return Object.values(errors).some(Boolean);
}

function collectNestedChildIDs(
  children: Category[] | undefined,
  blockedIDs: Set<string>,
) {
  children?.forEach((child) => {
    blockedIDs.add(child.id);
    collectNestedChildIDs(child.children, blockedIDs);
  });
}

function collectFlatChildIDs(
  parentID: string,
  categories: Category[],
  blockedIDs: Set<string>,
) {
  categories.forEach((category) => {
    if (category.parent_id !== parentID || blockedIDs.has(category.id)) {
      return;
    }

    blockedIDs.add(category.id);
    collectNestedChildIDs(category.children, blockedIDs);
    collectFlatChildIDs(category.id, categories, blockedIDs);
  });
}

function getBlockedParentIDs(
  category: Category | null,
  categories: Category[],
): Set<string> {
  const blockedIDs = new Set<string>();

  if (!category) {
    return blockedIDs;
  }

  blockedIDs.add(category.id);
  collectNestedChildIDs(category.children, blockedIDs);
  collectFlatChildIDs(category.id, categories, blockedIDs);

  return blockedIDs;
}

function CategorySkeleton() {
  return (
    <div className="admin-categories-list" aria-label="Loading categories">
      {Array.from({ length: 5 }, (_, index) => (
        <div className="admin-categories-skeleton-row" key={index}>
          <div className="admin-categories-skeleton-line short" />
          <div className="admin-categories-skeleton-line" />
          <div className="admin-categories-skeleton-line tiny" />
        </div>
      ))}
    </div>
  );
}

function AdminCategories() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [parentOptions, setParentOptions] = useState<Category[]>([]);
  const [meta, setMeta] = useState<CategoryListMeta>(EMPTY_META);
  const [page, setPage] = useState(1);
  const [reloadKey, setReloadKey] = useState(0);
  const { showToast } = useToast();

  const [modalMode, setModalMode] = useState<CategoryModalMode | null>(null);
  const [editingCategory, setEditingCategory] = useState<Category | null>(null);
  const [form, setForm] = useState<CategoryFormState>(EMPTY_FORM);
  const [formErrors, setFormErrors] = useState<CategoryFormErrors>({});

  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [deletingCategoryID, setDeletingCategoryID] = useState<string | null>(
    null,
  );
  const [error, setError] = useState<string | null>(null);

  const hasCategories = categories.length > 0;
  const isModalOpen = modalMode !== null;

  const blockedParentIDs = useMemo(
    () => getBlockedParentIDs(editingCategory, parentOptions),
    [editingCategory, parentOptions],
  );

  const availableParentOptions = useMemo(
    () =>
      parentOptions.filter((category) => !blockedParentIDs.has(category.id)),
    [blockedParentIDs, parentOptions],
  );

  const pageSummary = useMemo(() => {
    if (meta.total === 0) {
      return "0 category";
    }

    const start = (meta.page - 1) * meta.limit + 1;
    const end = Math.min(meta.page * meta.limit, meta.total);

    return `${start}-${end} of ${meta.total} categories`;
  }, [meta]);

  useEffect(() => {
    let isActive = true;

    async function loadCategories() {
      setIsLoading(true);
      setError(null);

      try {
        const [pageResult, optionsResult] = await Promise.all([
          listAdminCategories({
            page,
            limit: ADMIN_CATEGORY_LIMIT,
          }),
          listAdminCategories({
            page: 1,
            limit: 100,
          }),
        ]);

        if (isActive) {
          setCategories(pageResult.data);
          setMeta(pageResult.meta);
          setParentOptions(optionsResult.data);
        }
      } catch (loadError) {
        if (isActive) {
          setCategories([]);
          setMeta({
            ...EMPTY_META,
            page,
          });
          setParentOptions([]);
          setError(
            getCategoryErrorMessage(
              loadError,
              "Admin category list could not be loaded.",
            ),
          );
        }
      } finally {
        if (isActive) {
          setIsLoading(false);
        }
      }
    }

    loadCategories();

    return () => {
      isActive = false;
    };
  }, [page, reloadKey]);

  function openCreateModal() {
    setModalMode("create");
    setEditingCategory(null);
    setForm(EMPTY_FORM);
    setFormErrors({});
  }

  function openEditModal(category: Category) {
    setModalMode("edit");
    setEditingCategory(category);
    setForm(buildCategoryForm(category));
    setFormErrors({});
  }

  function closeModal() {
    if (isSaving) {
      return;
    }

    setModalMode(null);
    setEditingCategory(null);
    setForm(EMPTY_FORM);
    setFormErrors({});
  }

  function updateField(field: keyof CategoryFormState, value: string) {
    setForm((currentForm) => ({
      ...currentForm,
      [field]: value,
    }));

    setFormErrors((currentErrors) => ({
      ...currentErrors,
      [field]: undefined,
      form: undefined,
    }));
  }

  function handlePageChange(nextPage: number) {
    setPage(Math.max(1, nextPage));
  }

  function handleRetry() {
    setReloadKey((current) => current + 1);
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    const nextErrors = validateCategoryForm(form);

    if (hasErrors(nextErrors)) {
      setFormErrors(nextErrors);
      return;
    }

    setIsSaving(true);
    setFormErrors({});
    setError(null);

    try {
      const payload = buildCategoryInput(form);

      if (modalMode === "edit" && editingCategory) {
        await updateCategory(editingCategory.id, payload);
        showToast({
          type: "success",
          message: `${payload.name} updated.`,
        });
      } else {
        await createCategory(payload);
        showToast({
          type: "success",
          message: `${payload.name} created.`,
        });
        setPage(1);
      }

      closeModal();
      setReloadKey((current) => current + 1);
    } catch (submitError) {
      setFormErrors({
        form: getCategoryErrorMessage(
          submitError,
          modalMode === "edit"
            ? "Category could not be updated. Check the fields and try again."
            : "Category could not be created. Check the fields and try again.",
        ),
      });
    } finally {
      setIsSaving(false);
    }
  }

  async function handleDeleteCategory(category: Category) {
    const confirmed = window.confirm(
      `Delete "${category.name}"? Categories with children may be rejected by the backend.`,
    );

    if (!confirmed) {
      return;
    }

    setDeletingCategoryID(category.id);
    setError(null);

    try {
      await deleteCategory(category.id);
      showToast({
        type: "success",
        message: `${category.name} deleted.`,
      });
      setReloadKey((current) => current + 1);
    } catch (deleteError) {
      showToast(
        {
          type: "error",
          message: getCategoryErrorMessage(
            deleteError,
            "Category could not be deleted. Remove child categories first if needed.",
          ),
        },
        { duration: 6000 },
      );
    } finally {
      setDeletingCategoryID(null);
    }
  }

  return (
    <section className="admin-page" aria-labelledby="admin-categories-title">
      <header className="admin-page-header">
        <span>Admin Categories</span>
        <h1 id="admin-categories-title">Categories.</h1>
        <p>
          Manage category names, parent hierarchy, descriptions, and storefront
          grouping.
        </p>
      </header>

      <div className="admin-categories-toolbar">
        <div>
          <span>Category Index</span>
          <strong>{pageSummary}</strong>
        </div>

        <button type="button" onClick={openCreateModal}>
          <Plus className="h-5 w-5" aria-hidden="true" />
          Add Category
        </button>
      </div>

      {error && hasCategories && (
        <div className="admin-products-notice is-error" role="alert">
          <AlertTriangle className="h-5 w-5" aria-hidden="true" />
          <span>{error}</span>
          <button type="button" onClick={handleRetry}>
            <RefreshCcw className="h-4 w-4" aria-hidden="true" />
            Retry
          </button>
        </div>
      )}

      {isLoading ? (
        <CategorySkeleton />
      ) : error && !hasCategories ? (
        <div className="admin-categories-empty" role="alert">
          <div>
            <AlertTriangle
              className="mx-auto mb-3 h-10 w-10"
              aria-hidden="true"
            />
            <h2>Category list jammed.</h2>
            <p>{error}</p>
            <button type="button" onClick={handleRetry}>
              <RefreshCcw className="h-5 w-5" aria-hidden="true" />
              Retry
            </button>
          </div>
        </div>
      ) : !hasCategories ? (
        <div className="admin-categories-empty">
          <div>
            <Tags className="mx-auto mb-3 h-10 w-10" aria-hidden="true" />
            <h2>No categories found.</h2>
            <p>
              Create the first category before assigning products to shelves.
            </p>
            <button type="button" onClick={openCreateModal}>
              <Plus className="h-5 w-5" aria-hidden="true" />
              Add Category
            </button>
          </div>
        </div>
      ) : (
        <>
          <div className="admin-categories-table" aria-label="Category list">
            <div className="admin-categories-table-head" aria-hidden="true">
              <span>Name</span>
              <span>Slug</span>
              <span>Description</span>
              <span>Parent</span>
              <span>Actions</span>
            </div>

            <div className="admin-categories-list">
              {categories.map((category) => {
                const isDeleting = deletingCategoryID === category.id;

                return (
                  <article className="admin-categories-row" key={category.id}>
                    <div className="admin-categories-name-cell">
                      <span className="admin-categories-icon">
                        <FolderTree className="h-5 w-5" aria-hidden="true" />
                      </span>
                      <strong>{category.name}</strong>
                    </div>

                    <span className="admin-categories-muted">
                      {category.slug}
                    </span>

                    <span className="admin-categories-muted">
                      {category.description || "No description"}
                    </span>

                    <span className="admin-categories-parent">
                      {getParentName(category, parentOptions)}
                    </span>

                    <div className="admin-categories-actions">
                      <button
                        type="button"
                        onClick={() => openEditModal(category)}
                      >
                        <Edit3 className="h-4 w-4" aria-hidden="true" />
                        Edit
                      </button>

                      <button
                        type="button"
                        disabled={isDeleting}
                        onClick={() => void handleDeleteCategory(category)}
                      >
                        <Trash2 className="h-4 w-4" aria-hidden="true" />
                        {isDeleting ? "Deleting..." : "Delete"}
                      </button>
                    </div>
                  </article>
                );
              })}
            </div>
          </div>

          <nav
            className="admin-categories-pagination"
            aria-label="Category pagination"
          >
            <button
              className="pagination-button"
              type="button"
              disabled={meta.page <= 1}
              onClick={() => handlePageChange(meta.page - 1)}
            >
              Previous
            </button>

            <span>
              {meta.page} / {Math.max(meta.total_pages, 1)}
            </span>

            <button
              className="pagination-button"
              type="button"
              disabled={meta.page >= meta.total_pages}
              onClick={() => handlePageChange(meta.page + 1)}
            >
              Next
            </button>
          </nav>
        </>
      )}

      {isModalOpen && (
        <div className="admin-category-modal" role="presentation">
          <button
            className="admin-category-modal-backdrop"
            type="button"
            aria-label="Close category form"
            onClick={closeModal}
          />

          <section
            className="admin-category-dialog"
            role="dialog"
            aria-modal="true"
            aria-labelledby="admin-category-dialog-title"
          >
            <form onSubmit={handleSubmit}>
              <header className="admin-category-dialog-header">
                <div>
                  <span>
                    {modalMode === "edit" ? "Edit Category" : "Create Category"}
                  </span>
                  <h2 id="admin-category-dialog-title">
                    {modalMode === "edit"
                      ? "Update category."
                      : "New category."}
                  </h2>
                </div>

                <button type="button" onClick={closeModal} disabled={isSaving}>
                  <X className="h-5 w-5" aria-hidden="true" />
                  <span className="sr-only">Close</span>
                </button>
              </header>

              {formErrors.form && (
                <div className="admin-category-form-error" role="alert">
                  <AlertTriangle className="h-5 w-5" aria-hidden="true" />
                  <span>{formErrors.form}</span>
                </div>
              )}

              <div className="admin-category-form-body">
                <div className="admin-product-field">
                  <label htmlFor="category-name">Name</label>
                  <input
                    id="category-name"
                    type="text"
                    value={form.name}
                    onChange={(event) =>
                      updateField("name", event.target.value)
                    }
                    aria-invalid={Boolean(formErrors.name)}
                    aria-describedby={
                      formErrors.name ? "category-name-error" : undefined
                    }
                  />
                  {formErrors.name && (
                    <p
                      id="category-name-error"
                      className="admin-product-field-error"
                    >
                      {formErrors.name}
                    </p>
                  )}
                </div>

                <div className="admin-product-field">
                  <label htmlFor="category-parent">Parent category</label>
                  <select
                    id="category-parent"
                    value={form.parentID}
                    onChange={(event) =>
                      updateField("parentID", event.target.value)
                    }
                  >
                    <option value="">Root category</option>
                    {availableParentOptions.map((category) => (
                      <option key={category.id} value={category.id}>
                        {category.name}
                      </option>
                    ))}
                  </select>
                </div>

                <div className="admin-product-field">
                  <label htmlFor="category-description">Description</label>
                  <textarea
                    id="category-description"
                    rows={4}
                    value={form.description}
                    onChange={(event) =>
                      updateField("description", event.target.value)
                    }
                  />
                </div>

                <div className="admin-product-field">
                  <label htmlFor="category-image-url">Image URL</label>
                  <input
                    id="category-image-url"
                    type="text"
                    value={form.imageUrl}
                    onChange={(event) =>
                      updateField("imageUrl", event.target.value)
                    }
                    placeholder="/uploads/category.jpg"
                  />
                </div>
              </div>

              <footer className="admin-category-dialog-actions">
                <button type="button" onClick={closeModal} disabled={isSaving}>
                  Cancel
                </button>
                <button type="submit" disabled={isSaving}>
                  <Save className="h-4 w-4" aria-hidden="true" />
                  {isSaving ? "Saving..." : "Save Category"}
                </button>
              </footer>
            </form>
          </section>
        </div>
      )}
    </section>
  );
}

export default AdminCategories;
