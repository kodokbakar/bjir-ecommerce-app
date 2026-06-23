import { APP_VERSION } from "./navigation";

function Footer() {
  return (
    <footer className="flex flex-col gap-3 border-t border-[var(--color-border)] bg-white px-4 py-5 text-xs font-bold text-[var(--color-text-muted)] sm:px-6 lg:flex-row lg:items-center lg:justify-between lg:px-8">
      <span className="text-center lg:text-left">
        &copy; {new Date().getFullYear()} Bjir E-commerce. Hak cipta dilindungi.
      </span>

      <nav
        className="flex flex-wrap items-center justify-center gap-x-4 gap-y-2"
        aria-label="Footer navigation"
      >
        <a
          className="text-[var(--color-text-muted)] no-underline hover:text-[var(--color-primary-dark)]"
          href="#about"
        >
          Tentang Kami
        </a>
        <a
          className="text-[var(--color-text-muted)] no-underline hover:text-[var(--color-primary-dark)]"
          href="#privacy"
        >
          Kebijakan Privasi
        </a>
        <a
          className="text-[var(--color-text-muted)] no-underline hover:text-[var(--color-primary-dark)]"
          href="#contact"
        >
          Hubungi Kami
        </a>
        <span>Versi {APP_VERSION}</span>
      </nav>
    </footer>
  );
}

export default Footer;
