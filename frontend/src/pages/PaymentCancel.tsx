import { Alert, Card } from "antd";
import { useEffect } from "react";
import { useNavigate } from "react-router-dom";

export default function PaymentCancel() {
  const navigate = useNavigate();
  useEffect(() => {
    navigate("/orders", { replace: true });
  }, []);
  return (
    <Card>
      <Alert
        type="warning"
        showIcon
        message="Payment canceled. Redirecting to My Orders..."
      />
    </Card>
  );
}
