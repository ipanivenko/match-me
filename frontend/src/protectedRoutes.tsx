import type { JSX } from "react";
import { Navigate, useLocation } from "react-router-dom";

const hasToken = () => !!localStorage.getItem("token");

export function RequireAuth({ children }: { children: JSX.Element }) {
  const loc = useLocation();
  // If not logged in, go to /login and remember where the user wanted to go
  return hasToken() ? children : <Navigate to="/login" replace state={{ from: loc }} />;
}

export function RequireGuest({ children }: { children: JSX.Element }) {
  // If already logged in, no need to see login/register again
  return hasToken() ? <Navigate to="/profile" replace /> : children;
}
