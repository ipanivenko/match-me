import { API } from "../registerform";

export async function reactToUser(
  route: string,
  userId: string,
  reaction: "like" | "dislike"
) {
  const token = localStorage.getItem("token");

  const res = await fetch(`${API}${route}/${userId}/reaction`, {
    method: "POST",
   headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
    body: JSON.stringify({ reaction }),
  });

  if (!res.ok) {
    const text = await res.text();
    console.error("[reaction error]", res.status, res.statusText, "→", text);
    throw new Error(text || `HTTP ${res.status} ${res.statusText}`);
  }

  return res.json() as Promise<{ reaction: string }>;
}
