import { Layout, Menu } from "antd";
import { Link, Route, Routes, useNavigate } from "react-router-dom";
import Home from "./pages/Home";
import Catalog from "./pages/Catalog";
import ProductDetails from "./pages/ProductDetails";
import Cart from "./pages/Cart";
import Checkout from "./pages/Checkout";
import Login from "./pages/Login";
import Register from "./pages/Register";
import Orders from "./pages/Orders";
import AdminDashboard from "./pages/AdminDashboard";
import { CartProvider } from "./store/CartContext";
import { useAuth } from "./store/AuthContext";

const { Header, Content, Footer } = Layout;

export default function App() {
  const { isAuthenticated, logout } = useAuth();
  const nav = useNavigate();
  return (
    <CartProvider>
      <Layout style={{ minHeight: "100vh" }}>
        <Header>
          <Menu theme="dark" mode="horizontal" selectable={false}>
            <Menu.Item key="home">
              <Link to="/">Home</Link>
            </Menu.Item>
            <Menu.Item key="catalog">
              <Link to="/catalog">Catalog</Link>
            </Menu.Item>
            <Menu.Item key="cart">
              <Link to="/cart">Cart</Link>
            </Menu.Item>
            {isAuthenticated ? (
              <>
                <Menu.Item key="orders">
                  <Link to="/orders">My Orders</Link>
                </Menu.Item>
                <Menu.Item key="admin">
                  <Link to="/admin">Admin</Link>
                </Menu.Item>
                <Menu.Item
                  key="logout"
                  style={{ marginLeft: "auto" }}
                  onClick={() => {
                    logout();
                    nav("/");
                  }}
                >
                  Logout
                </Menu.Item>
              </>
            ) : (
              <>
                <Menu.Item key="login" style={{ marginLeft: "auto" }}>
                  <Link to="/login">Login</Link>
                </Menu.Item>
                <Menu.Item key="register">
                  <Link to="/register">Register</Link>
                </Menu.Item>
              </>
            )}
          </Menu>
        </Header>
        <Content style={{ padding: 24 }}>
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/catalog" element={<Catalog />} />
            <Route path="/product/:id" element={<ProductDetails />} />
            <Route path="/cart" element={<Cart />} />
            <Route path="/checkout" element={<Checkout />} />
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
            <Route path="/orders" element={<Orders />} />
            <Route path="/admin" element={<AdminDashboard />} />
          </Routes>
        </Content>
        <Footer style={{ textAlign: "center" }}>
          Магазин за мебели © 2025
        </Footer>
      </Layout>
    </CartProvider>
  );
}
