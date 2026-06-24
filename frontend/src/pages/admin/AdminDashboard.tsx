function AdminDashboard() {
  return (
    <section className="admin-page" aria-labelledby="admin-dashboard-title">
      <header className="admin-page-header">
        <span>Admin Dashboard</span>
        <h1 id="admin-dashboard-title">Control room.</h1>
        <p>
          Data-focused overview for catalog, category, and order operations.
        </p>
      </header>

      <div className="admin-panel-grid">
        <article className="admin-panel">
          <span>Catalog</span>
          <strong>Products</strong>
          <p>Manage buyer-facing product data from the admin section.</p>
        </article>

        <article className="admin-panel">
          <span>Taxonomy</span>
          <strong>Categories</strong>
          <p>Keep storefront grouping clean and easy to browse.</p>
        </article>

        <article className="admin-panel">
          <span>Fulfillment</span>
          <strong>Orders</strong>
          <p>Review customer orders and payment state.</p>
        </article>
      </div>
    </section>
  );
}

export default AdminDashboard;
