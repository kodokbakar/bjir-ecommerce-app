import { Link } from "react-router-dom";
import type { LucideIcon } from "lucide-react";

interface NavItemProps {
  to: string;
  label: string;
  Icon: LucideIcon;
  isActive: boolean;
  isCollapsed?: boolean;
  onClick?: () => void;
}

function NavItem({
  to,
  label,
  Icon,
  isActive,
  isCollapsed = false,
  onClick,
}: NavItemProps) {
  const linkClassName = [
    "group relative flex min-h-11 items-center gap-3 rounded-xl border-l-4 px-3 text-sm font-bold transition-all duration-200",
    "focus-visible:outline focus-visible:outline-3 focus-visible:outline-offset-2 focus-visible:outline-[var(--color-brutal-hot)]",
    isCollapsed ? "justify-center" : "justify-start",
    isActive
      ? "border-l-[var(--color-brutal-hot)] bg-[var(--color-secondary)] text-[var(--color-text-dark)] shadow-[3px_3px_0_var(--color-brutal-ink)]"
      : "border-l-transparent text-[var(--color-text-muted)] hover:bg-[var(--color-brutal-accent)] hover:text-[var(--color-brutal-ink)]",
  ]
    .filter(Boolean)
    .join(" ");

  return (
    <Link
      className={linkClassName}
      to={to}
      title={isCollapsed ? label : undefined}
      aria-current={isActive ? "page" : undefined}
      onClick={onClick}
    >
      <Icon className="h-5 w-5 shrink-0" aria-hidden="true" />

      {!isCollapsed && <span className="truncate">{label}</span>}
    </Link>
  );
}

export default NavItem;
