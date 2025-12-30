import { api } from "./client";

export const createOrder = (payload: any) =>
  api.post("/orders", payload).then((r) => r.data);
export const myOrders = () => api.get("/user/orders").then((r) => r.data);
export const myOrder = (id: number) =>
  api.get(`/user/orders/${id}`).then((r) => r.data);
export const payOrder = (id: number) =>
  api.post(`/user/orders/${id}/pay`).then((r) => r.data);
