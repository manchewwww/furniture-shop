import { Button, Card, Form, Input, Row, Col, Typography, message } from "antd";
import { register } from "../api/auth";
import { useNavigate, useLocation } from "react-router-dom";
import { useI18n } from "../store/I18nContext";

export default function Register() {
  const nav = useNavigate();
  const location = useLocation() as any;
  const { t } = useI18n();
  const onFinish = async (v: any) => {
    try {
      const payload = {
        name: v.name,
        email: v.email,
        password: v.password,
        address: v.address,
        phone: v.phone,
      };
      await register(payload);
      message.success(t("register.success"));
      nav("/login", {
        replace: true,
        state: { from: location.state?.from, registered: true },
      });
    } catch {
      message.error(t("register.error"));
    }
  };
  return (
    <Row justify="center" align="middle" style={{ minHeight: "60vh" }}>
      <Col>
        <Card title={t("register.title")} style={{ maxWidth: 520 }}>
          <Form layout="vertical" onFinish={onFinish}>
            <Form.Item
              name="name"
              label={t("register.name")}
              rules={[
                { required: true, message: t("validation.required") },
                { min: 2, message: t("validation.min").replace("{n}", "2") },
              ]}
            >
              <Input />
            </Form.Item>
            <Form.Item
              name="email"
              label={t("register.email")}
              rules={[
                { required: true, message: t("validation.required") },
                { type: "email", message: t("validation.email") },
              ]}
            >
              <Input />
            </Form.Item>
            <Form.Item
              name="password"
              label={t("register.password")}
              rules={[
                { required: true, message: t("validation.required") },
                {
                  min: 8,
                  message: t("validation.password_min").replace("{n}", "8"),
                },
                {
                  pattern: /^(?=.*[A-Za-z])(?=.*\d).+$/,
                  message: t("validation.password_complexity"),
                },
              ]}
              hasFeedback
            >
              <Input.Password />
            </Form.Item>
            <Form.Item
              name="confirmPassword"
              label={t("register.confirm_password")}
              dependencies={["password"]}
              hasFeedback
              rules={[
                { required: true, message: t("validation.required") },
                ({ getFieldValue }) => ({
                  validator(_, value) {
                    if (!value || getFieldValue("password") === value) {
                      return Promise.resolve();
                    }
                    return Promise.reject(
                      new Error(t("validation.password_match"))
                    );
                  },
                }),
              ]}
            >
              <Input.Password />
            </Form.Item>
            <Form.Item
              name="address"
              label={t("register.address")}
              rules={[{ required: true, message: t("validation.required") }]}
            >
              <Input />
            </Form.Item>
            <Form.Item
              name="phone"
              label={t("register.phone")}
              rules={[
                { required: true, message: t("validation.required") },
                {
                  pattern: /^[0-9+()\-\s]{7,20}$/,
                  message: t("validation.phone"),
                },
              ]}
            >
              <Input />
            </Form.Item>
            <Button type="primary" htmlType="submit">
              {t("register.submit")}
            </Button>
          </Form>
        </Card>
      </Col>
    </Row>
  );
}
