import { useNavigate } from "react-router-dom";
import { clearToken } from "./token";

export function useLogout() {
  const navigate = useNavigate();

  return (redirectTo = "/login") => {
    clearToken();
    navigate(redirectTo, { replace: true });
  };
}
