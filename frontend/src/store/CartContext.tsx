import React, {
  createContext,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import { useAuth } from "./AuthContext";
import {
  addItem as apiAddItem,
  clearCart as apiClear,
  getCart as apiGet,
  removeItem as apiRemove,
  replaceCart as apiReplace,
  updateItem as apiUpdate,
} from "../api/cart";
import { fetchProduct } from "../api/catalog";

export type CartItem = {
  product: any;
  quantity: number;
  options: { id: number; type: string }[];
};

type CartCtxType = {
  items: CartItem[];
  add: (item: CartItem) => Promise<void>;
  remove: (productId: number) => Promise<void>;
  increment: (productId: number) => Promise<void>;
  decrement: (productId: number) => Promise<void>;
  clear: () => Promise<void>;
  clearLocal: () => Promise<void>;
};

const CartCtx = createContext<CartCtxType | undefined>(undefined);

export const CartProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const { isAuthenticated } = useAuth();
  const [items, setItems] = useState<CartItem[]>(() => {
    const saved = localStorage.getItem("cart");
    return saved ? JSON.parse(saved) : [];
  });
  useEffect(() => {
    localStorage.setItem("cart", JSON.stringify(items));
  }, [items]);

  useEffect(() => {
    const sync = async () => {
      if (!isAuthenticated) return;
      try {
        const server = await apiGet();
        const serverItemsRaw: any[] = server.items || [];
        const serverItems: CartItem[] = await Promise.all(
          serverItemsRaw.map(async (it: any) => ({
            product: await fetchProduct(it.product_id),
            quantity: it.quantity,
            options: JSON.parse(it.selected_options_json || "[]"),
          }))
        );

        if (items.length) {
          const keyOf = (ci: CartItem) =>
            `${ci.product.id}:${JSON.stringify(
              (ci.options || []).slice().sort((a, b) => a.id - b.id)
            )}`;
          const map = new Map<string, CartItem>();
          for (const it of serverItems) map.set(keyOf(it), { ...it });
          for (const it of items) {
            const k = keyOf(it);
            if (map.has(k)) {
              map.get(k)!.quantity += it.quantity;
            } else {
              map.set(k, { ...it });
            }
          }
          const merged = Array.from(map.values());
          // Persist merged to backend
          const payload = merged.map((it) => ({
            product_id: it.product.id,
            quantity: it.quantity,
            options: it.options,
          }));
          await apiReplace(payload);
          const refreshed = await apiGet();
          const mergedHydrated: CartItem[] = await Promise.all(
            (refreshed.items || []).map(async (it: any) => ({
              product: await fetchProduct(it.product_id),
              quantity: it.quantity,
              options: JSON.parse(it.selected_options_json || "[]"),
            }))
          );
          setItems(mergedHydrated);
        } else {
          setItems(serverItems);
        }
      } catch {}
    };
    sync();
  }, [isAuthenticated]);
  const value = useMemo(
    () => ({
      items,
      add: async (item: CartItem) => {
        if (isAuthenticated) {
          await apiAddItem({
            product_id: item.product.id,
            quantity: item.quantity,
            options: item.options,
          });
          const s = await apiGet();
          const mapped: CartItem[] = await Promise.all(
            (s.items || []).map(async (it: any) => ({
              product: await fetchProduct(it.product_id),
              quantity: it.quantity,
              options: JSON.parse(it.selected_options_json || "[]"),
            }))
          );
          setItems(mapped);
        } else {
          setItems((prev) => {
            const existing = prev.find((p) => p.product.id === item.product.id);
            if (existing) {
              return prev.map((p) =>
                p.product.id === item.product.id
                  ? { ...p, quantity: p.quantity + item.quantity }
                  : p
              );
            }
            return [...prev, item];
          });
        }
      },
      remove: async (id: number) => {
        if (isAuthenticated) {
          const s = await apiGet();
          const found = (s.items || []).find((i: any) => i.product_id === id);
          if (found) await apiRemove(found.id);
          const ref = await apiGet();
          const mapped: CartItem[] = await Promise.all(
            (ref.items || []).map(async (it: any) => ({
              product: await fetchProduct(it.product_id),
              quantity: it.quantity,
              options: JSON.parse(it.selected_options_json || "[]"),
            }))
          );
          setItems(mapped);
        } else {
          setItems((prev) => prev.filter((p) => p.product.id !== id));
        }
      },
      increment: async (id: number) => {
        if (isAuthenticated) {
          const s = await apiGet();
          const found = (s.items || []).find((i: any) => i.product_id === id);
          if (found)
            await apiUpdate(found.id, {
              quantity: found.quantity + 1,
              options: JSON.parse(found.selected_options_json || "[]"),
            });
          const ref = await apiGet();
          const mapped: CartItem[] = await Promise.all(
            (ref.items || []).map(async (it: any) => ({
              product: await fetchProduct(it.product_id),
              quantity: it.quantity,
              options: JSON.parse(it.selected_options_json || "[]"),
            }))
          );
          setItems(mapped);
        } else {
          setItems((prev) =>
            prev.map((p) =>
              p.product.id === id ? { ...p, quantity: p.quantity + 1 } : p
            )
          );
        }
      },
      decrement: async (id: number) => {
        if (isAuthenticated) {
          const s = await apiGet();
          const found = (s.items || []).find((i: any) => i.product_id === id);
          if (found) {
            const next = found.quantity - 1;
            if (next <= 0) await apiRemove(found.id);
            else
              await apiUpdate(found.id, {
                quantity: next,
                options: JSON.parse(found.selected_options_json || "[]"),
              });
          }
          const ref = await apiGet();
          const mapped: CartItem[] = await Promise.all(
            (ref.items || []).map(async (it: any) => ({
              product: await fetchProduct(it.product_id),
              quantity: it.quantity,
              options: JSON.parse(it.selected_options_json || "[]"),
            }))
          );
          setItems(mapped);
        } else {
          setItems((prev) =>
            prev.flatMap((p) => {
              if (p.product.id !== id) return [p];
              const nextQty = p.quantity - 1;
              return nextQty <= 0 ? [] : [{ ...p, quantity: nextQty }];
            })
          );
        }
      },
      clear: async () => {
        if (isAuthenticated) await apiClear();
        setItems([]);
      },
      clearLocal: async () => {
        setItems([]);
      },
    }),
    [items, isAuthenticated]
  );
  return <CartCtx.Provider value={value}>{children}</CartCtx.Provider>;
};

export const useCart = () => {
  const ctx = useContext(CartCtx);
  if (!ctx) throw new Error("CartProvider is missing");
  return ctx;
};
