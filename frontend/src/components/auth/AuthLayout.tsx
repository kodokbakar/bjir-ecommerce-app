import { Grid2X2 } from "lucide-react";
import type { ReactNode } from "react";

import BrandPanel, { type BrandPanelVariant } from "./BrandPanel";

interface AuthLayoutProps {
  variant: BrandPanelVariant;
  children: ReactNode;
}

function AuthLayout({ variant, children }: AuthLayoutProps) {
  return (
    <main className="grid min-h-screen place-items-center bg-[var(--color-brutal-paper)] px-4 py-6 sm:px-6">
      <section className="grid w-full max-w-5xl overflow-hidden border-4 border-[var(--color-brutal-ink)] bg-[var(--color-brutal-surface)] shadow-[8px_8px_0_var(--color-brutal-ink)] md:grid-cols-[minmax(320px,0.9fr)_minmax(360px,1fr)]">
        <BrandPanel variant={variant} />

        <div className="flex min-h-[100svh] flex-col justify-center bg-[var(--color-brutal-surface)] p-6 sm:min-h-[620px] sm:p-10">
          <div className="mb-8 flex items-center gap-3 md:hidden">
            <span className="grid h-10 w-10 place-items-center border-2 border-[var(--color-brutal-ink)] bg-[var(--color-primary)] text-white shadow-[3px_3px_0_var(--color-brutal-ink)]">
              <Grid2X2 className="h-5 w-5" aria-hidden="true" />
            </span>
            <span className="text-base font-black text-[var(--color-brutal-ink)]">
              Bjir E-Commerce
            </span>
          </div>

          {children}
        </div>
      </section>
    </main>
  );
}

export default AuthLayout;
