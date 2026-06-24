function AdminCategories() {
  return (
    <section className="admin-page" aria-labelledby="admin-categories-title">
      <header className="admin-page-header">
        <span>Admin Categories</span>
        <h1 id="admin-categories-title">Categories.</h1>
        <p>
          Category management page shell. CRUD table comes in the next issue.
        </p>
      </header>

      <div className="admin-panel">
        <span>Next</span>
        <strong>Category table</strong>
        <p>This route is ready for category listing and hierarchy actions.</p>
      </div>
    </section>
  );
}

export default AdminCategories;
