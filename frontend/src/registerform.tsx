import { Link } from "react-router-dom";
import { useNavigate } from "react-router-dom";
import "bulma/css/bulma.min.css";
import "./styles/reg_login.css";
import axios from "axios";

export const API = import.meta.env.VITE_API_BASE_URL;

type RegisterResponse = {
  id: string;
  email: string;
  createdAt?: string;
};

export type ErrorResponse = {
  message: string;
};

export default function RegisterForm() {
  const navigate = useNavigate();
  const handleSubmit: React.FormEventHandler<HTMLFormElement> = async (e) => {
    e.preventDefault();

    const fd = new FormData(e.currentTarget);
    const email = fd.get("email")?.toString().trim() ?? "";
    const password = fd.get("password")?.toString() ?? "";
    const confirm = fd.get("confirm")?.toString() ?? "";

    if (!email || !password) {
      alert("Email and password are required");
      return;
    }

    if (password !== confirm) {
      alert("Passwords don't match ❌");
      return;
    }
    try {
      const res = await axios.post<RegisterResponse>(`${API}/users/register`, {
        email,
        password,
      });
      alert("Account created ✅");
      navigate("/login", { replace:true });
      console.log(res.data);
    } catch (err) {
      if (axios.isAxiosError(err)) {
        const data = err.response?.data as ErrorResponse | undefined;
        if (data?.message) {
          alert(`${data.message}`);
        } else {
          alert(`Request failed ❌ (${err.response?.status ?? "no status"})`);
        }
      } else {
        alert("Unexpected error ❌");
        console.error(err);
      }
    }
  };

  return (
    <div className="page-container has-background-light">
      <form className="register-form" onSubmit={handleSubmit}>
        <p className="has-text-right">
          Already have an account? <Link to="/login">Login</Link>
        </p>
        <h2 className="title is-4 has-text-centered">
          Welcome to the Match-me-children app!
        </h2>
        <br />
        <h2 className="subtitle is-5">Create your account:</h2>
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
            placeholder="At least 6 characters, one letter"
            required
            autoComplete="new-password"
            minLength={6}
            className="input"
          />
        </label>

        <label>
          Confirm password
          <input
            type="password"
            name="confirm"
            placeholder="Repeat your password"
            required
            autoComplete="new-password"
            minLength={6}
            className="input"
          />
        </label>

        <button type="submit" className="button is-primary is-fullwidth">
          Create account
        </button>
      </form>
    </div>
  );
}
