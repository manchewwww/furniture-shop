import { Button, Card, Form, Input, Typography, message } from "antd";
import { login } from "../api/auth";
import { useNavigate, Link } from "react-router-dom";
import { useAuth } from "../store/AuthContext";
import { useI18n } from "../store/I18nContext";

export default function Login() {
  const nav = useNavigate();
  const { refresh } = useAuth();
  const { t } = useI18n();
  const onFinish = async (v: any) => {
    try {
      await login(v.email, v.password);
      await refresh();
      message.success(t("login.success"));
      nav("/");
    } catch {
      message.error(t("login.error"));
    }
  };
  return (
    <Card title={t("login.title")} style={{ maxWidth: 420 }}>
      <Form layout="vertical" onFinish={onFinish}>
        <Form.Item
          name="email"
          label={t("login.email")}
          rules={[{ required: true }]}
        >
          <Input />
        </Form.Item>
        <Form.Item
          name="password"
          label={t("login.password")}
          rules={[{ required: true }]}
        >
          <Input.Password />
        </Form.Item>
        <Button type="primary" htmlType="submit">
          {t("login.submit")}
        </Button>
      </Form>
      <Typography.Paragraph style={{ marginTop: 12 }}>
        {t("login.register_cta")}{" "}
        <Link to="/register">{t("nav.register")}</Link>
      </Typography.Paragraph>
    </Card>
  );
}
