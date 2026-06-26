// src/api/client.ts
import { API } from "../registerform";

export async function get<T>(path: string): Promise<T> {
  const token = localStorage.getItem("token") || "";

  const res = await fetch(`${API}${path}`, {
    headers: {
      Accept: "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
    credentials: "include",
  });

  const data = await res.json().catch(() => null);

  // Handle unauthorized
  if (res.status === 401) {
    localStorage.removeItem("token");
    window.location.replace("/login");
    throw new Error("Unauthorized");
  }

  // Handle any other errors
  if (!res.ok) {
    const message =
      (data && data.message) || `Request failed (${res.status})`;
    throw new Error(message);
  }

  // Everything OK — return parsed JSON
  return data as T;
}
