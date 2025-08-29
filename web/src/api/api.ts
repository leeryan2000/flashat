const API_BASE = import.meta.env.VITE_API_URL + "/api";

export async function api<T>(path: string, init: RequestInit = {}): Promise<T> {
  const res = await fetch(API_BASE + path, {
    credentials: "include", // 👈 send cookies with the request
    headers: { "Content-Type": "application/json", ...(init.headers || {}) },
    ...init,
  });

  if (res.status === 401) throw new Error("unauthorized");
  if (!res.ok) throw new Error(`HTTP ${res.status}`);

  return res.status === 204 ? (undefined as T) : res.json();
}