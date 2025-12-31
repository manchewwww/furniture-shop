import { Button, Card, Empty, Row, Col, Space, Table, Typography } from "antd";
import { Link, useNavigate } from "react-router-dom";
import { useCart } from "../store/CartContext";
import { useI18n } from "../store/I18nContext";
import { getApiOrigin } from "../api/client";

function unitPrice(product: any, selected: { id: number; type: string }[]) {
  let price = Number(product.base_price || 0);
  const byId: Record<number, any> = {};
  (product.options || []).forEach((o: any) => (byId[o.id] = o));
  selected.forEach((so) => {
    const opt = byId[so.id];
    if (!opt) return;
    if (opt.price_modifier_type === "absolute") {
      price += Number(opt.price_modifier_value || 0);
    } else if (opt.price_modifier_type === "percent") {
      price = price * (1 + Number(opt.price_modifier_value || 0) / 100);
    }
  });
  return price;
}

function itemProductionDays(
  product: any,
  selected: { id: number; type: string }[]
) {
  let days = Number(product.base_production_time_days || 0) || 1;
  const byId: Record<number, any> = {};
  (product.options || []).forEach((o: any) => (byId[o.id] = o));
  selected.forEach((so) => {
    const opt = byId[so.id];
    if (!opt) return;
    days += Number(opt.production_time_modifier_days || 0);
    if (opt.production_time_modifier_percent != null) {
      days = Math.round(
        days * (1 + Number(opt.production_time_modifier_percent) / 100)
      );
    }
  });
  return Math.max(1, days);
}

export default function Cart() {
  const { items, remove, clear, increment, decrement } = useCart();
  const { t } = useI18n();
  const navigate = useNavigate();

  const totals = items.map(
    (it) => unitPrice(it.product, it.options || []) * it.quantity
  );
  const total = totals.reduce((a, b) => a + b, 0);
  const etaDays = items.reduce(
    (max, it) =>
      Math.max(max, itemProductionDays(it.product, it.options || [])),
    0
  );

  return (
    <div>
      <Typography.Title level={2}>{t("cart.title")}</Typography.Title>
      {items.length === 0 ? (
        <Empty description="Your cart is empty">
          <Button type="primary" onClick={() => navigate("/catalog")}>
            Go to catalog
          </Button>
        </Empty>
      ) : (
        <Row gutter={16}>
          <Col xs={24} md={16}>
            <Table
              rowKey={(r) => String(r.product.id)}
              dataSource={items}
              pagination={false}
              columns={[
                {
                  title: "Product",
                  render: (_: any, it: any) => {
                    const origin = getApiOrigin();
                    const img =
                      it.product.image_url &&
                      !/^https?:/i.test(it.product.image_url)
                        ? origin + it.product.image_url
                        : it.product.image_url;
                    return (
                      <Space align="start">
                        {img ? (
                          <img
                            src={img}
                            alt={it.product.name}
                            style={{
                              width: 64,
                              height: 48,
                              objectFit: "cover",
                              borderRadius: 4,
                            }}
                          />
                        ) : null}
                        <div>
                          <Link to={`/product/${it.product.id}`}>
                            {it.product.name}
                          </Link>
                        </div>
                      </Space>
                    );
                  },
                },
                {
                  title: "Unit Price",
                  align: "right" as const,
                  render: (_: any, it: any) =>
                    unitPrice(it.product, it.options || []).toFixed(2),
                },
                {
                  title: "Quantity",
                  align: "center" as const,
                  render: (_: any, it: any) => (
                    <Space>
                      <Button
                        size="small"
                        onClick={() => decrement(it.product.id)}
                        disabled={it.quantity <= 1}
                      >
                        -
                      </Button>
                      <span style={{ minWidth: 16, textAlign: "center" }}>
                        {it.quantity}
                      </span>
                      <Button
                        size="small"
                        onClick={() => increment(it.product.id)}
                      >
                        +
                      </Button>
                    </Space>
                  ),
                },
                {
                  title: "Line Total",
                  align: "right" as const,
                  render: (_: any, it: any) =>
                    (
                      unitPrice(it.product, it.options || []) * it.quantity
                    ).toFixed(2),
                },
                {
                  title: "",
                  align: "right" as const,
                  render: (_: any, it: any) => (
                    <Button
                      danger
                      size="small"
                      onClick={() => remove(it.product.id)}
                    >
                      Remove
                    </Button>
                  ),
                },
              ]}
            />
          </Col>
          <Col xs={24} md={8}>
            <Card>
              <Typography.Title level={4} style={{ marginTop: 0 }}>
                Summary
              </Typography.Title>
              <div style={{ display: "flex", justifyContent: "space-between" }}>
                <span>{t("orders.col.total")}:</span>
                <span>{total.toFixed(2)}</span>
              </div>
              <div style={{ display: "flex", justifyContent: "space-between" }}>
                <span>{t("orders.col.eta_days")}:</span>
                <span>{etaDays || 0}</span>
              </div>
              <div style={{ marginTop: 12 }}>
                <Button
                  type="primary"
                  block
                  onClick={() => navigate("/checkout")}
                >
                  {t("checkout.title")}
                </Button>
                <Button danger block style={{ marginTop: 8 }} onClick={clear}>
                  Clear Cart
                </Button>
              </div>
            </Card>
          </Col>
        </Row>
      )}
    </div>
  );
}
