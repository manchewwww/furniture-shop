import React, { useMemo, useState } from "react";
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

function onlyDigits(value: unknown) {
  return String(value ?? "").replace(/\s+/g, "");
}

export default function Checkout() {
  const { items, clear } = useCart();
  const { t } = useI18n();

  const [orderId, setOrderId] = useState<number | null>(null);
  const [paymentMethod, setPaymentMethod] = useState<PaymentMethod>("card");
  const [placing, setPlacing] = useState(false);
  const [paying, setPaying] = useState(false);

  const [orderForm] = Form.useForm();
  const [cardForm] = Form.useForm();

  const rules = useMemo(() => {
    const requiredTrim = (msg: string) => ({
      required: true,
      validator: async (_: any, v: any) => {
        const s = String(v ?? "").trim();
        if (!s) throw new Error(msg);
      },
    });

    const minTrim = (min: number, msg: string) => ({
      validator: async (_: any, v: any) => {
        const s = String(v ?? "").trim();
        if (s.length < min) throw new Error(msg);
      },
    });

    const emailRule = () => ({
      type: "email" as const,
      message: "Enter a valid email",
    });

    const phoneRule = () => ({
      pattern: /^[0-9+\-()\s]{7,20}$/,
      message: "Enter a valid phone number",
    });

    const cardNumberRule = () => ({
      validator: async (_: any, v: any) => {
        const digits = onlyDigits(v);
        if (!/^\d{13,19}$/.test(digits)) throw new Error("Invalid card number");
        if (!luhnValid(digits)) throw new Error("Invalid card number");
      },
    });

    const monthRule = () => ({
      validator: async (_: any, v: any) => {
        const m = Number(v);
        if (!Number.isInteger(m) || m < 1 || m > 12)
          throw new Error("Invalid month");
      },
    });

    const yearRule = () => ({
      validator: async (_: any, v: any) => {
        const now = new Date();
        const raw = String(v ?? "").trim();
        let y = Number(raw);
        if (!Number.isInteger(y)) throw new Error("Invalid year");
        if (raw.length === 2) y += 2000;

        const m = Number(cardForm.getFieldValue("expiry_month")) || 1; // 1..12
        const expEnd = new Date(y, m, 0); // last day of month
        if (expEnd < now) throw new Error("Card expired");
      },
    });

    const cvvRule = () => ({ pattern: /^\d{3,4}$/, message: "3-4 digits" });

    return {
      requiredTrim,
      minTrim,
      emailRule,
      phoneRule,
      cardNumberRule,
      monthRule,
      yearRule,
      cvvRule,
    };
  }, [cardForm]);

  const onFinishOrder = async (values: any) => {
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
      message.success(t("checkout.success") || "Order created successfully");

      if (payload.payment_method !== "card") {
        clear();
        orderForm.resetFields();
      }
    } catch {
      message.error(t("checkout.error") || "Failed to create order");
    } finally {
      setPlacing(false);
    }
  };

  const onFinishCard = async (values: any) => {
    if (!orderId) {
      message.error(t("checkout.pay.no_order") || "No order to pay for");
      return;
    }

    setPaying(true);
    try {
      await payByCard({
        order_id: orderId,
        cardholder_name: String(values.cardholder_name).trim(),
        card_number: onlyDigits(values.card_number),
        expiry_month: String(values.expiry_month).trim(),
        expiry_year: String(values.expiry_year).trim(),
        cvv: String(values.cvv).trim(),
      });

      message.success(t("checkout.pay.success") || "Payment successful");
      clear();
      orderForm.resetFields();
      cardForm.resetFields();
      setOrderId(null);
      setPaymentMethod("card");
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
            initialValues={{ payment_method: "card" as PaymentMethod }}
            onFinish={onFinishOrder}
            validateTrigger={["onBlur", "onSubmit"]}
          >
            <Form.Item
              name="name"
              label={t("checkout.name")}
              rules={[
                rules.requiredTrim("Name is required"),
                rules.minTrim(2, "Name must be at least 2 characters"),
              ]}
            >
              <Input autoComplete="name" />
            </Form.Item>

            <Form.Item
              name="email"
              label={t("checkout.email")}
              rules={[
                rules.requiredTrim("Email is required"),
                rules.emailRule(),
              ]}
            >
              <Input autoComplete="email" />
            </Form.Item>

            <Form.Item
              name="phone"
              label={t("checkout.phone")}
              rules={[
                rules.requiredTrim("Phone is required"),
                rules.phoneRule(),
              ]}
            >
              <Input autoComplete="tel" />
            </Form.Item>

            <Form.Item
              name="address"
              label={t("checkout.address")}
              rules={[
                rules.requiredTrim("Address is required"),
                rules.minTrim(5, "Address must be at least 5 characters"),
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

      {orderId && paymentMethod !== "card" && (
        <Alert
          type="success"
          showIcon
          message={t("checkout.success") || "Order created successfully."}
        />
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

          <Form
            form={cardForm}
            layout="vertical"
            onFinish={onFinishCard}
            validateTrigger={["onBlur", "onSubmit"]}
          >
            <Form.Item
              name="cardholder_name"
              label={t("checkout.cardholder_name")}
              rules={[
                rules.requiredTrim("Cardholder is required"),
                rules.minTrim(2, "Cardholder must be at least 2 characters"),
              ]}
            >
              <Input autoComplete="cc-name" />
            </Form.Item>

            <Form.Item
              name="card_number"
              label={t("checkout.card_number")}
              rules={[
                rules.requiredTrim("Card number is required"),
                rules.cardNumberRule(),
              ]}
            >
              <Input autoComplete="cc-number" inputMode="numeric" />
            </Form.Item>

            <Form.Item
              name="expiry_month"
              label={t("checkout.exp_month")}
              rules={[rules.requiredTrim("Month required"), rules.monthRule()]}
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
              rules={[rules.requiredTrim("Year required"), rules.yearRule()]}
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
              rules={[rules.requiredTrim("CVV required"), rules.cvvRule()]}
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
