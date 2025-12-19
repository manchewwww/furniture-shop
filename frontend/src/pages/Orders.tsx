import { Table, Tag, Typography, message } from "antd";
import { useEffect, useState } from "react";
import { myOrders } from "../api/orders";
import { useI18n } from "../store/I18nContext";

export default function Orders() {
  const [orders, setOrders] = useState<any[]>([]);
  const { t } = useI18n();
  useEffect(() => {
    myOrders()
      .then(setOrders)
      .catch(() => message.error(t("orders.error")));
  }, []);
  return (
    <div>
      <Typography.Title level={2}>{t("orders.title")}</Typography.Title>
      <Table
        rowKey="id"
        dataSource={orders}
        columns={[
          { title: t("orders.col.id"), dataIndex: "id" },
          {
            title: t("orders.col.status"),
            dataIndex: "status",
            render: (s) => <Tag>{s}</Tag>,
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
        ]}
      />
    </div>
  );
}
