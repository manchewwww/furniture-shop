import { Button, Card, Form, Input, Modal, Table, Upload, message } from "antd";
import { UploadOutlined } from "@ant-design/icons";
import { useEffect, useState } from "react";
import { api } from "../api/client";
import { useI18n } from "../store/I18nContext";
import { Link, useNavigate } from "react-router-dom";

export default function AdminDepartments() {
  const { t } = useI18n();
  const nav = useNavigate();
  const [depts, setDepts] = useState<any[]>([]);
  const [openDept, setOpenDept] = useState(false);
  const [editing, setEditing] = useState<any | null>(null);
  const [deptForm] = Form.useForm();

  const load = async () => {
    try {
      const d = await api.get("/admin/departments");
      setDepts(d.data);
    } catch {
      message.error("Failed to load departments");
    }
  };

  useEffect(() => {
    load();
  }, []);

  const submitDept = async () => {
    const v = await deptForm.validateFields();
    if (editing) {
      await api.put(`/admin/departments/${editing.id}`, v);
    } else {
      await api.post("/admin/departments", v);
    }
    setOpenDept(false);
    setEditing(null);
    deptForm.resetFields();
    load();
  };

  return (
    <div>
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
            {
              title: t("actions"),
              render: (_: any, r: any) => (
                <>
                  <Button
                    size="small"
                    onClick={() => {
                      setEditing(r);
                      deptForm.setFieldsValue({
                        name: r.name,
                        description: r.description,
                        image_url: r.image_url,
                      });
                      setOpenDept(true);
                    }}
                    style={{ marginRight: 8 }}
                  >
                    Edit
                  </Button>
                  <Button
                    danger
                    size="small"
                    onClick={async () => {
                      await api.delete(`/admin/departments/${r.id}`);
                      load();
                    }}
                  >
                    Delete
                  </Button>
                </>
              ),
            },
          ]}
        />
        <Modal
          title={editing ? t("edit_department") : t("create_department")}
          open={openDept}
          onOk={submitDept}
          onCancel={() => {
            setOpenDept(false);
            setEditing(null);
            deptForm.resetFields();
          }}
        >
          <Form layout="vertical" form={deptForm}>
            <Form.Item
              name="name"
              label={t("department_name")}
              rules={[{ required: true }]}
            >
              <Input />
            </Form.Item>
            <Form.Item name="description" label={t("department_description")}>
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
                  deptForm.setFieldsValue({ image_url: finalUrl });
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
