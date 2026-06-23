interface CartBadgeProps {
  count: number;
}

function getBadgeLabel(count: number): string {
  if (count > 99) {
    return "99+";
  }

  return String(count);
}

function CartBadge({ count }: CartBadgeProps) {
  if (count <= 0) {
    return null;
  }

  return (
    <span
      className="absolute -right-1 -top-1 grid min-h-5 min-w-5 place-items-center rounded-full border-2 border-white bg-[var(--color-stock-out)] px-1 text-[10px] font-black leading-none text-white shadow-[2px_2px_0_var(--color-brutal-ink)]"
      aria-label={`${count} cart item${count === 1 ? "" : "s"}`}
    >
      {getBadgeLabel(count)}
    </span>
  );
}

export default CartBadge;
