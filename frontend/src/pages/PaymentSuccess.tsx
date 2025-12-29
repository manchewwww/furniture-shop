import { Alert, Card, Descriptions, Table, Typography, message } from "antd";
import { useEffect, useMemo, useState } from "react";
import { useSearchParams } from "react-router-dom";
import { myOrder } from "../api/orders";
import { useCart } from "../store/CartContext";

export default function PaymentSuccess() {
  const [sp] = useSearchParams();
  const orderIdParam = sp.get("order_id");
  const orderId = useMemo(() => Number(orderIdParam || 0) || 0, [orderIdParam]);
  const [order, setOrder] = useState<any | null>(null);
  const { clear } = useCart();

  useEffect(() => {
    if (!orderId) return;
    myOrder(orderId)
      .then((o) => {
        setOrder(o);
        clear();
      })
      .catch(() => message.error("Unable to load order details"));
  }, [orderId]);

  const readyDate = useMemo(() => {
    if (!order) return null;
    const created = new Date(order.created_at);
    const d = new Date(created);
    d.setDate(d.getDate() + (order.estimated_production_time_days || 0));
    return d;
  }, [order]);

  return (
    <Card>
      <Alert type="success" showIcon message="Payment completed." />
      {order && (
        <>
          <Descriptions bordered style={{ marginTop: 16 }} size="small">
            <Descriptions.Item label="Order ID">{order.id}</Descriptions.Item>
            <Descriptions.Item label="Payment Status">
              {order.payment_status}
            </Descriptions.Item>
            <Descriptions.Item label="Total Paid">
              {order.total_price} EUR
            </Descriptions.Item>
            <Descriptions.Item label="Ordered On">
              {new Date(order.created_at).toLocaleString()}
            </Descriptions.Item>
            <Descriptions.Item label="ETA (days)">
              {order.estimated_production_time_days}
            </Descriptions.Item>
            <Descriptions.Item label="Estimated Ready By">
              {readyDate ? readyDate.toLocaleDateString() : "-"}
            </Descriptions.Item>
          </Descriptions>

          <Typography.Title level={5} style={{ marginTop: 16 }}>
            Items
          </Typography.Title>
          <Table
            rowKey="id"
            size="small"
            pagination={false}
            dataSource={order.items || []}
            columns={[
              { title: "Product", dataIndex: "product_id" },
              { title: "Qty", dataIndex: "quantity" },
              { title: "Unit Price", dataIndex: "unit_price" },
              { title: "Line Total", dataIndex: "line_total" },
            ]}
          />
        </>
      )}
    </Card>
  );
}
