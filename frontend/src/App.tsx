import { Navigate, Route, Routes } from "react-router-dom";

import Layout from "./components/layout/Layout";
import ProtectedRoute from "./components/ProtectedRoute";
import Cart from "./pages/Cart";
import CategoryProducts from "./pages/CategoryProducts";
import Checkout from "./pages/Checkout";
import Dashboard from "./pages/Dashboard";
import Login from "./pages/Login";
import MyOrders from "./pages/MyOrders";
import OrderDetail from "./pages/OrderDetail";
import ProductDetail from "./pages/ProductDetail";
import Products from "./pages/Products";
import Register from "./pages/Register";
import Profile from "./pages/Profile";
import Payment from "./pages/Payment";

import AdminProtectedRoute from "./components/AdminProtectedRoute";
import AdminLayout from "./components/admin/AdminLayout";
import AdminCategories from "./pages/admin/AdminCategories";
import AdminDashboard from "./pages/admin/AdminDashboard";
import AdminOrders from "./pages/admin/AdminOrders";
import AdminProducts from "./pages/admin/AdminProducts";
import ProductForm from "./pages/admin/ProductForm";

import GlobalNetworkError from "./components/GlobalNetworkError";
import ToastViewport from "./components/ToastViewport";

import LandingPage from "./pages/LandingPage";

function App() {
  return (
    <>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/" element={<LandingPage />} />

        <Route element={<ProtectedRoute />}>
          <Route
            path="/dashboard"
            element={
              <Layout>
                <Dashboard />
              </Layout>
            }
          />

          <Route
            path="/products"
            element={
              <Layout>
                <Products />
              </Layout>
            }
          />

          <Route
            path="/products/:slug"
            element={
              <Layout>
                <ProductDetail />
              </Layout>
            }
          />

          <Route
            path="/categories/:slug"
            element={
              <Layout>
                <CategoryProducts />
              </Layout>
            }
          />

          <Route
            path="/cart"
            element={
              <Layout>
                <Cart />
              </Layout>
            }
          />

          <Route
            path="/checkout"
            element={
              <Layout>
                <Checkout />
              </Layout>
            }
          />

          <Route
            path="/orders"
            element={
              <Layout>
                <MyOrders />
              </Layout>
            }
          />

          <Route
            path="/orders/:id"
            element={
              <Layout>
                <OrderDetail />
              </Layout>
            }
          />

          <Route
            path="/profile"
            element={
              <Layout>
                <Profile />
              </Layout>
            }
          />

          <Route
            path="/payment"
            element={
              <Layout>
                <Payment />
              </Layout>
            }
          />

          <Route
            path="/settings"
            element={<Navigate to="/profile" replace />}
          />
        </Route>

        <Route element={<AdminProtectedRoute />}>
          <Route
            path="/admin"
            element={<Navigate to="/admin/dashboard" replace />}
          />

          <Route
            path="/admin/dashboard"
            element={
              <AdminLayout>
                <AdminDashboard />
              </AdminLayout>
            }
          />

          <Route
            path="/admin/products/new"
            element={
              <AdminLayout>
                <ProductForm />
              </AdminLayout>
            }
          />

          <Route
            path="/admin/products/edit"
            element={
              <AdminLayout>
                <ProductForm />
              </AdminLayout>
            }
          />

          <Route
            path="/admin/products"
            element={
              <AdminLayout>
                <AdminProducts />
              </AdminLayout>
            }
          />

          <Route
            path="/admin/categories"
            element={
              <AdminLayout>
                <AdminCategories />
              </AdminLayout>
            }
          />

          <Route
            path="/admin/orders"
            element={
              <AdminLayout>
                <AdminOrders />
              </AdminLayout>
            }
          />
        </Route>
      </Routes>

      <GlobalNetworkError />
      <ToastViewport />
    </>
  );
}

export default App;
