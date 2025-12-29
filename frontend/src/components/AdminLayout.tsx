import { Outlet } from "react-router-dom";
import { Card } from "antd";
import AdminNav from "./AdminNav";

export default function AdminLayout() {
  return (
    <div>
      <AdminNav />
      <Card bodyStyle={{ padding: 0 }} style={{ border: "none" }}>
        <Outlet />
      </Card>
    </div>
  );
}
