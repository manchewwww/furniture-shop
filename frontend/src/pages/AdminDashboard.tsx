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
import { useEffect, useMemo, useState } from "react";
import { api } from "../api/client";
import { useI18n } from "../store/I18nContext";

export default function AdminDashboard() {
  const [depts, setDepts] = useState<any[]>([]);
  const [orders, setOrders] = useState<any[]>([]);
  const [products, setProducts] = useState<any[]>([]);
  const [categories, setCategories] = useState<any[]>([]);
  const [openDept, setOpenDept] = useState(false);
  const [deptForm] = Form.useForm();
  const [openCategory, setOpenCategory] = useState(false);
  const [categoryForm] = Form.useForm();
  const [openProduct, setOpenProduct] = useState(false);
  const [productForm] = Form.useForm();
  const [editing, setEditing] = useState<any | null>(null);
  const [colorOptions, setColorOptions] = useState<string[]>([]);
  const commonColours = useMemo(
    () => [
      "White",
      "Black",
      "Gray",
      "Silver",
      "Beige",
      "Brown",
      "Oak",
      "Walnut",
      "Mahogany",
      "Pine",
      "Blue",
      "Navy",
      "Sky Blue",
      "Red",
      "Burgundy",
      "Green",
      "Olive",
      "Teal",
      "Yellow",
      "Gold",
      "Orange",
      "Purple",
      "Pink",
    ],
    []
  );
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

  const createCategory = async () => {
    const v = await categoryForm.validateFields();
    await api.post("/admin/categories", v);
    setOpenCategory(false);
    categoryForm.resetFields();
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
      try {
        const existing = await api.get(`/admin/product_options`, {
          params: { product_id: editing.id },
        });
        const colors = (existing.data || []).filter(
          (o: any) => o.option_type === "color"
        );
        await Promise.all(
          colors.map((o: any) => api.delete(`/admin/product_options/${o.id}`))
        );
      } catch {}
      if (colorOptions.length) {
        await Promise.all(
          colorOptions.map((name) =>
            api.post(`/admin/product_options`, {
              product_id: editing.id,
              option_type: "color",
              option_name: name,
              price_modifier_type: "absolute",
              price_modifier_value: 0,
            })
          )
        );
      }
    } else {
      const created = await api.post("/admin/products", payload);
      const newId = created.data?.id;
      if (newId && colorOptions.length) {
        await Promise.all(
          colorOptions.map((name) =>
            api.post(`/admin/product_options`, {
              product_id: newId,
              option_type: "color",
              option_name: name,
              price_modifier_type: "absolute",
              price_modifier_value: 0,
            })
          )
        );
      }
    }
    setOpenProduct(false);
    setEditing(null);
    productForm.resetFields();
    setColorOptions([]);
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
            label: t("departments"),
            children: (
              <Card
                title={t("departments")}
                extra={
                  <Button onClick={() => setOpenDept(true)}>
                    {t("create_department")}
                  </Button>
                }
              >
                <Table
                  rowKey="id"
                  dataSource={depts}
                  columns={[
                    { title: t("department_name"), dataIndex: "name" },
                    {
                      title: t("department_description"),
                      dataIndex: "description",
                    },
                  ]}
                />
                <Modal
                  title={t("create_department")}
                  open={openDept}
                  onOk={createDept}
                  onCancel={() => setOpenDept(false)}
                >
                  <Form layout="vertical" form={deptForm}>
                    <Form.Item
                      name="name"
                      label={t("department_name")}
                      rules={[{ required: true }]}
                    >
                      <Input />
                    </Form.Item>
                    <Form.Item
                      name="description"
                      label={t("department_description")}
                    >
                      <Input />
                    </Form.Item>
                    <Form.Item name="image_url" label={t("department_image")}>
                      <Input placeholder="data:image/...;base64,... or upload below" />
                    </Form.Item>
                    <Upload
                      accept="image/*"
                      showUploadList={false}
                      customRequest={async (opts: any) => {
                        const formData = new FormData();
                        formData.append("file", opts.file);
                        try {
                          const res = await api.post(
                            "/admin/upload",
                            formData,
                            {
                              headers: {
                                "Content-Type": "multipart/form-data",
                              },
                            }
                          );
                          const base = (api.defaults.baseURL as string) || "";
                          let origin = "";
                          try {
                            origin = new URL(base).origin;
                          } catch {}
                          const finalUrl = /^https?:/i.test(res.data.url)
                            ? res.data.url
                            : origin + res.data.url;
                          deptForm.setFieldsValue({ image_url: finalUrl });
                          message.success(t("upload_success"));
                          opts.onSuccess?.(res.data);
                        } catch (e) {
                          message.error(t("upload_fail"));
                          opts.onError?.(e);
                        }
                      }}
                    >
                      <Button icon={<UploadOutlined />}>
                        {t("upload_image")}
                      </Button>
                    </Upload>
                  </Form>
                </Modal>
              </Card>
            ),
          },
          {
            key: "categories",
            label: "Categories",
            children: (
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
                    {
                      title: t("category_description"),
                      dataIndex: "description",
                    },
                    {
                      title: t("department"),
                      dataIndex: "department_id",
                      render: (id: number) =>
                        depts.find((d) => d.id === id)?.name || id,
                    },
                    {
                      title: t("actions"),
                      render: (_: any, r: any) => (
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
                      ),
                    },
                  ]}
                />
                <Modal
                  title={t("create_category")}
                  open={openCategory}
                  onOk={createCategory}
                  onCancel={() => setOpenCategory(false)}
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
                    <Form.Item
                      name="description"
                      label={t("category_description")}
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
            label: t("orders"),
            children: (
              <Card title={t("orders_title")}>
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
                      title: t("actions"),
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
        title={t("products")}
        style={{ marginTop: 16 }}
        extra={
          <Button
            onClick={() => {
              setEditing(null);
              setColorOptions([]);
              productForm.resetFields();
              setOpenProduct(true);
            }}
          >
            {t("create_product")}
          </Button>
        }
      >
        <Table
          rowKey="id"
          dataSource={products}
          columns={[
            { title: "ID", dataIndex: "id", width: 60 },
            { title: t("product_name"), dataIndex: "name" },
            {
              title: t("product_description"),
              dataIndex: "short_description",
            },
            { title: t("product_price"), dataIndex: "base_price" },
            {
              title: t("orders.col.eta_days"),
              dataIndex: "base_production_time_days",
            },
            {
              title: t("product_image"),
              dataIndex: "image_url",
              render: (u: string) => {
                if (!u) return null;
                let origin = "";
                try {
                  origin = new URL(api.defaults.baseURL as string).origin;
                } catch {}
                const url = /^https?:/i.test(u) ? u : origin + u;
                return (
                  <img
                    src={url}
                    alt=""
                    style={{ width: 60, height: 40, objectFit: "cover" }}
                  />
                );
              },
            },
            {
              title: t("actions"),
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
                      api
                        .get(`/admin/product_options`, {
                          params: { product_id: r.id },
                        })
                        .then((res) => {
                          const colors = (res.data || [])
                            .filter((o: any) => o.option_type === "color")
                            .map((o: any) => o.option_name);
                          setColorOptions(colors);
                        })
                        .catch(() => setColorOptions([]));
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
          title={editing ? t("edit_product") : t("create_product")}
          open={openProduct}
          onOk={submitProduct}
          onCancel={() => {
            setOpenProduct(false);
            setEditing(null);
            setColorOptions([]);
          }}
        >
          <Form layout="vertical" form={productForm}>
            <Form.Item name="department_id" label={t("department")}>
              <Select
                placeholder={t("department")}
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
              label={t("category")}
              rules={[{ required: true }]}
            >
              <Select placeholder={t("category")}>
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
              label={t("product_name")}
              rules={[{ required: true }]}
            >
              <Input />
            </Form.Item>
            <Form.Item
              name="description"
              label={t("product_description")}
              rules={[{ required: true }]}
            >
              <Input.TextArea rows={3} />
            </Form.Item>
            <Form.Item
              name="price"
              label={t("product_price")}
              rules={[{ required: true }]}
            >
              <InputNumber min={0} step={0.01} style={{ width: "100%" }} />
            </Form.Item>
            <Form.Item
              name="production_days"
              label={t("product_production_days")}
              rules={[{ required: true }]}
            >
              <InputNumber min={0} step={1} style={{ width: "100%" }} />
            </Form.Item>
            <Form.Item
              name="image"
              label={t("product_image")}
              rules={[{ required: true }]}
            >
              <Input placeholder="https://... or upload below" />
            </Form.Item>
            <Form.Item label="Colours">
              <Select
                mode="tags"
                value={colorOptions}
                onChange={(vals) => setColorOptions(vals as string[])}
                options={commonColours.map((c) => ({ label: c, value: c }))}
                placeholder="Type a colour and press Enter (or choose from list)"
                tokenSeparators={[",", " ", ";"]}
                allowClear
              />
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
                  const base = (api.defaults.baseURL as string) || "";
                  let origin = "";
                  try {
                    origin = new URL(base).origin;
                  } catch {}
                  const finalUrl = /^https?:/i.test(res.data.url)
                    ? res.data.url
                    : origin + res.data.url;
                  productForm.setFieldsValue({ image: finalUrl });
                  message.success(t("upload_success"));
                  opts.onSuccess?.(res.data);
                } catch (e) {
                  message.error(t("upload_fail"));
                  opts.onError?.(e);
                }
              }}
            >
              <Button icon={<UploadOutlined />}>{t("upload_image")}</Button>
            </Upload>
          </Form>
        </Modal>
      </Card>
    </div>
  );
}
