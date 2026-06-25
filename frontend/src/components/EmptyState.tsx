import type { ReactNode } from "react";

interface EmptyStateProps {
  title: string;
  description: string;
  eyebrow?: string;
  icon?: ReactNode;
  action?: ReactNode;
  tone?: "neutral" | "error";
  className?: string;
}

function EmptyState({
  title,
  description,
  eyebrow = "Empty State",
  icon,
  action,
  tone = "neutral",
  className = "",
}: EmptyStateProps) {
  const classNames = [
    "empty-state",
    tone === "error" ? "is-error" : "",
    className,
  ]
    .filter(Boolean)
    .join(" ");

  return (
    <section
      className={classNames}
      role={tone === "error" ? "alert" : "status"}
      aria-live={tone === "error" ? "assertive" : "polite"}
    >
      <div className="empty-state-inner">
        <span className="empty-state-eyebrow">{eyebrow}</span>

        {icon && (
          <span className="empty-state-icon" aria-hidden="true">
            {icon}
          </span>
        )}

        <h2>{title}</h2>
        <p>{description}</p>

        {action && <div className="empty-state-action">{action}</div>}
      </div>
    </section>
  );
}

export default EmptyState;
