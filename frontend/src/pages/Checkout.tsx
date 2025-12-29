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
import { createOrder, payByCard } from "../api/orders";
import { useI18n } from "../store/I18nContext";

type PaymentMethod = "card" | "bank";

function luhnValid(digits: string) {
  let sum = 0;
  let dbl = false;
  for (let i = digits.length - 1; i >= 0; i--) {
    let d = Number(digits[i]);
    if (Number.isNaN(d)) return false;
    if (dbl) {
      d *= 2;
      if (d > 9) d -= 9;
    }
    sum += d;
    dbl = !dbl;
  }
  return sum % 10 === 0;
}

export default function Checkout() {
  const { items, clear } = useCart();
  const { t } = useI18n();
  const [orderId, setOrderId] = useState<number | null>(null);
  const [paymentMethod, setPaymentMethod] = useState<PaymentMethod>("card");
  const [placing, setPlacing] = useState(false);
  const [paying, setPaying] = useState(false);
  const [bankInfo, setBankInfo] = useState<any | null>(null);
  const [orderForm] = Form.useForm();
  const [cardForm] = Form.useForm();

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
      setOrderId(res.order_id);
      setPaymentMethod(payload.payment_method);
      if (payload.payment_method === "card" && res.checkout_url) {
        window.location.href = res.checkout_url as string;
        return;
      }
      if (payload.payment_method === "bank" && res.instructions) {
        setBankInfo(res.instructions);
      }
    } catch {
      message.error(t("checkout.error") || "Failed to create order");
    } finally {
      setPlacing(false);
    }
  };

  const onPayCard = async (values: any) => {
    if (!orderId) {
      message.error(t("checkout.pay.no_order") || "No order to pay for");
      return;
    }
    setPaying(true);
    try {
      await payByCard({ order_id: orderId, ...values });
      message.success(t("checkout.pay.success") || "Payment successful");
      clear();
      cardForm.resetFields();
      setOrderId(null);
    } catch {
      message.error(t("checkout.pay.error") || "Payment failed");
    } finally {
      setPaying(false);
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
                  {
                    value: "bank",
                    label: t("checkout.payment.bank") || "Bank Transfer",
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

      {orderId && paymentMethod === "bank" && bankInfo && (
        <Card title={t("checkout.payment.bank") || "Bank Transfer"}>
          <Alert
            type="success"
            showIcon
            style={{ marginBottom: 12 }}
            message={t("checkout.success") || "Order created successfully."}
          />
          <Typography.Paragraph>
            <b>{t("orders.col.total") || "Total"}:</b> {bankInfo.amount}{" "}
            {bankInfo.currency}
          </Typography.Paragraph>
          <Typography.Paragraph>
            <b>Beneficiary:</b> {bankInfo.beneficiary_name}
            <br />
            <b>Bank:</b> {bankInfo.bank_name}
            <br />
            <b>IBAN:</b> {bankInfo.iban}
            <br />
            <b>BIC:</b> {bankInfo.bic}
          </Typography.Paragraph>
          <Alert
            type="info"
            showIcon
            message={
              (t("checkout.noncard.instructions") ||
                "Use this reference in the transfer:") +
              ` ${bankInfo.payment_reference}`
            }
          />
        </Card>
      )}

      {orderId && paymentMethod === "card" && (
        <Card title={t("checkout.card.title")}>
          <Alert
            type="info"
            showIcon
            style={{ marginBottom: 16 }}
            message={
              t("checkout.card.order_created") ||
              "Order created. Please enter your card details to pay."
            }
          />
          <Form form={cardForm} layout="vertical" onFinish={onPayCard}>
            <Form.Item
              name="cardholder_name"
              label={t("checkout.cardholder_name")}
              rules={[
                requiredTrim("Cardholder is required"),
                { min: 2, message: "Cardholder must be at least 2 characters" },
              ]}
            >
              <Input autoComplete="cc-name" />
            </Form.Item>
            <Form.Item
              name="card_number"
              label={t("checkout.card_number")}
              rules={[
                { required: true, message: "Card number is required" },
                {
                  validator: async (_: any, v: string) => {
                    const d = String(v || "").replace(/\s+/g, "");
                    if (!/^\d{13,19}$/.test(d) || !luhnValid(d))
                      throw new Error("Invalid card number");
                  },
                },
              ]}
            >
              <Input autoComplete="cc-number" inputMode="numeric" />
            </Form.Item>
            <Form.Item
              name="expiry_month"
              label={t("checkout.exp_month")}
              rules={[
                { required: true, message: "Month required" },
                {
                  validator: async (_: any, v: string) => {
                    const m = Number(v);
                    if (!Number.isInteger(m) || m < 1 || m > 12)
                      throw new Error("Invalid month");
                  },
                },
              ]}
            >
              <Input
                autoComplete="cc-exp-month"
                inputMode="numeric"
                placeholder="MM"
              />
            </Form.Item>
            <Form.Item
              name="expiry_year"
              label={t("checkout.exp_year")}
              dependencies={["expiry_month"]}
              rules={[
                { required: true, message: "Year required" },
                {
                  validator: async ({ getFieldValue }: any, v: string) => {
                    const now = new Date();
                    let y = Number(v);
                    if (!Number.isInteger(y)) throw new Error("Invalid year");
                    if ((v || "").length === 2) y += 2000;
                    const m = Number(getFieldValue("expiry_month")) || 1;
                    if (new Date(y, m, 0) < now)
                      throw new Error("Card expired");
                  },
                },
              ]}
            >
              <Input
                autoComplete="cc-exp-year"
                inputMode="numeric"
                placeholder="YY or YYYY"
              />
            </Form.Item>
            <Form.Item
              name="cvv"
              label="CVV"
              rules={[
                { required: true, message: "CVV required" },
                { pattern: /^\d{3,4}$/, message: "3-4 digits" },
              ]}
            >
              <Input.Password autoComplete="cc-csc" inputMode="numeric" />
            </Form.Item>
            <Button type="primary" htmlType="submit" loading={paying}>
              {t("checkout.pay") || "Pay"}
            </Button>
          </Form>
        </Card>
      )}
    </div>
  );
}
