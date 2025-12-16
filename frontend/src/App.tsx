import { Layout, Menu } from 'antd'
import { Link, Route, Routes } from 'react-router-dom'
import Home from './pages/Home'
import Catalog from './pages/Catalog'
import ProductDetails from './pages/ProductDetails'
import Cart from './pages/Cart'
import Checkout from './pages/Checkout'
import Login from './pages/Login'
import Register from './pages/Register'
import Orders from './pages/Orders'
import AdminDashboard from './pages/AdminDashboard'
import { CartProvider } from './store/CartContext'

const { Header, Content, Footer } = Layout

export default function App() {
  return (
    <CartProvider>
      <Layout style={{ minHeight: '100vh' }}>
        <Header>
          <Menu theme="dark" mode="horizontal" selectable={false}>
            <Menu.Item key="home"><Link to="/">Начало</Link></Menu.Item>
            <Menu.Item key="catalog"><Link to="/catalog">Каталог</Link></Menu.Item>
            <Menu.Item key="cart"><Link to="/cart">Количка</Link></Menu.Item>
            <Menu.Item key="orders"><Link to="/orders">Моите поръчки</Link></Menu.Item>
            <Menu.Item key="admin"><Link to="/admin">Админ</Link></Menu.Item>
            <Menu.Item key="login" style={{ marginLeft: 'auto' }}><Link to="/login">Вход</Link></Menu.Item>
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
        <Footer style={{ textAlign: 'center' }}>Онлайн магазин за мебели © 2025</Footer>
      </Layout>
    </CartProvider>
  )
}

