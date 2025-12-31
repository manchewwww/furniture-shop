import { api } from "./client";

export const getCart = async () => (await api.get("/user/cart")).data;
export const replaceCart = async (items: any[]) =>
  (await api.put("/user/cart", { items })).data;
export const addItem = async (item: any) =>
  (await api.post("/user/cart/items", item)).data;
export const updateItem = async (id: number, item: any) =>
  (await api.patch(`/user/cart/items/${id}`, item)).data;
export const removeItem = async (id: number) =>
  (await api.delete(`/user/cart/items/${id}`)).data;
export const clearCart = async () => (await api.delete(`/user/cart`)).data;
