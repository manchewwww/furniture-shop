import {
  Button,
  Card,
  Col,
  InputNumber,
  Row,
  Select,
  Typography,
  message,
} from "antd";
import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { fetchProduct, fetchRecommendations } from "../api/catalog";
import { useCart } from "../store/CartContext";
import { useI18n } from "../store/I18nContext";
import { getApiOrigin } from "../api/client";

export default function ProductDetails() {
  const { id } = useParams();
  const [product, setProduct] = useState<any>();
  const [rec, setRec] = useState<any[]>([]);
  const [qty, setQty] = useState<number>(1);
  const [selected, setSelected] = useState<number[]>([]);
  const { add } = useCart();
  const { t } = useI18n();
  useEffect(() => {
    if (id) {
      fetchProduct(Number(id)).then(setProduct);
      fetchRecommendations(Number(id)).then(setRec);
    }
  }, [id]);
  if (!product) return null;
  return (
    <div>
      <Row gutter={24}>
        <Col
          md={12}
          xs={24}
          style={{
            borderRight: "1px solid #f0f0f0",
            paddingRight: 16,
            marginBottom: 16,
          }}
        >
          {(() => {
            const origin = getApiOrigin();
            const img =
              product.image_url && !/^https?:/i.test(product.image_url)
                ? origin + product.image_url
                : product.image_url;
            return (
              <img
                src={img}
                alt={product.name}
                style={{
                  width: "100%",
                  maxHeight: 360,
                  objectFit: "cover",
                  borderRadius: 4,
                }}
              />
            );
          })()}
        </Col>
        <Col md={12} xs={24} style={{ paddingLeft: 16 }}>
          <Typography.Title level={3}>{product.name}</Typography.Title>
          <p>{product.long_description}</p>
          <p>
            {t("product.base_price")}: {product.base_price.toFixed(2)}
          </p>
          <p>
            {t("product.base_prod_time")}: {product.base_production_time_days}
          </p>
          <p>
            Dimensions (cm): {product.default_width}W × {product.default_height}
            H × {product.default_depth}D
          </p>
          <div style={{ margin: "12px 0" }}>
            <Typography.Text>{t("product.options")}:</Typography.Text>
            {product.options && product.options.length > 0 ? (
              <Select
                mode="multiple"
                placeholder={
                  t("product.select_options") ||
                  "Select options (colours, materials, extras)"
                }
                style={{ minWidth: 320 }}
                onChange={(vals) => setSelected(vals as number[])}
                options={(product.options || []).map((o: any) => ({
                  value: o.id,
                  label: `${o.option_name} (${o.option_type})`,
                }))}
              />
            ) : (
              <Typography.Text type="secondary" style={{ marginLeft: 8 }}>
                This product has no configurable options.
              </Typography.Text>
            )}
          </div>
          <div style={{ display: "flex", gap: 12, alignItems: "center" }}>
            <Typography.Text>{t("product.quantity")}:</Typography.Text>
            <InputNumber
              min={1}
              value={qty}
              onChange={(v) => setQty(Number(v))}
            />
            <Button
              type="primary"
              onClick={() => {
                add({
                  product,
                  quantity: qty,
                  options: selected.map((id) => ({ id, type: "extra" })),
                });
                message.success(t("product.added"));
              }}
            >
              {t("product.add_to_cart")}
            </Button>
          </div>
        </Col>
      </Row>
      {!!rec.length && (
        <div style={{ marginTop: 24 }}>
          <Typography.Title level={4} style={{ textAlign: "center" }}>
            {t("product.recommended")}
          </Typography.Title>
          <Row gutter={[16, 16]} justify="start">
            {rec.map((r) => {
              const origin = getApiOrigin();
              const rimg =
                r.image_url && !/^https?:/i.test(r.image_url)
                  ? origin + r.image_url
                  : r.image_url;
              return (
                <Col key={r.id} xs={12} sm={8} md={6} lg={6}>
                  <Link to={`/product/${r.id}`} style={{ display: "block" }}>
                    <Card
                      hoverable
                      bodyStyle={{ textAlign: "center" }}
                      cover={
                        rimg ? (
                          <img
                            src={rimg}
                            alt={r.name}
                            style={{
                              width: "100%",
                              height: 120,
                              objectFit: "cover",
                            }}
                          />
                        ) : undefined
                      }
                    >
                      <Typography.Text>{r.name}</Typography.Text>
                    </Card>
                  </Link>
                </Col>
              );
            })}
          </Row>
        </div>
      )}
    </div>
  );
}
