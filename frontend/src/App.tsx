import { Layout, Menu, Select } from "antd";
import { Link, Navigate, Route, Routes, useNavigate } from "react-router-dom";
import Home from "./pages/Home";
import Catalog from "./pages/Catalog";
import ProductDetails from "./pages/ProductDetails";
import Cart from "./pages/Cart";
import Checkout from "./pages/Checkout";
import Login from "./pages/Login";
import Register from "./pages/Register";
import Orders from "./pages/Orders";
import PaymentSuccess from "./pages/PaymentSuccess";
import PaymentCancel from "./pages/PaymentCancel";
import AdminLayout from "./components/AdminLayout";
import AdminDepartments from "./pages/AdminDepartments";
import AdminCategories from "./pages/AdminCategories";
import AdminProducts from "./pages/AdminProducts";
import AdminOrders from "./pages/AdminOrders";
import { useCart } from "./store/CartContext";
import { useAuth } from "./store/AuthContext";
import { useI18n } from "./store/I18nContext";
import {
  RequireAuth,
  RequireRole,
  ForbidAuth,
  ForbidRole,
} from "./components/RouteGuards";

const { Header, Content, Footer } = Layout;

export default function App() {
  const { isAuthenticated, user, logout } = useAuth();
  const { t, lang, setLang } = useI18n();
  const nav = useNavigate();
  const { clear } = useCart();
  const handleLogout = async () => {
    try {
      await clear();
    } catch {}
    logout();
    nav("/");
  };
  return (
    <Layout style={{ minHeight: "100vh" }}>
      <Header>
        <Menu theme="dark" mode="horizontal" selectable={false}>
          <Menu.Item key="home">
            <Link to="/">{t("nav.home")}</Link>
          </Menu.Item>
          <Menu.Item key="catalog">
            <Link to="/catalog">{t("nav.catalog")}</Link>
          </Menu.Item>
          <Menu.Item key="cart">
            <Link to="/cart">{t("nav.cart")}</Link>
          </Menu.Item>
          {isAuthenticated ? (
            <>
              {user?.role !== "admin" && (
                <Menu.Item key="orders">
                  <Link to="/orders">{t("nav.orders")}</Link>
                </Menu.Item>
              )}
              {user?.role === "admin" && (
                <Menu.Item key="admin">
                  <Link to="/admin">{t("nav.admin")}</Link>
                </Menu.Item>
              )}
              <Menu.Item
                key="logout"
                style={{ marginLeft: "auto" }}
                onClick={handleLogout}
              >
                {t("nav.logout")}
              </Menu.Item>
            </>
          ) : (
            <>
              <Menu.Item key="login" style={{ marginLeft: "auto" }}>
                <Link to="/login">{t("nav.login")}</Link>
              </Menu.Item>
              <Menu.Item key="register">
                <Link to="/register">{t("nav.register")}</Link>
              </Menu.Item>
            </>
          )}
          <Menu.Item key="lang" style={{ marginLeft: 12 }}>
            <Select
              value={lang}
              style={{ width: 140 }}
              onChange={(v) => setLang(v as any)}
              options={[
                { value: "en", label: "English" },
                { value: "bg", label: "Български" },
              ]}
            />
          </Menu.Item>
        </Menu>
      </Header>
      <Content style={{ padding: 24 }}>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/catalog" element={<Catalog />} />
          <Route path="/product/:id" element={<ProductDetails />} />
          <Route path="/cart" element={<Cart />} />
          <Route
            path="/checkout"
            element={
              <RequireAuth>
                <Checkout />
              </RequireAuth>
            }
          />
          <Route path="/payment/success" element={<PaymentSuccess />} />
          <Route path="/payment/cancel" element={<PaymentCancel />} />
          <Route
            path="/login"
            element={
              <ForbidAuth>
                <Login />
              </ForbidAuth>
            }
          />
          <Route
            path="/register"
            element={
              <ForbidAuth>
                <Register />
              </ForbidAuth>
            }
          />
          <Route
            path="/orders"
            element={
              <RequireAuth>
                <ForbidRole role="admin">
                  <Orders />
                </ForbidRole>
              </RequireAuth>
            }
          />
          <Route
            path="/admin"
            element={<RequireRole role="admin">{<AdminLayout />}</RequireRole>}
          >
            <Route index element={<Navigate to="products" replace />} />
            <Route path="departments" element={<AdminDepartments />} />
            <Route path="categories" element={<AdminCategories />} />
            <Route path="products" element={<AdminProducts />} />
            <Route path="orders" element={<AdminOrders />} />
          </Route>
        </Routes>
      </Content>
      <Footer style={{ textAlign: "center" }}>Магазин за мебели © 2025</Footer>
    </Layout>
  );
}
