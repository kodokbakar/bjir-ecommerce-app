function AdminOrders() {
  return (
    <section className="admin-page" aria-labelledby="admin-orders-title">
      <header className="admin-page-header">
        <span>Admin Orders</span>
        <h1 id="admin-orders-title">Orders.</h1>
        <p>
          Order management page shell. Admin order table comes in the next
          issue.
        </p>
      </header>

      <div className="admin-panel">
        <span>Next</span>
        <strong>Order table</strong>
        <p>
          This route is ready for customer order review and status operations.
        </p>
      </div>
    </section>
  );
}

export default AdminOrders;
