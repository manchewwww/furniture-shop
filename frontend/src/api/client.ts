import axios from "axios";

const API_URL = "http://localhost:8080/api";

export const api = axios.create({
  baseURL: API_URL,
});

export const setAuthToken = (token?: string) => {
  if (token) {
    api.defaults.headers.common["Authorization"] = `Bearer ${token}`;
    localStorage.setItem("token", token);
  } else {
    delete api.defaults.headers.common["Authorization"];
    localStorage.removeItem("token");
  }
};

const existing = localStorage.getItem("token");
if (existing) setAuthToken(existing);

export const getApiOrigin = () => {
  try {
    return new URL(API_URL).origin;
  } catch {
    return window.location.origin;
  }
};
