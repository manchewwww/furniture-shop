import { Spin } from "antd";
import { useEffect, useMemo } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { useCart } from "../store/CartContext";

export default function PaymentSuccess() {
  const [sp] = useSearchParams();
  const orderIdParam = sp.get("order_id");
  const orderId = useMemo(() => Number(orderIdParam || 0) || 0, [orderIdParam]);
  const { clear } = useCart();
  const nav = useNavigate();

  useEffect(() => {
    if (!orderId) return;
    clear();
    nav(`/orders?open=${orderId}`, { replace: true });
  }, [orderId]);

  return (
    <div style={{ display: "flex", justifyContent: "center", padding: 24 }}>
      <Spin />
    </div>
  );
}
