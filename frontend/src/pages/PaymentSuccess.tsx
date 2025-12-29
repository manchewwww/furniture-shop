import { Alert, Card, Typography } from "antd";
import { useSearchParams } from "react-router-dom";

export default function PaymentSuccess() {
  const [sp] = useSearchParams();
  const orderId = sp.get("order_id");
  return (
    <Card>
      <Alert type="success" showIcon message="Payment completed (demo)." />
      <Typography.Paragraph style={{ marginTop: 12 }}>
        Order ID: {orderId}
      </Typography.Paragraph>
    </Card>
  );
}
