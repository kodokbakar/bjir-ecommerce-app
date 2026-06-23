interface PaginationProps {
  page: number;
  limit: number;
  total: number;
  totalPages: number;
  onPageChange: (page: number) => void;
}

function buildPageNumbers(page: number, totalPages: number): number[] {
  const start = Math.max(1, page - 2);
  const end = Math.min(totalPages, page + 2);

  return Array.from({ length: end - start + 1 }, (_, index) => start + index);
}

function Pagination({
  page,
  limit,
  total,
  totalPages,
  onPageChange,
}: PaginationProps) {
  if (totalPages <= 0) {
    return null;
  }

  const safePage = Math.min(Math.max(page, 1), totalPages);
  const startItem = total === 0 ? 0 : (safePage - 1) * limit + 1;
  const endItem = Math.min(safePage * limit, total);
  const pages = buildPageNumbers(safePage, totalPages);

  return (
    <nav className="pagination" aria-label="Product pagination">
      <p className="pagination-summary">
        Showing {startItem}-{endItem} of {total} products
      </p>

      <div className="pagination-controls">
        <button
          className="pagination-button"
          type="button"
          disabled={safePage <= 1}
          onClick={() => onPageChange(safePage - 1)}
        >
          Previous
        </button>

        <span className="pagination-mobile-label">
          {safePage} / {totalPages}
        </span>

        <div className="pagination-controls pagination-pages">
          {pages[0] > 1 && (
            <>
              <button
                className="pagination-page"
                type="button"
                onClick={() => onPageChange(1)}
              >
                1
              </button>
              {pages[0] > 2 && <span aria-hidden="true">...</span>}
            </>
          )}

          {pages.map((item) => (
            <button
              key={item}
              className="pagination-page"
              type="button"
              aria-current={item === safePage ? "page" : undefined}
              onClick={() => onPageChange(item)}
            >
              {item}
            </button>
          ))}

          {pages[pages.length - 1] < totalPages && (
            <>
              {pages[pages.length - 1] < totalPages - 1 && (
                <span aria-hidden="true">...</span>
              )}
              <button
                className="pagination-page"
                type="button"
                onClick={() => onPageChange(totalPages)}
              >
                {totalPages}
              </button>
            </>
          )}
        </div>

        <button
          className="pagination-button"
          type="button"
          disabled={safePage >= totalPages}
          onClick={() => onPageChange(safePage + 1)}
        >
          Next
        </button>
      </div>
    </nav>
  );
}

export default Pagination;