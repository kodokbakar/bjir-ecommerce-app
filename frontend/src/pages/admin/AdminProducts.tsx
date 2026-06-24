function AdminProducts() {
  return (
    <section className="admin-page" aria-labelledby="admin-products-title">
      <header className="admin-page-header">
        <span>Admin Products</span>
        <h1 id="admin-products-title">Products.</h1>
        <p>
          Product management page shell. CRUD table comes in the next issue.
        </p>
      </header>

      <div className="admin-panel">
        <span>Next</span>
        <strong>Product table</strong>
        <p>This route is ready for admin product listing and actions.</p>
      </div>
    </section>
  );
}

export default AdminProducts;
