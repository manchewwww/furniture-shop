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
import { useI18n } from "../store/I18nContext";

export default function AdminDashboard() {
  const [depts, setDepts] = useState<any[]>([]);
  const [orders, setOrders] = useState<any[]>([]);
  const [products, setProducts] = useState<any[]>([]);
  const [categories, setCategories] = useState<any[]>([]);
  const [openDept, setOpenDept] = useState(false);
  const [deptForm] = Form.useForm();
  const [openProduct, setOpenProduct] = useState(false);
  const [productForm] = Form.useForm();
  const [editing, setEditing] = useState<any | null>(null);
  const { t } = useI18n();
  const selectedDept: number | undefined = Form.useWatch(
    "department_id",
    productForm
  );

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
      message.error("Failed to load admin data");
    }
  };
  useEffect(() => {
    load();
  }, []);

  const createDept = async () => {
    const v = await deptForm.validateFields();
    await api.post("/admin/departments", v);
    setOpenDept(false);
    deptForm.resetFields();
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
      base_production_time_days: v.production_days ?? 0,
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
            label: t("admin.departments"),
            children: (
              <Card
                title={t("admin.departments")}
                extra={
                  <Button onClick={() => setOpenDept(true)}>
                    {t("admin.create_department")}
                  </Button>
                }
              >
                <Table
                  rowKey="id"
                  dataSource={depts}
                  columns={[
                    { title: t("admin.department_name"), dataIndex: "name" },
                    {
                      title: t("admin.department_description"),
                      dataIndex: "description",
                    },
                  ]}
                />
                <Modal
                  title={t("admin.create_department")}
                  open={openDept}
                  onOk={createDept}
                  onCancel={() => setOpenDept(false)}
                >
                  <Form layout="vertical" form={deptForm}>
                    <Form.Item
                      name="name"
                      label={t("admin.department_name")}
                      rules={[{ required: true }]}
                    >
                      <Input />
                    </Form.Item>
                    <Form.Item
                      name="description"
                      label={t("admin.department_description")}
                    >
                      <Input />
                    </Form.Item>
                  </Form>
                </Modal>
              </Card>
            ),
          },
          {
            key: "orders",
            label: t("admin.orders"),
            children: (
              <Card title={t("admin.orders_title")}>
                <Table
                  rowKey="id"
                  dataSource={orders}
                  columns={[
                    { title: t("orders.col.id"), dataIndex: "id" },
                    {
                      title: t("orders.col.status"),
                      dataIndex: "status",
                      render: (s: string) => <Tag>{s}</Tag>,
                    },
                    {
                      title: t("orders.col.payment_status"),
                      dataIndex: "payment_status",
                    },
                    { title: t("orders.col.total"), dataIndex: "total_price" },
                    {
                      title: t("admin.actions"),
                      render: (_: any, r: any) => (
                        <>
                          {[
                            "нов",
                            "в производство",
                            "готов за доставка",
                            "доставен",
                            "отказан",
                            "в обработка",
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
        title={t("admin.products")}
        style={{ marginTop: 16 }}
        extra={
          <Button
            onClick={() => {
              setEditing(null);
              productForm.resetFields();
              setOpenProduct(true);
            }}
          >
            {t("admin.create_product")}
          </Button>
        }
      >
        <Table
          rowKey="id"
          dataSource={products}
          columns={[
            { title: "ID", dataIndex: "id", width: 60 },
            { title: t("admin.product_name"), dataIndex: "name" },
            {
              title: t("admin.product_description"),
              dataIndex: "short_description",
            },
            { title: t("admin.product_price"), dataIndex: "base_price" },
            {
              title: t("orders.col.eta_days"),
              dataIndex: "base_production_time_days",
            },
            {
              title: t("admin.product_image"),
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
              title: t("admin.actions"),
              render: (_: any, r: any) => (
                <>
                  <Button
                    size="small"
                    style={{ marginRight: 8 }}
                    onClick={() => {
                      setEditing(r);
                      setOpenProduct(true);
                      productForm.setFieldsValue({
                        department_id: categories.find(
                          (c: any) => c.id === r.category_id
                        )?.department_id,
                        category_id: r.category_id,
                        name: r.name,
                        description: r.short_description || r.long_description,
                        price: r.base_price,
                        production_days: r.base_production_time_days,
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
          title={editing ? t("admin.edit_product") : t("admin.create_product")}
          open={openProduct}
          onOk={submitProduct}
          onCancel={() => {
            setOpenProduct(false);
            setEditing(null);
          }}
        >
          <Form layout="vertical" form={productForm}>
            <Form.Item name="department_id" label={t("admin.department")}>
              <Select
                placeholder={t("admin.department")}
                onChange={() =>
                  productForm.setFieldsValue({ category_id: undefined })
                }
              >
                {depts.map((d: any) => (
                  <Select.Option key={d.id} value={d.id}>
                    {d.name}
                  </Select.Option>
                ))}
              </Select>
            </Form.Item>
            <Form.Item
              name="category_id"
              label={t("admin.category")}
              rules={[{ required: true }]}
            >
              <Select placeholder={t("admin.category")}>
                {categories
                  .filter(
                    (c: any) =>
                      !selectedDept || c.department_id === selectedDept
                  )
                  .map((c: any) => (
                    <Select.Option key={c.id} value={c.id}>
                      {c.name}
                    </Select.Option>
                  ))}
              </Select>
            </Form.Item>
            <Form.Item
              name="name"
              label={t("admin.product_name")}
              rules={[{ required: true }]}
            >
              <Input />
            </Form.Item>
            <Form.Item
              name="description"
              label={t("admin.product_description")}
              rules={[{ required: true }]}
            >
              <Input.TextArea rows={3} />
            </Form.Item>
            <Form.Item
              name="price"
              label={t("admin.product_price")}
              rules={[{ required: true }]}
            >
              <InputNumber min={0} step={0.01} style={{ width: "100%" }} />
            </Form.Item>
            <Form.Item
              name="production_days"
              label={t("admin.product_production_days")}
              rules={[{ required: true }]}
            >
              <InputNumber min={0} step={1} style={{ width: "100%" }} />
            </Form.Item>
            <Form.Item
              name="image"
              label={t("admin.product_image")}
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
                  message.success(t("admin.upload_success"));
                  opts.onSuccess?.(res.data);
                } catch (e) {
                  message.error(t("admin.upload_fail"));
                  opts.onError?.(e);
                }
              }}
            >
              <Button icon={<UploadOutlined />}>
                {t("admin.upload_image")}
              </Button>
            </Upload>
          </Form>
        </Modal>
      </Card>
    </div>
  );
}
