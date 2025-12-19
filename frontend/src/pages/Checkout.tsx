import {
  Alert,
  Button,
  Card,
  Form,
  Input,
  Select,
  Typography,
  message,
} from "antd";
import { useCart } from "../store/CartContext";
import { createOrder, payByCard } from "../api/orders";
import { useState } from "react";
import { useI18n } from "../store/I18nContext";

export default function Checkout() {
  const { items, clear } = useCart();
  const [orderId, setOrderId] = useState<number | null>(null);
  const [paymentMethod, setPaymentMethod] = useState<string>("card");
  const { t } = useI18n();

  const onFinish = async (values: any) => {
    const payload = {
      name: values.name,
      email: values.email,
      phone: values.phone,
      address: values.address,
      payment_method: values.payment_method,
      items: items.map((it) => ({
        product_id: it.product.id,
        quantity: it.quantity,
        options: it.options,
      })),
    };
    try {
      const res = await createOrder(payload);
      setOrderId(res.order_id);
      setPaymentMethod(values.payment_method);
      message.success(t("checkout.success"));
    } catch {
      message.error(t("checkout.error"));
    }
  };

  const onPayCard = async (values: any) => {
    try {
      await payByCard({ order_id: orderId, ...values });
      message.success(t("checkout.pay.success"));
      clear();
    } catch {
      message.error(t("checkout.pay.error"));
    }
  };

  return (
    <div>
      <Typography.Title level={2}>{t("checkout.title")}</Typography.Title>
      {!orderId && (
        <Card title={t("checkout.form.title")}>
          <Form
            layout="vertical"
            onFinish={onFinish}
            initialValues={{ payment_method: "card" }}
          >
            <Form.Item
              name="name"
              label={t("checkout.name")}
              rules={[{ required: true }]}
            >
              <Input />
            </Form.Item>
            <Form.Item
              name="email"
              label={t("checkout.email")}
              rules={[{ required: true }]}
            >
              <Input />
            </Form.Item>
            <Form.Item
              name="phone"
              label={t("checkout.phone")}
              rules={[{ required: true }]}
            >
              <Input />
            </Form.Item>
            <Form.Item
              name="address"
              label={t("checkout.address")}
              rules={[{ required: true }]}
            >
              <Input />
            </Form.Item>
            <Form.Item
              name="payment_method"
              label={t("checkout.payment_method")}
              rules={[{ required: true }]}
            >
              <Select
                options={[
                  { value: "card", label: "Card" },
                  { value: "cod", label: "Cash on Delivery" },
                  { value: "bank", label: "Bank Transfer" },
                ]}
              />
            </Form.Item>
            <Button type="primary" htmlType="submit">
              {t("checkout.place_order")}
            </Button>
          </Form>
        </Card>
      )}
      {orderId && paymentMethod !== "card" && (
        <Alert type="success" message={t("checkout.success")} />
      )}
      {orderId && paymentMethod === "card" && (
        <Card title={t("checkout.card.title")}>
          <Form layout="vertical" onFinish={onPayCard}>
            <Form.Item
              name="cardholder_name"
              label={t("checkout.cardholder_name")}
              rules={[{ required: true }]}
            >
              <Input />
            </Form.Item>
            <Form.Item
              name="card_number"
              label={t("checkout.card_number")}
              rules={[{ required: true }]}
            >
              <Input />
            </Form.Item>
            <Form.Item
              name="expiry_month"
              label={t("checkout.exp_month")}
              rules={[{ required: true }]}
            >
              <Input />
            </Form.Item>
            <Form.Item
              name="expiry_year"
              label={t("checkout.exp_year")}
              rules={[{ required: true }]}
            >
              <Input />
            </Form.Item>
            <Form.Item name="cvv" label="CVV" rules={[{ required: true }]}>
              <Input />
            </Form.Item>
            <Button type="primary" htmlType="submit">
              {t("checkout.pay")}
            </Button>
          </Form>
        </Card>
      )}
    </div>
  );
}
