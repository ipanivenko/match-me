import { Link } from "react-router-dom";
import "bulma/css/bulma.min.css";
import "./styles/reg_login.css";
import axios from "axios";
import type { ErrorResponse } from "./registerform";
import { API } from "./registerform";
import { useNavigate } from "react-router-dom";

type LoginResponse = {
  user_id: string;
  access_token: string;
};

export default function LoginForm() {
  const navigate = useNavigate();
  const handleSubmit: React.FormEventHandler<HTMLFormElement> = async (e) => {
    e.preventDefault();

    const fd = new FormData(e.currentTarget);
    const email = fd.get("email")?.toString().trim() ?? "";
    const password = fd.get("password")?.toString() ?? "";

    if (!email || !password) {
      alert("Email and password are required");
      return;
    }

    try {
      const res = await axios.post<LoginResponse>(`${API}/users/login`, {
        email,
        password,
      });
      alert("Loged in ✅");

    const data = res.data; 
    localStorage.setItem("token", data.access_token);
    localStorage.setItem("userId", data.user_id);
    
      setTimeout(() => {
      navigate("/profile", { replace: true });
      }, 100);
    } catch (err) {
      if (axios.isAxiosError(err)) {
        const data = err.response?.data as ErrorResponse | undefined;
        if (data?.message) {
          alert(`${data.message}`);
        } else {
          alert(`Request failed ❌ (${err.response?.status ?? "Please, try again later"})`);
        }
      } else {
        alert("Unexpected error ❌");
      }
    }
  };

  return (
    <div className="page-container has-background-light">
      <form className="login-form" onSubmit={handleSubmit}>
        <p className="has-text-right">
          <Link to="/register">Register</Link> if you do not have an account
        </p>
        <h2 className="title is-4 has-text-centered">
          Welcome to the Match-me-children app!
        </h2>
        <br />
        <h2 className="subtitle is-5">Login to your account:</h2>
        <label>
          Email
          <input
            type="email"
            name="email"
            placeholder="you@example.com"
            required
            autoComplete="email"
            className="input"
          />
        </label>

        <label>
          Password
          <input
            type="password"
            name="password"
            placeholder="enter your password"
            required
            className="input"
          />
        </label>

        <button type="submit" className="button is-primary is-fullwidth">
          Login
        </button>
      </form>
    </div>
  );
}
