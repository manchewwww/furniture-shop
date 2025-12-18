import { Button, Card, Form, Input, Typography, message } from "antd";
import { login } from "../api/auth";
import { useNavigate, Link } from "react-router-dom";
import { useAuth } from "../store/AuthContext";

export default function Login() {
  const nav = useNavigate();
  const { refresh } = useAuth();
  const onFinish = async (v: any) => {
    try {
      await login(v.email, v.password);
      await refresh();
      message.success("Вход успешен");
      nav("/");
    } catch {
      message.error("Грешен имейл или парола");
    }
  };
  return (
    <Card title="Вход" style={{ maxWidth: 420 }}>
      <Form layout="vertical" onFinish={onFinish}>
        <Form.Item name="email" label="Имейл" rules={[{ required: true }]}>
          <Input />
        </Form.Item>
        <Form.Item name="password" label="Парола" rules={[{ required: true }]}>
          <Input.Password />
        </Form.Item>
        <Button type="primary" htmlType="submit">
          Вход
        </Button>
      </Form>
      <Typography.Paragraph style={{ marginTop: 12 }}>
        Don’t have an account? <Link to="/register">Register</Link>
      </Typography.Paragraph>
    </Card>
  );
}
