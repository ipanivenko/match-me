const API = import.meta.env.VITE_API_BASE_URL;

export async function acceptOrRejectConnection(
  connectionId: string,
  action: "accept" | "reject"
) {
  const token = localStorage.getItem("token");

  const res = await fetch(`${API}/connections/${connectionId}/action`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ action }),
  });

  if (!res.ok) {
    const text = await res.text();
    console.error("[connection action error]", res.status, text);
    throw new Error(text || `Failed to ${action} connection`);
  }

  return res.json();
}