import { Card, Col, Row, Select, Typography } from "antd";
import { useEffect, useState } from "react";
import {
  fetchCategories,
  fetchDepartments,
  fetchProductsByCategory,
} from "../api/catalog";
import { Link } from "react-router-dom";
import { useI18n } from "../store/I18nContext";

export default function Catalog() {
  const [depts, setDepts] = useState<any[]>([]);
  const [cats, setCats] = useState<any[]>([]);
  const [products, setProducts] = useState<any[]>([]);
  const [deptId, setDeptId] = useState<number | undefined>();
  const [catId, setCatId] = useState<number | undefined>();
  const { t } = useI18n();
  useEffect(() => {
    fetchDepartments().then(setDepts);
  }, []);
  useEffect(() => {
    if (deptId) fetchCategories(deptId).then(setCats);
  }, [deptId]);
  useEffect(() => {
    if (catId) fetchProductsByCategory(catId).then(setProducts);
  }, [catId]);
  return (
    <div>
      <Typography.Title level={2}>{t("catalog.title")}</Typography.Title>
      <div style={{ display: "flex", gap: 12, marginBottom: 16 }}>
        <Select
          placeholder={t("catalog.select.department")}
          style={{ minWidth: 220 }}
          onChange={setDeptId}
          options={depts.map((d) => ({ value: d.id, label: d.name }))}
        />
        <Select
          placeholder={t("catalog.select.category")}
          style={{ minWidth: 220 }}
          onChange={setCatId}
          options={cats.map((c) => ({ value: c.id, label: c.name }))}
        />
      </div>
      <Row gutter={[16, 16]}>
        {products.map((p) => (
          <Col key={p.id} xs={24} sm={12} md={8}>
            <Card title={p.name} cover={<img src={p.image_url} alt={p.name} />}>
              <p>{p.short_description}</p>
              <Link to={`/product/${p.id}`}>{t("catalog.view")}</Link>
            </Card>
          </Col>
        ))}
      </Row>
    </div>
  );
}
