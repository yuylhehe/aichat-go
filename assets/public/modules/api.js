import { API_BASE } from "./config.js";
import { Store } from "./store.js";

export const API = {
  headers() {
    const h = { "Content-Type": "application/json" };
    if (Store.token) h["Authorization"] = `Bearer ${Store.token}`;
    return h;
  },

  async request(method, path, body = null) {
    const opts = { method, headers: this.headers() };
    if (body) opts.body = JSON.stringify(body);

    try {
      const res = await fetch(API_BASE + path, opts);

      if (res.status === 401) {
        // Dispatch event for Auth module to handle
        window.dispatchEvent(new CustomEvent("auth:unauthorized"));
        throw new Error("Unauthorized");
      }

      const data = await res.json();

      if (!res.ok) {
        throw new Error(
          data.message?.message || data.message || res.statusText
        );
      }

      return data.data ?? data;
    } catch (err) {
      console.error("API Error:", err);
      throw err;
    }
  },

  get: (path) => API.request("GET", path),
  post: (path, body) => API.request("POST", path, body),
  put: (path, body) => API.request("PUT", path, body),
  del: (path) => API.request("DELETE", path),
};
