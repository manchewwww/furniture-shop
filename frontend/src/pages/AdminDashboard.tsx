import {
  Button,
  Card,
  Form,
  Input,
  InputNumber,
  Modal,
  Select,
  Table,
  Tabs,
  Tag,
  Upload,
  message,
  Popconfirm,
} from "antd";
import { UploadOutlined } from "@ant-design/icons";
import { useEffect, useState } from "react";
import { api } from "../api/client";

export default function AdminDashboard() {
  const [depts, setDepts] = useState<any[]>([]);
  const [orders, setOrders] = useState<any[]>([]);
  const [products, setProducts] = useState<any[]>([]);
  const [categories, setCategories] = useState<any[]>([]);
  const [open, setOpen] = useState(false);
  const [form] = Form.useForm();
  const [openProduct, setOpenProduct] = useState(false);
  const [productForm] = Form.useForm();
  const [editing, setEditing] = useState<any | null>(null);
  const load = async () => {
    try {
      const [d, o, p, c] = await Promise.all([
        api.get("/admin/departments"),
        api.get("/admin/orders"),
        api.get("/admin/products"),
        api.get("/admin/categories"),
      ]);
      setDepts(d.data);
      setOrders(o.data);
      setProducts(p.data);
      setCategories(c.data);
    } catch {
      message.error("Нужен е админ вход");
    }
  };
  useEffect(() => {
    load();
  }, []);

  const createDept = async () => {
    const v = await form.validateFields();
    await api.post("/admin/departments", v);
    setOpen(false);
    form.resetFields();
    load();
  };

  const setStatus = async (id: number, status: string) => {
    await api.patch(`/admin/orders/${id}/status`, { status });
    load();
  };

  const submitProduct = async () => {
    const v = await productForm.validateFields();
    const payload = {
      category_id: v.category_id,
      name: v.name,
      short_description: v.description,
      long_description: v.description,
      base_price: v.price,
      image_url: v.image,
    };
    if (editing) {
      await api.put(`/admin/products/${editing.id}`, payload);
    } else {
      await api.post("/admin/products", payload);
    }
    setOpenProduct(false);
    setEditing(null);
    productForm.resetFields();
    load();
  };

  const removeProduct = async (id: number) => {
    await api.delete(`/admin/products/${id}`);
    load();
  };

  return (
    <div>
      <Tabs
        items={[
          {
            key: "depts",
            label: "Отдели",
            children: (
              <Card
                title="Отдели"
                extra={<Button onClick={() => setOpen(true)}>Нов отдел</Button>}
              >
                <Table
                  rowKey="id"
                  dataSource={depts}
                  columns={[
                    { title: "Име", dataIndex: "name" },
                    { title: "Описание", dataIndex: "description" },
                  ]}
                />
                <Modal
                  title="Нов отдел"
                  open={open}
                  onOk={createDept}
                  onCancel={() => setOpen(false)}
                >
                  <Form layout="vertical" form={form}>
                    <Form.Item
                      name="name"
                      label="Име"
                      rules={[{ required: true }]}
                    >
                      <Input />
                    </Form.Item>
                    <Form.Item name="description" label="Описание">
                      <Input />
                    </Form.Item>
                  </Form>
                </Modal>
              </Card>
            ),
          },
          {
            key: "orders",
            label: "Поръчки",
            children: (
              <Card title="Всички поръчки">
                <Table
                  rowKey="id"
                  dataSource={orders}
                  columns={[
                    { title: "№", dataIndex: "id" },
                    {
                      title: "Статус",
                      dataIndex: "status",
                      render: (s, r) => <Tag>{s}</Tag>,
                    },
                    { title: "Плащане", dataIndex: "payment_status" },
                    { title: "Общо", dataIndex: "total_price" },
                    {
                      title: "Действия",
                      render: (_: any, r: any) => (
                        <>
                          {[
                            "нова",
                            "потвърдена",
                            "впроизводство",
                            "изпратена",
                            "доставена",
                            "отказана",
                          ].map((st) => (
                            <Button
                              key={st}
                              size="small"
                              style={{ marginRight: 4 }}
                              onClick={() => setStatus(r.id, st)}
                            >
                              {st}
                            </Button>
                          ))}
                        </>
                      ),
                    },
                  ]}
                />
              </Card>
            ),
          },
        ]}
      />

      <Card
        title="????????"
        style={{ marginTop: 16 }}
        extra={
          <Button
            onClick={() => {
              setEditing(null);
              productForm.resetFields();
              setOpenProduct(true);
            }}
          >
            ????? ?????
          </Button>
        }
      >
        <Table
          rowKey="id"
          dataSource={products}
          columns={[
            { title: "?", dataIndex: "id", width: 60 },
            { title: "????", dataIndex: "name" },
            { title: "????????", dataIndex: "short_description" },
            { title: "????", dataIndex: "base_price" },
            {
              title: "????????",
              dataIndex: "image_url",
              render: (u: string) =>
                u ? (
                  <img
                    src={u}
                    alt=""
                    style={{ width: 60, height: 40, objectFit: "cover" }}
                  />
                ) : null,
            },
            {
              title: "",
              render: (_: any, r: any) => (
                <>
                  <Button
                    size="small"
                    style={{ marginRight: 8 }}
                    onClick={() => {
                      setEditing(r);
                      setOpenProduct(true);
                      productForm.setFieldsValue({
                        category_id: r.category_id,
                        name: r.name,
                        description: r.short_description || r.long_description,
                        price: r.base_price,
                        image: r.image_url,
                      });
                    }}
                  >
                    Edit
                  </Button>
                  <Popconfirm
                    title="Delete product?"
                    onConfirm={() => removeProduct(r.id)}
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
          title={editing ? "???????? ?????" : "????? ?????"}
          open={openProduct}
          onOk={submitProduct}
          onCancel={() => {
            setOpenProduct(false);
            setEditing(null);
          }}
        >
          <Form layout="vertical" form={productForm}>
            <Form.Item
              name="category_id"
              label="?????????"
              rules={[{ required: true }]}
            >
              <Select placeholder="????? ?????????">
                {categories.map((c: any) => (
                  <Select.Option key={c.id} value={c.id}>
                    {c.name}
                  </Select.Option>
                ))}
              </Select>
            </Form.Item>
            <Form.Item name="name" label="????" rules={[{ required: true }]}>
              <Input />
            </Form.Item>
            <Form.Item
              name="description"
              label="????????"
              rules={[{ required: true }]}
            >
              <Input.TextArea rows={3} />
            </Form.Item>
            <Form.Item name="price" label="????" rules={[{ required: true }]}>
              <InputNumber min={0} step={0.01} style={{ width: "100%" }} />
            </Form.Item>
            <Form.Item
              name="image"
              label="???? (URL)"
              rules={[{ required: true }]}
            >
              <Input placeholder="https://... or upload below" />
            </Form.Item>
            <Upload
              accept="image/*"
              showUploadList={false}
              customRequest={async (opts: any) => {
                const formData = new FormData();
                formData.append("file", opts.file);
                try {
                  const res = await api.post("/admin/upload", formData, {
                    headers: { "Content-Type": "multipart/form-data" },
                  });
                  productForm.setFieldsValue({ image: res.data.url });
                  message.success("Uploaded");
                  opts.onSuccess?.(res.data);
                } catch (e) {
                  message.error("Upload failed");
                  opts.onError?.(e);
                }
              }}
            >
              <Button icon={<UploadOutlined />}>Upload Image</Button>
            </Upload>
          </Form>
        </Modal>
      </Card>
    </div>
  );
}
