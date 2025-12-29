import { Button, Card, Form, Input, Typography, message } from "antd";
import { register } from "../api/auth";
import { useNavigate, useLocation } from "react-router-dom";
import { useAuth } from "../store/AuthContext";
import { useI18n } from "../store/I18nContext";

export default function Register() {
  const nav = useNavigate();
  const location = useLocation() as any;
  const { refresh } = useAuth();
  const { t } = useI18n();
  const onFinish = async (v: any) => {
    try {
      await register(v);
      await refresh();
      message.success(t("register.success"));
      const to = location.state?.from?.pathname || "/";
      nav(to, { replace: true });
    } catch {
      message.error(t("register.error"));
    }
  };
  return (
    <Card title={t("register.title")} style={{ maxWidth: 520 }}>
      <Form layout="vertical" onFinish={onFinish}>
        <Form.Item
          name="name"
          label={t("register.name")}
          rules={[{ required: true }]}
        >
          <Input />
        </Form.Item>
        <Form.Item
          name="email"
          label={t("register.email")}
          rules={[{ required: true }]}
        >
          <Input />
        </Form.Item>
        <Form.Item
          name="password"
          label={t("register.password")}
          rules={[{ required: true }]}
        >
          <Input.Password />
        </Form.Item>
        <Form.Item
          name="address"
          label={t("register.address")}
          rules={[{ required: true }]}
        >
          <Input />
        </Form.Item>
        <Form.Item
          name="phone"
          label={t("register.phone")}
          rules={[{ required: true }]}
        >
          <Input />
        </Form.Item>
        <Button type="primary" htmlType="submit">
          {t("register.submit")}
        </Button>
      </Form>
    </Card>
  );
}
