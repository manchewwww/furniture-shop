import { Button, Card } from "antd";
import { useLocation, useNavigate } from "react-router-dom";

export default function AdminNav() {
  const nav = useNavigate();
  const { pathname } = useLocation();
  const items = [
    { path: "/admin/departments", label: "Departments" },
    { path: "/admin/categories", label: "Categories" },
    { path: "/admin/products", label: "Products" },
    { path: "/admin/orders", label: "Orders" },
  ];
  return (
    <Card style={{ marginBottom: 16 }}>
      {items.map((it) => (
        <Button
          key={it.path}
          type="link"
          style={{ fontWeight: pathname.startsWith(it.path) ? 600 : 400 }}
          onClick={() => nav(it.path)}
        >
          {it.label}
        </Button>
      ))}
    </Card>
  );
}
