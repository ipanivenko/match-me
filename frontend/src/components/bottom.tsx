import "../styles/bottomPanel.css";
import { useNavigate, useLocation } from "react-router-dom";

const IconUser = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true">
    <path d="M12 12a5 5 0 1 0-5-5 5 5 0 0 0 5 5Zm0 2c-5 0-9 2.5-9 5.5A1.5 1.5 0 0 0 4.5 21h15A1.5 1.5 0 0 0 21 19.5C21 16.5 17 14 12 14Z" />
  </svg>
);

const IconChild = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true">
    <path d="M12 4a3 3 0 1 1-3 3 3 3 0 0 1 3-3Zm7 13a5 5 0 0 0-5-5H10a5 5 0 0 0-5 5 2 2 0 0 0 2 2h10a2 2 0 0 0 2-2Z" />
  </svg>
);

const IconStar = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true">
    <path d="m12 2 3 6 7 .9-5 4.8 1.4 6.9L12 17l-6.4 3.6L7 13.7 2 8.9 9 8z" />
  </svg>
);

const IconConnections = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true">
    <path
      d="M12 2a3 3 0 1 1-3 3 3 3 0 0 1 3-3zm7 7a3 3 0 1 1-3 3 3 3 0 0 1 3-3zm-14 0a3 3 0 1 1-3 3 3 3 0 0 1 3-3zm7 7a3 3 0 1 1-3 3 3 3 0 0 1 3-3zm-5.2-1.8 3.4-2m3.6 0 3.4 2m-7 4.8L12 18l2.8 1.7"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.6"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
);

const IconChat = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true">
    <path d="M4 4h16v12H7l-3 3V4z" />
  </svg>
);

export default function BottomPanel() {
  const nav = useNavigate();
  const { pathname } = useLocation();

  const Btn = (to: string, label: string, Icon: React.FC) => {
    const isConnections =
      to === "/connections" &&
      (pathname === "/connections" || pathname === "/connections/requests");

    const isActive = isConnections || pathname === to;

    return (
      <button
        key={to}
        type="button"
        className={`button is-small is-rounded icon-btn ${
          isActive ? "is-primary" : "is-light"
        }`}
        onClick={() => nav(to)}
        aria-label={label}
        title={label}>
        <Icon />
      </button>
    );
  };

  return (
    <nav className="bottom-panel">
      {Btn("/profile", "Profile", IconUser)}
      {Btn("/child", "Child profile", IconChild)}
      {Btn("/recommendations", "Recommendations", IconStar)}
      {Btn("/connections", "Connections", IconConnections)}
      {Btn("/chats", "Chats", IconChat)}
    </nav>
  );
}