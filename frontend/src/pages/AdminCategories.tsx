import {
  Button,
  Card,
  Form,
  Input,
  Modal,
  Popconfirm,
  Select,
  Table,
  message,
} from "antd";
import { useEffect, useState } from "react";
import { api } from "../api/client";
import { useI18n } from "../store/I18nContext";
import { useNavigate } from "react-router-dom";

export default function AdminCategories() {
  const { t } = useI18n();
  const nav = useNavigate();
  const [depts, setDepts] = useState<any[]>([]);
  const [categories, setCategories] = useState<any[]>([]);
  const [openCategory, setOpenCategory] = useState(false);
  const [editing, setEditing] = useState<any | null>(null);
  const [categoryForm] = Form.useForm();

  const load = async () => {
    try {
      const [d, c] = await Promise.all([
        api.get("/admin/departments"),
        api.get("/admin/categories"),
      ]);
      setDepts(d.data);
      setCategories(c.data);
    } catch {
      message.error("Failed to load categories");
    }
  };

  useEffect(() => {
    load();
  }, []);

  const submitCategory = async () => {
    const v = await categoryForm.validateFields();
    if (editing) {
      await api.put(`/admin/categories/${editing.id}`, v);
    } else {
      await api.post("/admin/categories", v);
    }
    setOpenCategory(false);
    setEditing(null);
    categoryForm.resetFields();
    load();
  };

  return (
    <div>
      <Card
        title="Categories"
        extra={
          <Button onClick={() => setOpenCategory(true)}>
            {t("create_category")}
          </Button>
        }
      >
        <Table
          rowKey="id"
          dataSource={categories}
          columns={[
            { title: "ID", dataIndex: "id", width: 60 },
            { title: t("category_name"), dataIndex: "name" },
            { title: t("category_description"), dataIndex: "description" },
            {
              title: t("department"),
              dataIndex: "department_id",
              render: (id: number) =>
                depts.find((d) => d.id === id)?.name || id,
            },
            {
              title: t("actions"),
              render: (_: any, r: any) => (
                <>
                  <Button
                    size="small"
                    onClick={() => {
                      setEditing(r);
                      categoryForm.setFieldsValue({
                        department_id: r.department_id,
                        name: r.name,
                        description: r.description,
                      });
                      setOpenCategory(true);
                    }}
                    style={{ marginRight: 8 }}
                  >
                    Edit
                  </Button>
                  <Popconfirm
                    title="Delete category?"
                    onConfirm={async () => {
                      await api.delete(`/admin/categories/${r.id}`);
                      load();
                    }}
                  >
                    <Button danger size="small">
                      Delete
                    </Button>
                  </Popconfirm>
                </>
              ),
            },
          ]}
        />
        <Modal
          title={editing ? t("edit_category") : t("create_category")}
          open={openCategory}
          onOk={submitCategory}
          onCancel={() => {
            setOpenCategory(false);
            setEditing(null);
            categoryForm.resetFields();
          }}
        >
          <Form layout="vertical" form={categoryForm}>
            <Form.Item
              name="department_id"
              label={t("department")}
              rules={[{ required: true }]}
            >
              <Select placeholder={t("department")}>
                {depts.map((d: any) => (
                  <Select.Option key={d.id} value={d.id}>
                    {d.name}
                  </Select.Option>
                ))}
              </Select>
            </Form.Item>
            <Form.Item
              name="name"
              label={t("category_name")}
              rules={[{ required: true }]}
            >
              <Input />
            </Form.Item>
            <Form.Item name="description" label={t("category_description")}>
              <Input />
            </Form.Item>
          </Form>
        </Modal>
      </Card>
    </div>
  );
}
