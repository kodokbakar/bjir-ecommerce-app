import { useEffect, useState, type FormEvent } from "react";

interface SearchBarProps {
  value: string;
  isLoading?: boolean;
  placeholder?: string;
  onSearch: (value: string) => void;
}

function SearchBar({
  value,
  isLoading = false,
  placeholder = "Search products...",
  onSearch,
}: SearchBarProps) {
  const [draft, setDraft] = useState(value);

  useEffect(() => {
    const timeoutId = window.setTimeout(() => {
      const nextSearch = draft.trim();

      if (nextSearch !== value.trim()) {
        onSearch(nextSearch);
      }
    }, 300);

    return () => {
      window.clearTimeout(timeoutId);
    };
  }, [draft, onSearch, value]);

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    onSearch(draft.trim());
  }

  function handleClear() {
    setDraft("");
    onSearch("");
  }

  return (
    <form className="search-bar" role="search" onSubmit={handleSubmit}>
      <label className="products-label" htmlFor="product-search">
        Search
      </label>

      <div className="search-bar-control">
        <svg
          aria-hidden="true"
          className="search-bar-icon"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2.4"
          strokeLinecap="round"
          strokeLinejoin="round"
        >
          <circle cx="11" cy="11" r="7" />
          <path d="m20 20-3.5-3.5" />
        </svg>

        <input
          className="search-bar-input"
          id="product-search"
          type="search"
          value={draft}
          onChange={(event) => setDraft(event.target.value)}
          placeholder={placeholder}
          autoComplete="off"
          aria-label="Search products"
        />

        {isLoading && (
          <span className="search-bar-loading" role="status" aria-label="Searching">
            Searching
          </span>
        )}

        {draft && (
          <button
            className="search-bar-clear"
            type="button"
            onClick={handleClear}
            aria-label="Clear search"
          >
            ×
          </button>
        )}
      </div>
    </form>
  );
}

export default SearchBar;
