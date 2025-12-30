import React, { useState } from "react";
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
import { createOrder } from "../api/orders";
import { useI18n } from "../store/I18nContext";

type PaymentMethod = "card";

export default function Checkout() {
  const { items, clear } = useCart();
  const { t } = useI18n();
  const [orderId, setOrderId] = useState<number | null>(null);
  const [paymentMethod, setPaymentMethod] = useState<PaymentMethod>("card");
  const [placing, setPlacing] = useState(false);
  const [orderForm] = Form.useForm();

  const requiredTrim = (msg: string) => ({
    required: true,
    validator: async (_: any, v: any) => {
      const s = String(v ?? "").trim();
      if (!s) throw new Error(msg);
    },
  });

  const onPlaceOrder = async (values: any) => {
    if (!items.length) {
      message.error(t("checkout.empty_cart") || "Your cart is empty.");
      return;
    }
    const payload = {
      name: String(values.name).trim(),
      email: String(values.email).trim(),
      phone: String(values.phone).trim(),
      address: String(values.address).trim(),
      payment_method: values.payment_method as PaymentMethod,
      items: items.map((it) => ({
        product_id: it.product.id,
        quantity: it.quantity,
        options: it.options,
      })),
    };
    setPlacing(true);
    try {
      const res = await createOrder(payload);
      if (payload.payment_method === "card" && res.checkout_url) {
        window.location.assign(res.checkout_url as string);
        return;
      }
      setOrderId(res.order_id);
      setPaymentMethod(payload.payment_method);
    } catch {
      message.error(t("checkout.error") || "Failed to create order");
    } finally {
      setPlacing(false);
    }
  };

  return (
    <div>
      <Typography.Title level={2}>{t("checkout.title")}</Typography.Title>

      {!orderId && (
        <Card title={t("checkout.form.title")}>
          <Form
            form={orderForm}
            layout="vertical"
            onFinish={onPlaceOrder}
            initialValues={{ payment_method: "card" }}
          >
            <Form.Item
              name="name"
              label={t("checkout.name")}
              rules={[
                requiredTrim("Name is required"),
                { min: 2, message: "Name must be at least 2 characters" },
              ]}
            >
              <Input autoComplete="name" />
            </Form.Item>
            <Form.Item
              name="email"
              label={t("checkout.email")}
              rules={[
                { required: true, message: "Email is required" },
                { type: "email", message: "Enter a valid email" },
              ]}
            >
              <Input autoComplete="email" />
            </Form.Item>
            <Form.Item
              name="phone"
              label={t("checkout.phone")}
              rules={[
                { required: true, message: "Phone is required" },
                {
                  pattern: /^[0-9+\-()\s]{7,20}$/,
                  message: "Enter a valid phone number",
                },
              ]}
            >
              <Input autoComplete="tel" />
            </Form.Item>
            <Form.Item
              name="address"
              label={t("checkout.address")}
              rules={[
                requiredTrim("Address is required"),
                { min: 5, message: "Address must be at least 5 characters" },
              ]}
            >
              <Input autoComplete="street-address" />
            </Form.Item>
            <Form.Item
              name="payment_method"
              label={t("checkout.payment_method")}
              rules={[{ required: true, message: "Choose a payment method" }]}
            >
              <Select
                options={[
                  {
                    value: "card",
                    label: t("checkout.payment.card") || "Card",
                  },
                ]}
              />
            </Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              loading={placing}
              disabled={!items.length}
            >
              {t("checkout.place_order")}
            </Button>
          </Form>
        </Card>
      )}
    </div>
  );
}
