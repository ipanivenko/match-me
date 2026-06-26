import { API } from "../registerform";

export function buildAvatarUrl(
  cloudName: string,
  publicId: string,
  version?: number,
  size = 256
) {
  if (!cloudName || !publicId) return "";

  // Encode each path segment, keep slashes
  const encodedId = publicId.split("/").map(encodeURIComponent).join("/");

  const v = typeof version === "number" && version > 0 ? `v${version}/` : "";
  const transform = `c_fill,w_${size},h_${size},g_face,f_auto,q_auto,dpr_auto`;

  return `https://res.cloudinary.com/${cloudName}/image/upload/${transform}/${v}${encodedId}`;
}

export async function uploadAvatar(file: File) {
  const token = localStorage.getItem("token") ?? "";

  // 1) get signed params from your backend (authorized)
  const s = await fetch(`${API}/me/cloudinary-sign`, {
    method: "GET",
    headers: {
      Accept: "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
  }).then((r) => {
    if (!r.ok) throw new Error(`Sign failed: ${r.status}`);
    return r.json();
  });

  // 2) direct upload to Cloudinary (no auth header)
  const fd = new FormData();
  fd.append("file", file);
  fd.append("api_key", String(s.api_key));
  fd.append("timestamp", String(s.timestamp));   // ensure string
  fd.append("signature", String(s.signature));
  fd.append("folder", String(s.folder));
  fd.append("public_id", String(s.public_id));   // "avatar"
  fd.append("overwrite", String(s.overwrite));   // "true"

  const upRes = await fetch(
    `https://api.cloudinary.com/v1_1/${s.cloud_name}/image/upload`,
    { method: "POST", body: fd }
  );
  const upText = await upRes.text();
  if (!upRes.ok) throw new Error(upText || "Cloudinary upload failed");
  const uploaded = JSON.parse(upText);

  // 3) persist stable handle on your backend (authorized)
  const saveRes = await fetch(`${API}/me/photo`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
    body: JSON.stringify({
      public_id: uploaded.public_id,
      version: uploaded.version,
    }),
  });
  const saveText = await saveRes.text();
  if (!saveRes.ok) throw new Error(`Persist failed: ${saveRes.status} — ${saveText}`);

  // 4) return a ready-to-use URL + ids
  return {
    public_id: uploaded.public_id as string,
    version: uploaded.version as number,
    url: buildAvatarUrl(s.cloud_name, uploaded.public_id, uploaded.version),
    cloud_name: s.cloud_name as string,
  };
}

export async function deleteAvatar() {
  const token = localStorage.getItem("token") ?? "";
  const res = await fetch(`${API}/me/photo`, {
    method: "DELETE",
    headers: {
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
  });
  if (!res.ok) throw new Error(`Delete failed: ${res.status}`);
}
