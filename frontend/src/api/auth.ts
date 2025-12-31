import { api, setAuthToken } from "./client";

export const register = async (data: any) => {
  const res = await api.post("/auth/register", data);
  return res.data;
};

export const login = async (email: string, password: string) => {
  const res = await api.post("/auth/login", { email, password });
  setAuthToken(res.data.token);
  return res.data;
};

export const me = async () => (await api.get("/user/me")).data;
