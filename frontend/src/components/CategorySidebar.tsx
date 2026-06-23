import type { Category } from "../types/product";

interface CategorySidebarProps {
    categories: Category[];
    activeCategory: string;
    isLoading: boolean;
    error: string | null;
    onSelect: (slug: string) => void;
    onRetry: () => void;
}

function CategoryCountBadge({ count }: { count?: number }) {
    if (typeof count !== "number") {
        return null;
    }

    return <span className="category-count-badge">{count}</span>;
}

function CategorySidebar({
    categories,
    activeCategory,
    isLoading,
    error,
    onSelect,
    onRetry,
}: CategorySidebarProps) {
    return (
        <aside className="category-sidebar" aria-label="Product categories">
            <div className="category-sidebar-header">
                <span>Category Index</span>
            </div>

            <button
                className={`category-sidebar-item ${!activeCategory ? "is-active" : ""}`}
                type="button"
                onClick={() => onSelect("")}
            >
                <span>All Products</span>
            </button>

            {isLoading && (
                <div className="category-sidebar-note" role="status">
                    Loading categories...
                </div>
            )}

            {error && (
                <div className="category-sidebar-error" role="status">
                    <span>Categories failed.</span>
                    <button type="button" onClick={onRetry}>
                        Retry
                    </button>
                </div>
            )}

            {!isLoading && !error && categories.length === 0 && (
                <div className="category-sidebar-note">No categories yet.</div>
            )}

            {!isLoading &&
                !error &&
                categories.map((category) => (
                    <button
                        className={`category-sidebar-item ${activeCategory === category.slug ? "is-active" : ""
                            }`}
                        key={category.id}
                        type="button"
                        onClick={() => onSelect(category.slug)}
                    >
                        <span>{category.name}</span>
                        <CategoryCountBadge count={category.product_count} />
                    </button>
                ))}
        </aside>
    );
}

export default CategorySidebar;
