import { Button, Card, Form, Input, Typography, message } from "antd";
import { register } from "../api/auth";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../store/AuthContext";

export default function Register() {
  const nav = useNavigate();
  const { refresh } = useAuth();
  const onFinish = async (v: any) => {
    try {
      await register(v);
      await refresh();
      message.success("Регистрацията е успешна");
      nav("/");
    } catch {
      message.error("Грешка при регистрация");
    }
  };
  return (
    <Card title="Регистрация" style={{ maxWidth: 520 }}>
      <Form layout="vertical" onFinish={onFinish}>
        <Form.Item name="name" label="Име" rules={[{ required: true }]}>
          <Input />
        </Form.Item>
        <Form.Item name="email" label="Имейл" rules={[{ required: true }]}>
          <Input />
        </Form.Item>
        <Form.Item name="password" label="Парола" rules={[{ required: true }]}>
          <Input.Password />
        </Form.Item>
        <Form.Item name="address" label="Адрес" rules={[{ required: true }]}>
          <Input />
        </Form.Item>
        <Form.Item name="phone" label="Телефон" rules={[{ required: true }]}>
          <Input />
        </Form.Item>
        <Button type="primary" htmlType="submit">
          Регистрация
        </Button>
      </Form>
    </Card>
  );
}
