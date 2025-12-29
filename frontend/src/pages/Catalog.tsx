import { Button, Card, Col, Row, Typography } from "antd";
import { useEffect, useState } from "react";
import {
  fetchCategories,
  fetchDepartments,
  fetchProductsByCategory,
} from "../api/catalog";
import { Link, useSearchParams } from "react-router-dom";
import { useI18n } from "../store/I18nContext";
import { api } from "../api/client";

export default function Catalog() {
  const [depts, setDepts] = useState<any[]>([]);
  const [cats, setCats] = useState<any[]>([]);
  const [products, setProducts] = useState<any[]>([]);
  const [deptId, setDeptId] = useState<number | undefined>();
  const [catId, setCatId] = useState<number | undefined>();
  const { t } = useI18n();
  const [searchParams, setSearchParams] = useSearchParams();
  useEffect(() => {
    fetchDepartments().then(setDepts);
  }, []);
  useEffect(() => {
    const s = searchParams.get("dept");
    if (s) setDeptId(Number(s));
  }, [searchParams]);
  useEffect(() => {
    if (deptId) fetchCategories(deptId).then(setCats);
  }, [deptId]);
  useEffect(() => {
    if (catId) fetchProductsByCategory(catId).then(setProducts);
  }, [catId]);
  return (
    <div>
      <Typography.Title level={2}>{t("catalog.title")}</Typography.Title>
      {(deptId || catId) && (
        <div style={{ display: "flex", gap: 8, marginBottom: 12 }}>
          {deptId && (
            <Button
              size="small"
              onClick={() => {
                setDeptId(undefined);
                setCatId(undefined);
                setSearchParams((sp) => {
                  const next = new URLSearchParams(sp);
                  next.delete("dept");
                  return next;
                });
              }}
            >
              All Departments
            </Button>
          )}
          {deptId && catId && (
            <Button size="small" onClick={() => setCatId(undefined)}>
              All Categories
            </Button>
          )}
        </div>
      )}

      {!deptId && (
        <Row gutter={[16, 16]} style={{ marginBottom: 16 }}>
          {depts.map((d) => {
            const origin = (() => {
              try {
                return new URL(api.defaults.baseURL as string).origin;
              } catch {
                return "";
              }
            })();
            const img =
              d.image_url && !/^https?:/i.test(d.image_url)
                ? origin + d.image_url
                : d.image_url;
            return (
              <Col key={d.id} xs={24} sm={12} md={8}>
                <Card
                  hoverable
                  title={d.name}
                  cover={img ? <img src={img} alt={d.name} /> : undefined}
                  onClick={() => {
                    setDeptId(d.id);
                    setCatId(undefined);
                    setSearchParams((sp) => {
                      const next = new URLSearchParams(sp);
                      next.set("dept", String(d.id));
                      return next;
                    });
                  }}
                >
                  <p>{d.description}</p>
                </Card>
              </Col>
            );
          })}
        </Row>
      )}
      {!!deptId && cats.length > 0 && !catId && (
        <Row gutter={[12, 12]} style={{ marginBottom: 16 }}>
          {cats.map((c) => (
            <Col key={c.id} xs={12} sm={8} md={6}>
              <Card hoverable onClick={() => setCatId(c.id)} title={c.name} />
            </Col>
          ))}
        </Row>
      )}
      <Row gutter={[16, 16]}>
        {products.map((p) => {
          const origin = (() => {
            try {
              return new URL(api.defaults.baseURL as string).origin;
            } catch {
              return "";
            }
          })();
          const img =
            p.image_url && !/^https?:/i.test(p.image_url)
              ? origin + p.image_url
              : p.image_url;
          return (
            <Col key={p.id} xs={24} sm={12} md={8}>
              <Card size="small">
                <div style={{ display: "flex", gap: 12 }}>
                  {img ? (
                    <img
                      src={img}
                      alt={p.name}
                      style={{
                        width: 100,
                        height: 80,
                        objectFit: "cover",
                        borderRadius: 4,
                      }}
                    />
                  ) : null}
                  <div style={{ flex: 1 }}>
                    <div style={{ fontWeight: 600, marginBottom: 4 }}>
                      {p.name}
                    </div>
                    <div
                      style={{ color: "#666", marginBottom: 6, fontSize: 12 }}
                    >
                      {p.short_description}
                    </div>
                    <div
                      style={{
                        display: "flex",
                        gap: 12,
                        marginBottom: 6,
                        fontSize: 12,
                      }}
                    >
                      <span>
                        {t("product.base_price")}:{" "}
                        {Number(p.base_price).toFixed(2)}
                      </span>
                      <span>
                        {t("product.base_prod_time")}:{" "}
                        {p.base_production_time_days}
                      </span>
                    </div>
                    <Link to={`/product/${p.id}`}>{t("catalog.view")}</Link>
                  </div>
                </div>
              </Card>
            </Col>
          );
        })}
      </Row>
    </div>
  );
}
