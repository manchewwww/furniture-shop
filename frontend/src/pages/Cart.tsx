import { Button, Card, List, Typography } from "antd";
import { Link } from "react-router-dom";
import { useCart } from "../store/CartContext";
import { useI18n } from "../store/I18nContext";

export default function Cart() {
  const { items, remove, clear } = useCart();
  const { t } = useI18n();
  const total = items.reduce(
    (s, it) => s + it.product.base_price * it.quantity,
    0
  );
  return (
    <div>
      <Typography.Title level={2}>{t("cart.title")}</Typography.Title>
      <List
        dataSource={items}
        renderItem={(it) => (
          <List.Item
            actions={[
              <Button danger onClick={() => remove(it.product.id)}>
                Remove
              </Button>,
            ]}
          >
            <List.Item.Meta
              title={it.product.name}
              description={`${t("product.quantity")}: ${it.quantity}`}
            />
            <div>{(it.product.base_price * it.quantity).toFixed(2)}</div>
          </List.Item>
        )}
      />
      <Card>
        <p>
          {t("orders.col.total")}: {total.toFixed(2)}
        </p>
        <Button type="primary" disabled={items.length === 0}>
          <Link to="/checkout">{t("checkout.title")}</Link>
        </Button>
        <Button style={{ marginLeft: 8 }} onClick={clear}>
          Clear
        </Button>
      </Card>
    </div>
  );
}
