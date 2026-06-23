import type { Category } from "../types/product";

interface CategoryBarProps {
    categories: Category[];
    activeCategory: string;
    isLoading: boolean;
    error: string | null;
    onSelect: (slug: string) => void;
    onRetry: () => void;
}

function CategoryBar({
    categories,
    activeCategory,
    isLoading,
    error,
    onSelect,
    onRetry,
}: CategoryBarProps) {
    return (
        <div className="category-bar" aria-label="Product categories">
            <button
                className={`category-chip ${!activeCategory ? "is-active" : ""}`}
                type="button"
                onClick={() => onSelect("")}
            >
                All
            </button>

            {categories.map((category) => (
                <button
                    className={`category-chip ${activeCategory === category.slug ? "is-active" : ""
                        }`}
                    key={category.id}
                    type="button"
                    onClick={() => onSelect(category.slug)}
                >
                    {category.name}
                    {typeof category.product_count === "number" && (
                        <span>{category.product_count}</span>
                    )}
                </button>
            ))}

            {isLoading && <span className="category-chip is-muted">Loading...</span>}

            {error && (
                <button className="category-chip is-error" type="button" onClick={onRetry}>
                    Retry categories
                </button>
            )}
        </div>
    );
}

export default CategoryBar;
