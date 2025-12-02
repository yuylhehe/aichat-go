export const api = {
  base: '/api/v1',
  token: null,
  setToken(t) {
    this.token = t;
    if (t) {
      localStorage.setItem('token', t);
    } else {
      localStorage.removeItem('token');
    }
  },
  getToken() {
    return this.token ?? localStorage.getItem('token');
  },
  headers() {
    const h = { 'Content-Type': 'application/json' };
    const t = this.getToken();
    if (t) h['Authorization'] = `Bearer ${t}`;
    return h;
  },
  async get(path) {
    const res = await fetch(this.base + path, { headers: this.headers() });
    return handle(res);
  },
  async post(path, body) {
    const res = await fetch(this.base + path, {
      method: 'POST',
      headers: this.headers(),
      body: JSON.stringify(body),
    });
    return handle(res);
  },
  async patch(path, body) {
    const res = await fetch(this.base + path, {
      method: 'PATCH',
      headers: this.headers(),
      body: JSON.stringify(body),
    });
    return handle(res);
  },
  async del(path) {
    const res = await fetch(this.base + path, {
      method: 'DELETE',
      headers: this.headers(),
    });
    return handle(res);
  },
};

async function handle(res) {
  const data = await res.json().catch(() => ({}));
  if (!res.ok || data.success === false)
    throw new Error(data.message?.message ?? data.message ?? res.statusText);
  return data.data ?? data;
}
