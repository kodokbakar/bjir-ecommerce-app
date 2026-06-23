import {
  BadgeCheck,
  Grid2X2,
  ShieldCheck,
  ShoppingBag,
  Star,
  UserPlus,
  Zap,
  type LucideIcon,
} from "lucide-react";

export type BrandPanelVariant = "login" | "register";

interface BrandPanelProps {
  variant: BrandPanelVariant;
}

interface BrandCopy {
  headline: string;
  description: string;
  pills: string[];
  icons: LucideIcon[];
}

const BRAND_COPY: Record<BrandPanelVariant, BrandCopy> = {
  login: {
    headline: "Masuk ke TOKO belanja paling berisik.",
    description:
      "Cari produk, cek stok, dan lanjutkan belanja tanpa drama. Katalog dibuat cepat, tajam, dan gampang dipakai.",
    pills: ["Checkout cepat", "Produk real", "Aman dipakai"],
    icons: [ShoppingBag, ShieldCheck, Zap],
  },
  register: {
    headline: "Buat akun, mulai berburu barang.",
    description:
      "Gabung sebagai member dan simpan akses ke katalog produk, order, dan promo yang siap kamu ambil.",
    pills: ["Gratis daftar", "Member ready", "Belanja aman"],
    icons: [UserPlus, BadgeCheck, Star],
  },
};

function BrandPanel({ variant }: BrandPanelProps) {
  const copy = BRAND_COPY[variant];

  return (
    <aside className="relative hidden min-h-[620px] overflow-hidden bg-[var(--color-primary)] p-9 text-[var(--color-brutal-paper)] md:flex md:flex-col md:justify-between">
      <span className="absolute -right-12 -top-12 h-44 w-44 rounded-full bg-[var(--color-brutal-accent)] opacity-80" />
      <span className="absolute -bottom-16 -left-14 h-48 w-48 rounded-full bg-[var(--color-primary-dark)] opacity-60" />
      <span className="absolute bottom-28 right-7 h-24 w-24 rounded-full bg-[var(--color-brutal-blue)] opacity-70" />

      <div className="relative z-10 flex items-center gap-3">
        <span className="grid h-11 w-11 place-items-center border-2 border-[var(--color-brutal-ink)] bg-[var(--color-brutal-paper)] text-[var(--color-primary)] shadow-[4px_4px_0_var(--color-brutal-ink)]">
          <Grid2X2 className="h-6 w-6" aria-hidden="true" />
        </span>
        <span className="text-lg font-black tracking-tight">Bjir E-Commerce</span>
      </div>

      <div className="relative z-10">
        <div className="mb-6 flex gap-3">
          {copy.icons.map((Icon, index) => (
            <span
              className={[
                "grid h-14 w-14 place-items-center border-2 border-[var(--color-brutal-ink)] shadow-[4px_4px_0_var(--color-brutal-ink)]",
                index === 2
                  ? "bg-[var(--color-brutal-blue)] text-[var(--color-brutal-ink)]"
                  : "bg-white/15 text-white",
              ].join(" ")}
              key={Icon.displayName ?? Icon.name}
            >
              <Icon className="h-6 w-6" aria-hidden="true" />
            </span>
          ))}
        </div>

        <h1 className="m-0 max-w-sm text-5xl font-black uppercase leading-[0.9] tracking-[-0.08em]">
          {copy.headline}
        </h1>

        <p className="mt-5 max-w-md text-sm font-bold leading-7 text-[rgba(255,248,232,0.86)]">
          {copy.description}
        </p>
      </div>

      <div className="relative z-10 flex flex-wrap gap-2">
        {copy.pills.map((label) => (
          <span
            className="border-2 border-white/25 bg-black/20 px-3 py-1 text-xs font-black uppercase tracking-[0.12em]"
            key={label}
          >
            {label}
          </span>
        ))}
      </div>
    </aside>
  );
}

export default BrandPanel;
