import { API } from "../registerform";

export async function saveProfile(payload: Record<string, unknown>, route: string) {
  if (Object.keys(payload).length === 0) return; // nothing to update

  const token = localStorage.getItem("token");
  const res = await fetch(`${API}${route}`, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
    body: JSON.stringify(payload),
  });

  if (!res.ok) {
    const text = await res.text().catch(() => "");
    throw new Error(`HTTP ${res.status}${text ? `: ${text}` : ""}`);
  }

  
  const ct = res.headers.get("content-type") || "";
  if (ct.includes("application/json")) {
    return res.json();
  }
  return null;
}
