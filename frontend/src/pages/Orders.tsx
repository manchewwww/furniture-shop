import {
  Button,
  Descriptions,
  Space,
  Table,
  Tag,
  Typography,
  message,
} from "antd";
import type { ColumnsType } from "antd/es/table";
import React, {
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import { Link, useSearchParams } from "react-router-dom";

import { myOrder, myOrders, payOrder } from "../api/orders";
import { fetchProduct } from "../api/catalog";
import { api } from "../api/client";
import { useI18n } from "../store/I18nContext";

type OrderRow = {
  id: number;
  status: string;
  payment_status: "paid" | string;
  total_price: number | string;
  estimated_production_time_days?: number;
  created_at?: string;
};

type OrderItem = {
  id?: number;
  product_id: number;
  quantity: number;
  unit_price: number | string;
  line_total: number | string;
};

type OrderDetails = OrderRow & {
  items?: OrderItem[];
  address?: string;
  phone?: string;
  created_at?: string;
  estimated_production_time_days?: number;
  total_price?: number | string;
};

type Product = {
  id?: number;
  name?: string;
  image_url?: string;
};

export default function Orders() {
  const { t } = useI18n();
  const [sp] = useSearchParams();
  const openParam = sp.get("open");
  const openId = useMemo(() => Number(openParam || 0) || 0, [openParam]);

  const [orders, setOrders] = useState<OrderRow[]>([]);
  const [detailsById, setDetailsById] = useState<
    Record<number, OrderDetails | undefined>
  >({});
  const [loadingDetailsById, setLoadingDetailsById] = useState<
    Record<number, boolean>
  >({});
  const [productCache, setProductCache] = useState<Record<number, Product>>({});

  const inFlightDetails = useRef<Set<number>>(new Set());
  const inFlightProducts = useRef<Set<number>>(new Set());
  const [expandedRowKeys, setExpandedRowKeys] = useState<React.Key[]>([]);

  const apiOrigin = useMemo(() => {
    try {
      return new URL(api.defaults.baseURL as string).origin;
    } catch {
      return "";
    }
  }, []);

  useEffect(() => {
    let alive = true;
    myOrders()
      .then((res) => {
        if (alive) setOrders(res as OrderRow[]);
      })
      .catch(() => message.error(t("orders.error")));
    return () => {
      alive = false;
    };
  }, [t]);

  // Auto-expand order passed via query param `open`
  useEffect(() => {
    if (!openId) return;
    const exists = orders.some((o) => Number(o.id) === Number(openId));
    if (exists) {
      setExpandedRowKeys([openId]);
      loadDetails(openId);
    }
  }, [openId, orders]);

  const startPayment = useCallback(async (orderId: number) => {
    try {
      const res = await payOrder(orderId);
      if (res?.checkout_url) window.location.assign(res.checkout_url);
    } catch {
      message.error("Unable to start payment");
    }
  }, []);

  const orderColumns: ColumnsType<OrderRow> = useMemo(
    () => [
      { title: t("orders.col.id"), dataIndex: "id" },
      {
        title: t("orders.col.status"),
        dataIndex: "status",
        render: (s: string) => <Tag>{s}</Tag>,
      },
      { title: t("orders.col.payment_status"), dataIndex: "payment_status" },
      { title: t("orders.col.total"), dataIndex: "total_price" },
      {
        title: t("orders.col.eta_days"),
        dataIndex: "estimated_production_time_days",
      },
      {
        title: "Actions",
        render: (_: unknown, row: OrderRow) => {
          const canPay = row.payment_status !== "paid";
          return (
            <Button
              size="small"
              type="primary"
              disabled={!canPay}
              onClick={() => startPayment(row.id)}
            >
              Re-pay
            </Button>
          );
        },
      },
    ],
    [startPayment, t]
  );

  const ensureProductsCached = useCallback(
    async (productIds: number[]) => {
      const unique = Array.from(
        new Set(productIds.map((x) => Number(x)).filter(Boolean))
      );

      const missing = unique.filter(
        (pid) => !productCache[pid] && !inFlightProducts.current.has(pid)
      );
      if (!missing.length) return;

      missing.forEach((pid) => inFlightProducts.current.add(pid));

      const results = await Promise.all(
        missing.map(async (pid) => {
          try {
            const p = (await fetchProduct(pid)) as Product;
            return [pid, p] as const;
          } catch {
            return [pid, null] as const;
          } finally {
            inFlightProducts.current.delete(pid);
          }
        })
      );

      const patch: Record<number, Product> = {};
      for (const [pid, p] of results) if (p) patch[pid] = p;

      if (Object.keys(patch).length) {
        setProductCache((prev) => ({ ...prev, ...patch }));
      }
    },
    [productCache]
  );

  const loadDetails = useCallback(
    async (orderId: number) => {
      if (detailsById[orderId]) return;
      if (inFlightDetails.current.has(orderId)) return;

      inFlightDetails.current.add(orderId);
      setLoadingDetailsById((prev) => ({ ...prev, [orderId]: true }));

      try {
        const d = (await myOrder(orderId)) as OrderDetails;

        setDetailsById((prev) => ({ ...prev, [orderId]: d }));

        const productIds = (d.items ?? []).map((it) => Number(it.product_id));
        await ensureProductsCached(productIds);
      } catch {
        message.error("Unable to load order details");
      } finally {
        inFlightDetails.current.delete(orderId);
        setLoadingDetailsById((prev) => ({ ...prev, [orderId]: false }));
      }
    },
    [detailsById, ensureProductsCached]
  );

  const itemsColumns: ColumnsType<OrderItem> = useMemo(
    () => [
      {
        title: "Product",
        dataIndex: "product_id",
        render: (_: unknown, it: OrderItem) => {
          const p = productCache[it.product_id];
          const img =
            p?.image_url && !/^https?:/i.test(p.image_url)
              ? apiOrigin + p.image_url
              : p?.image_url;

          return (
            <Space>
              {img ? (
                <img
                  src={img}
                  alt={p?.name || String(it.product_id)}
                  style={{
                    width: 56,
                    height: 42,
                    objectFit: "cover",
                    borderRadius: 4,
                  }}
                />
              ) : null}

              <div>
                <Link to={`/product/${it.product_id}`}>
                  {p?.name || `#${it.product_id}`}
                </Link>
              </div>
            </Space>
          );
        },
      },
      { title: "Qty", dataIndex: "quantity" },
      { title: "Unit Price", dataIndex: "unit_price" },
      { title: "Line Total", dataIndex: "line_total" },
    ],
    [apiOrigin, productCache]
  );

  const expandedRowRender = useCallback(
    (row: OrderRow) => {
      const d = detailsById[row.id];
      const createdAt = d?.created_at || row.created_at;
      const created = createdAt ? new Date(createdAt) : null;

      const etaDays =
        d?.estimated_production_time_days ??
        row.estimated_production_time_days ??
        0;

      const readyBy = (() => {
        if (!created) return null;
        const r = new Date(created);
        r.setDate(r.getDate() + Number(etaDays || 0));
        return r;
      })();

      return (
        <Space direction="vertical" style={{ width: "100%" }} size="large">
          <Descriptions bordered size="small" column={3}>
            <Descriptions.Item label="Order ID">{row.id}</Descriptions.Item>
            <Descriptions.Item label="Status">{row.status}</Descriptions.Item>
            <Descriptions.Item label="Payment">
              {row.payment_status}
            </Descriptions.Item>

            <Descriptions.Item label="Created">
              {created ? created.toLocaleString() : "-"}
            </Descriptions.Item>

            <Descriptions.Item label="ETA (days)">{etaDays}</Descriptions.Item>

            <Descriptions.Item label="Ready By">
              {readyBy ? readyBy.toLocaleDateString() : "-"}
            </Descriptions.Item>

            <Descriptions.Item label="Total">
              {d?.total_price ?? row.total_price}
            </Descriptions.Item>

            {d?.address && (
              <Descriptions.Item label="Address" span={2}>
                {d.address}
              </Descriptions.Item>
            )}

            {d?.phone && (
              <Descriptions.Item label="Phone">{d.phone}</Descriptions.Item>
            )}
          </Descriptions>

          <Table<OrderItem>
            rowKey={(r) => String(r.id ?? r.product_id)}
            size="small"
            pagination={false}
            loading={!!loadingDetailsById[row.id] && !detailsById[row.id]}
            dataSource={d?.items ?? []}
            columns={itemsColumns}
          />
        </Space>
      );
    },
    [detailsById, itemsColumns, loadingDetailsById]
  );

  const onExpand = useCallback(
    (expanded: boolean, record: OrderRow) => {
      if (expanded) loadDetails(record.id);
      setExpandedRowKeys((prev) => {
        const key = record.id as unknown as React.Key;
        if (expanded) return [key];
        return prev.filter((k) => k !== key);
      });
    },
    [loadDetails]
  );

  return (
    <div>
      <Typography.Title level={2}>{t("orders.title")}</Typography.Title>

      <Table<OrderRow>
        rowKey="id"
        dataSource={orders}
        columns={orderColumns}
        expandable={{
          onExpand,
          expandedRowRender,
          expandedRowKeys,
        }}
      />
    </div>
  );
}
