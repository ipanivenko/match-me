import { API } from "../registerform";

export async function disconnectUser(
  userId: string
) {
  const token = localStorage.getItem("token");

  const res = await fetch(`${API}/reactions/disconnect`, {
    method: "POST",
   headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
    body: JSON.stringify({ target_user_id: userId }),
  });

  if (!res.ok) {
    const text = await res.text();
    console.error("[reaction error]", res.status, res.statusText, "→", text);
    throw new Error(text || `HTTP ${res.status} ${res.statusText}`);
  }

  return res.json() as Promise<{ status: string }>;
}
