import { Alert, Card } from "antd";

export default function PaymentCancel() {
  return (
    <Card>
      <Alert type="warning" showIcon message="Payment canceled." />
    </Card>
  );
}
