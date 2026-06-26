import { useEffect, useState } from "react";
import { useLogout } from "../auth/useLogout";

const API = import.meta.env.VITE_API_BASE_URL;

export default function UserHeader() {
  const [email, setEmail] = useState<string>("");
  const logout = useLogout();

  useEffect(() => {
    const fetchEmail = async () => {
      try {
        const token = localStorage.getItem("token");
        const res = await fetch(`${API}/me/email`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        
        if (res.ok) {
          const data = await res.json();
          setEmail(data.email);
        }
      } catch (err) {

      }
    };

    fetchEmail();
  }, []);

  return (
    <div className="user-header">
      <div className="header-right">
        {email && <div className="user-email">👤 {email}</div>}
        <button className="button is-dark logout" onClick={() => logout()}>
          Log out
        </button>
      </div>
    </div>
  );
}