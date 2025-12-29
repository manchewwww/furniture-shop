import { Button, Card, Select, Table, Tag, message } from "antd";
import { useEffect, useState } from "react";
import { api } from "../api/client";
import { useI18n } from "../store/I18nContext";
import { useNavigate } from "react-router-dom";

const ORDER_STATUSES = [
  "new",
  "processing",
  "in_production",
  "shipped",
  "delivered",
  "cancelled",
];

export default function AdminOrders() {
  const { t } = useI18n();
  const nav = useNavigate();
  const [orders, setOrders] = useState<any[]>([]);
  const [status, setStatus] = useState<string | undefined>();

  const load = async () => {
    try {
      const res = await api.get("/admin/orders", { params: { status } });
      setOrders(res.data);
    } catch {
      message.error("Failed to load orders");
    }
  };

  useEffect(() => {
    load();
  }, [status]);

  const updateStatus = async (id: number, s: string) => {
    try {
      await api.patch(`/admin/orders/${id}/status`, { status: s });
      load();
    } catch {
      message.error("Failed to update status");
    }
  };

  return (
    <div>
      <Card style={{ marginBottom: 16 }}>
        <Button type="link" onClick={() => nav("/admin/departments")}>
          Departments
        </Button>
        <Button type="link" onClick={() => nav("/admin/categories")}>
          Categories
        </Button>
        <Button type="link" onClick={() => nav("/admin/products")}>
          Products
        </Button>
        <Button type="link" onClick={() => nav("/admin/orders")}>
          Orders
        </Button>
      </Card>
      <Card
        title={t("orders_title")}
        extra={
          <Select
            allowClear
            placeholder="Filter by status"
            style={{ minWidth: 180 }}
            value={status as any}
            onChange={(v) => setStatus(v as string | undefined)}
            options={ORDER_STATUSES.map((s) => ({ value: s, label: s }))}
          />
        }
      >
        <Table
          rowKey="id"
          dataSource={orders}
          columns={[
            { title: t("orders.col.id"), dataIndex: "id" },
            {
              title: t("orders.col.status"),
              dataIndex: "status",
              render: (s: string) => <Tag>{s}</Tag>,
            },
            {
              title: t("orders.col.payment_status"),
              dataIndex: "payment_status",
            },
            { title: t("orders.col.total"), dataIndex: "total_price" },
            {
              title: t("orders.col.eta_days"),
              dataIndex: "estimated_production_time_days",
            },
            {
              title: t("actions"),
              render: (_: any, r: any) => (
                <Select
                  placeholder="Set status"
                  style={{ minWidth: 160 }}
                  onChange={(v) => updateStatus(r.id, v as string)}
                  options={ORDER_STATUSES.map((s) => ({ value: s, label: s }))}
                />
              ),
            },
          ]}
        />
      </Card>
    </div>
  );
}
