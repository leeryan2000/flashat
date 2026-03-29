const API_BASE = import.meta.env.VITE_API_URL + "/api";

// ***** Optimize the api calls to not use application/json all the time and 
export async function api<T>(path: string, init: RequestInit = {}): Promise<T> {
  const res = await fetch(API_BASE + path, {
    credentials: "include", // 👈 send cookies with the request
    headers: { "Content-Type": "application/json", ...(init.headers || {}) },
    ...init,
  });

  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error(body.error || `HTTP ${res.status}`);
  }

  const text = await res.text();
  return (text ? JSON.parse(text) : undefined) as T;
}