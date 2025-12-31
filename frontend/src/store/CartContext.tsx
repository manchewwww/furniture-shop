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

export type CartItem = {
  product: any;
  quantity: number;
  options: { id: number; type: string }[];
};

type CartCtxType = {
  items: CartItem[];
  add: (item: CartItem) => void;
  remove: (productId: number) => void;
  increment: (productId: number) => void;
  decrement: (productId: number) => void;
  clear: () => void;
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
        const serverItems: CartItem[] = (server.items || []).map((it: any) => ({
          product: { id: it.product_id },
          quantity: it.quantity,
          options: JSON.parse(it.selected_options_json || "[]"),
        }));
        if (items.length && serverItems.length === 0) {
          const payload = items.map((it) => ({
            product_id: it.product.id,
            quantity: it.quantity,
            options: it.options,
          }));
          await apiReplace(payload);
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
          const mapped: CartItem[] = (s.items || []).map((it: any) => ({
            product: { id: it.product_id },
            quantity: it.quantity,
            options: JSON.parse(it.selected_options_json || "[]"),
          }));
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
          const mapped: CartItem[] = (ref.items || []).map((it: any) => ({
            product: { id: it.product_id },
            quantity: it.quantity,
            options: JSON.parse(it.selected_options_json || "[]"),
          }));
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
          const mapped: CartItem[] = (ref.items || []).map((it: any) => ({
            product: { id: it.product_id },
            quantity: it.quantity,
            options: JSON.parse(it.selected_options_json || "[]"),
          }));
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
          const mapped: CartItem[] = (ref.items || []).map((it: any) => ({
            product: { id: it.product_id },
            quantity: it.quantity,
            options: JSON.parse(it.selected_options_json || "[]"),
          }));
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
