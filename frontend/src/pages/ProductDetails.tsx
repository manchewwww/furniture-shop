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
import { useParams } from "react-router-dom";
import { fetchProduct, fetchRecommendations } from "../api/catalog";
import { useCart } from "../store/CartContext";
import { useI18n } from "../store/I18nContext";

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
        <Col md={14} xs={24}>
          <img
            src={product.image_url}
            alt={product.name}
            style={{ maxWidth: "100%" }}
          />
          <Typography.Title level={3}>{product.name}</Typography.Title>
          <p>{product.long_description}</p>
          <p>
            {t("product.base_price")}: {product.base_price.toFixed(2)}
          </p>
          <p>
            {t("product.base_prod_time")}: {product.base_production_time_days}
          </p>
          <div style={{ margin: "12px 0" }}>
            <Typography.Text>{t("product.options")}:</Typography.Text>
            <Select
              mode="multiple"
              placeholder={t("product.select_options")}
              style={{ minWidth: 320 }}
              onChange={(vals) => setSelected(vals as number[])}
              options={(product.options || []).map((o: any) => ({
                value: o.id,
                label: `${o.option_name} (${o.option_type})`,
              }))}
            />
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
        <Col md={10} xs={24}>
          <Typography.Title level={4}>
            {t("product.recommended")}
          </Typography.Title>
          {rec.map((r) => (
            <Card
              key={r.id}
              size="small"
              style={{ marginBottom: 8 }}
              title={r.name}
            ></Card>
          ))}
        </Col>
      </Row>
    </div>
  );
}
