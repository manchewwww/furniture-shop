import {
  Button,
  Card,
  Form,
  Input,
  InputNumber,
  Modal,
  Popconfirm,
  Select,
  Table,
  Upload,
  message,
} from "antd";
import { UploadOutlined } from "@ant-design/icons";
import { useEffect, useMemo, useState } from "react";
import { api } from "../api/client";
import { useI18n } from "../store/I18nContext";
import { useNavigate } from "react-router-dom";

export default function AdminProducts() {
  const { t } = useI18n();
  const nav = useNavigate();
  const [depts, setDepts] = useState<any[]>([]);
  const [categories, setCategories] = useState<any[]>([]);
  const [products, setProducts] = useState<any[]>([]);
  const [openProduct, setOpenProduct] = useState(false);
  const [productForm] = Form.useForm();
  const [editing, setEditing] = useState<any | null>(null);
  const [colorOptions, setColorOptions] = useState<string[]>([]);
  const [materialOptions, setMaterialOptions] = useState<string[]>([]);
  const [extraOptions, setExtraOptions] = useState<string[]>([]);
  const [imagePreview, setImagePreview] = useState<string | null>(null);
  const selectedDept: number | undefined = Form.useWatch(
    "department_id",
    productForm
  );
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

  const load = async () => {
    try {
      const [d, c, p] = await Promise.all([
        api.get("/admin/departments"),
        api.get("/admin/categories"),
        api.get("/admin/products"),
      ]);
      setDepts(d.data);
      setCategories(c.data);
      setProducts(p.data);
    } catch {
      message.error("Failed to load products");
    }
  };

  useEffect(() => {
    load();
  }, []);

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
      default_width: v.default_width,
      default_height: v.default_height,
      default_depth: v.default_depth,
      base_material: v.base_material,
      quantity: v.quantity ?? 0,
    };
    if (editing) {
      await api.put(`/admin/products/${editing.id}`, payload);
      try {
        const existing = await api.get(`/admin/product_options`, {
          params: { product_id: editing.id },
        });
        const toDelete = (existing.data || []).filter((o: any) =>
          ["color", "material", "size", "extra"].includes(o.option_type)
        );
        await Promise.all(
          toDelete.map((o: any) => api.delete(`/admin/product_options/${o.id}`))
        );
      } catch {}
      const posts: Promise<any>[] = [];
      colorOptions.forEach((name) =>
        posts.push(
          api.post(`/admin/product_options`, {
            product_id: editing.id,
            option_type: "color",
            option_name: name,
            price_modifier_type: "absolute",
            price_modifier_value: 0,
          })
        )
      );
      materialOptions.forEach((name) =>
        posts.push(
          api.post(`/admin/product_options`, {
            product_id: editing.id,
            option_type: "material",
            option_name: name,
            price_modifier_type: "absolute",
            price_modifier_value: 0,
          })
        )
      );
      extraOptions.forEach((name) =>
        posts.push(
          api.post(`/admin/product_options`, {
            product_id: editing.id,
            option_type: "extra",
            option_name: name,
            price_modifier_type: "absolute",
            price_modifier_value: 0,
          })
        )
      );
      if (posts.length) await Promise.all(posts);
    } else {
      const created = await api.post("/admin/products", payload);
      const newId = created.data?.id;
      if (newId) {
        const posts: Promise<any>[] = [];
        colorOptions.forEach((name) =>
          posts.push(
            api.post(`/admin/product_options`, {
              product_id: newId,
              option_type: "color",
              option_name: name,
              price_modifier_type: "absolute",
              price_modifier_value: 0,
            })
          )
        );
        materialOptions.forEach((name) =>
          posts.push(
            api.post(`/admin/product_options`, {
              product_id: newId,
              option_type: "material",
              option_name: name,
              price_modifier_type: "absolute",
              price_modifier_value: 0,
            })
          )
        );
        extraOptions.forEach((name) =>
          posts.push(
            api.post(`/admin/product_options`, {
              product_id: newId,
              option_type: "extra",
              option_name: name,
              price_modifier_type: "absolute",
              price_modifier_value: 0,
            })
          )
        );
        if (posts.length) await Promise.all(posts);
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
      <Card
        title={t("products")}
        extra={
          <Button
            onClick={() => {
              setEditing(null);
              setColorOptions([]);
              setMaterialOptions([]);
              setExtraOptions([]);
              setImagePreview(null);
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
            { title: t("product_description"), dataIndex: "short_description" },
            { title: "Quantity", dataIndex: "quantity" },
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
                      const origin = (() => {
                        try {
                          return new URL(api.defaults.baseURL as string).origin;
                        } catch {
                          return "";
                        }
                      })();
                      const imageVal =
                        r.image_url && !/^https?:/i.test(r.image_url)
                          ? origin + r.image_url
                          : r.image_url;
                      productForm.setFieldsValue({
                        department_id: categories.find(
                          (c: any) => c.id === r.category_id
                        )?.department_id,
                        category_id: r.category_id,
                        name: r.name,
                        description: r.short_description || r.long_description,
                        price: r.base_price,
                        quantity: r.quantity,
                        production_days: r.base_production_time_days,
                        image: imageVal,
                        default_width: r.default_width,
                        default_height: r.default_height,
                        default_depth: r.default_depth,
                        base_material: r.base_material,
                      });
                      setImagePreview(imageVal || null);
                      api
                        .get(`/admin/product_options`, {
                          params: { product_id: r.id },
                        })
                        .then((res) => {
                          const items = (res.data || []) as any[];
                          setColorOptions(
                            items
                              .filter((o) => o.option_type === "color")
                              .map((o) => o.option_name)
                          );
                          setMaterialOptions(
                            items
                              .filter((o) => o.option_type === "material")
                              .map((o) => o.option_name)
                          );
                          setExtraOptions(
                            items
                              .filter((o) => o.option_type === "extra")
                              .map((o) => o.option_name)
                          );
                        })
                        .catch(() => {
                          setColorOptions([]);
                          setMaterialOptions([]);
                          setExtraOptions([]);
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
          title={editing ? t("edit_product") : t("create_product")}
          open={openProduct}
          onOk={submitProduct}
          onCancel={() => {
            setOpenProduct(false);
            setEditing(null);
            setColorOptions([]);
            setMaterialOptions([]);
            setExtraOptions([]);
            setImagePreview(null);
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
              name="quantity"
              label="Quantity"
              rules={[{ required: true }]}
            >
              <InputNumber min={0} step={1} style={{ width: "100%" }} />
            </Form.Item>
            <Form.Item label="Dimensions (cm)">
              <div style={{ display: "flex", gap: 8 }}>
                <Form.Item
                  name="default_width"
                  rules={[{ required: true }]}
                  style={{ marginBottom: 0, flex: 1 }}
                >
                  <InputNumber
                    min={1}
                    placeholder="Width"
                    style={{ width: "100%" }}
                  />
                </Form.Item>
                <Form.Item
                  name="default_height"
                  rules={[{ required: true }]}
                  style={{ marginBottom: 0, flex: 1 }}
                >
                  <InputNumber
                    min={1}
                    placeholder="Height"
                    style={{ width: "100%" }}
                  />
                </Form.Item>
                <Form.Item
                  name="default_depth"
                  rules={[{ required: true }]}
                  style={{ marginBottom: 0, flex: 1 }}
                >
                  <InputNumber
                    min={1}
                    placeholder="Depth"
                    style={{ width: "100%" }}
                  />
                </Form.Item>
              </div>
            </Form.Item>
            <Form.Item
              name="production_days"
              label={t("product_production_days")}
              rules={[{ required: true }]}
            >
              <InputNumber min={0} step={1} style={{ width: "100%" }} />
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
            <Form.Item label="Materials">
              <Select
                mode="tags"
                value={materialOptions}
                onChange={(vals) => setMaterialOptions(vals as string[])}
                placeholder="Type material names and press Enter"
                tokenSeparators={[",", " ", ";"]}
                allowClear
              />
            </Form.Item>
            <Form.Item name="base_material" label="Base Material">
              <Input placeholder="e.g., MDF, Wood, Metal" />
            </Form.Item>

            <Form.Item label="Initial Stock (optional)">
              <Form.Item name="stock_quantity" style={{ marginBottom: 0 }}>
                <InputNumber min={0} placeholder="Quantity" />
              </Form.Item>
            </Form.Item>
            <Form.Item label="Extras">
              <Select
                mode="tags"
                value={extraOptions}
                onChange={(vals) => setExtraOptions(vals as string[])}
                placeholder="Type extras and press Enter"
                tokenSeparators={[",", " ", ";"]}
                allowClear
              />
            </Form.Item>
            <Form.Item label={t("product_image")}>
              {imagePreview ? (
                <div
                  style={{
                    display: "flex",
                    alignItems: "center",
                    gap: 12,
                    marginBottom: 8,
                  }}
                >
                  <img
                    src={imagePreview}
                    alt="preview"
                    style={{
                      width: 160,
                      height: 120,
                      objectFit: "cover",
                      borderRadius: 4,
                      flex: "0 0 auto",
                    }}
                  />
                  <div
                    style={{
                      fontSize: 12,
                      color: "#555",
                      wordBreak: "break-all",
                    }}
                  >
                    {imagePreview}
                  </div>
                </div>
              ) : null}
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
                  setImagePreview(finalUrl);
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
