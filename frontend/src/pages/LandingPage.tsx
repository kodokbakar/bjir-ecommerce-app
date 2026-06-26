import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import {
  ArrowRight,
  BadgeCheck,
  Code2,
  Grid2X2,
  Mail,
  MapPin,
  Menu,
  Phone,
  RotateCcw,
  ShieldCheck,
  ShoppingCart,
  Sparkles,
  Truck,
  X,
} from "lucide-react";

import ProductImage from "../components/ProductImage";
import { useToast } from "../context/toast";
import { useAuth } from "../hooks/useAuth";
import {
  getContactErrorMessage,
  sendContactMessage,
} from "../services/contactService";
import { listProducts } from "../services/productService";
import type { Product } from "../types/product";
import { getDashboardPath } from "../utils/authRouting";
import { formatRupiah, getProductImage, getStockState } from "../utils/product";

const FEATURED_LIMIT = 8;

const FOOTER_QUICK_LINKS = [
  { label: "Produk", to: "/products" },
  { label: "Pesanan Saya", to: "/orders" },
  { label: "Keranjang", to: "/cart" },
];

const FOOTER_SOCIAL_LINKS = [
  {
    label: "GitHub",
    href: "https://github.com/kodokbakar/bjir-ecommerce-app",
    icon: Code2,
  },
  {
    label: "Email",
    href: "mailto:yudisbaek@gmail.com",
    icon: Mail,
  },
];

const CONTACT_EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

interface ContactFormState {
  name: string;
  email: string;
  message: string;
}

type ContactFormErrors = Partial<Record<keyof ContactFormState, string>>;

const CONTACT_FORM_INITIAL_STATE: ContactFormState = {
  name: "",
  email: "",
  message: "",
};

type FeaturedStatus = "loading" | "ready" | "error";

function LandingFeaturedSkeleton() {
  return (
    <div
      className="landing-featured-grid"
      aria-label="Loading featured products"
    >
      {Array.from({ length: FEATURED_LIMIT }, (_, index) => (
        <article className="landing-product-card is-skeleton" key={index}>
          <div className="landing-product-skeleton media" />
          <div className="landing-product-skeleton body">
            <span />
            <strong />
            <small />
          </div>
        </article>
      ))}
    </div>
  );
}

function LandingPage() {
  const [featuredProducts, setFeaturedProducts] = useState<Product[]>([]);
  const [featuredStatus, setFeaturedStatus] =
    useState<FeaturedStatus>("loading");
  const [featuredError, setFeaturedError] = useState("");
  const [isMobileNavOpen, setIsMobileNavOpen] = useState(false);
  const [contactForm, setContactForm] = useState<ContactFormState>(
    CONTACT_FORM_INITIAL_STATE,
  );
  const [contactErrors, setContactErrors] = useState<ContactFormErrors>({});
  const [isContactSubmitting, setIsContactSubmitting] = useState(false);
  const { showToast } = useToast();
  const { user, isAuthenticated, isLoading: isAuthLoading } = useAuth();

  const dashboardPath = getDashboardPath(user);
  const primaryHeroPath = isAuthenticated ? dashboardPath : "/products";
  const primaryHeroLabel = isAuthenticated ? "Buka Dashboard" : "Mulai Belanja";
  const secondaryHeroPath = isAuthenticated ? "/products" : "/register";
  const secondaryHeroLabel = isAuthenticated
    ? "Lanjut Belanja"
    : "Daftar Sekarang";
  const currentYear = new Date().getFullYear();

  useEffect(() => {
    let isActive = true;

    async function loadFeaturedProducts() {
      setFeaturedStatus("loading");
      setFeaturedError("");

      try {
        const result = await listProducts({
          page: 1,
          limit: FEATURED_LIMIT,
          sort_by: "created_at",
          sort_order: "desc",
        });

        if (isActive) {
          setFeaturedProducts(result.data.slice(0, FEATURED_LIMIT));
          setFeaturedStatus("ready");
        }
      } catch {
        if (isActive) {
          setFeaturedProducts([]);
          setFeaturedStatus("error");
          setFeaturedError("Produk belum tersedia");
        }
      }
    }

    loadFeaturedProducts();

    return () => {
      isActive = false;
    };
  }, []);

  useEffect(() => {
    if (!isMobileNavOpen) {
      return;
    }

    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === "Escape") {
        setIsMobileNavOpen(false);
      }
    }

    window.addEventListener("keydown", handleKeyDown);

    return () => {
      window.removeEventListener("keydown", handleKeyDown);
    };
  }, [isMobileNavOpen]);

  function validateContactForm(input: ContactFormState): ContactFormErrors {
    const nextErrors: ContactFormErrors = {};

    if (!input.name.trim()) {
      nextErrors.name = "Nama wajib diisi.";
    }

    if (!input.email.trim()) {
      nextErrors.email = "Email wajib diisi.";
    } else if (!CONTACT_EMAIL_PATTERN.test(input.email.trim())) {
      nextErrors.email = "Format email tidak valid.";
    }

    if (!input.message.trim()) {
      nextErrors.message = "Pesan wajib diisi.";
    }

    return nextErrors;
  }

  function handleContactInputChange(
    fieldName: keyof ContactFormState,
    value: string,
  ) {
    setContactForm((currentForm) => ({
      ...currentForm,
      [fieldName]: value,
    }));

    setContactErrors((currentErrors) => {
      if (!currentErrors[fieldName]) {
        return currentErrors;
      }

      const nextErrors = { ...currentErrors };
      delete nextErrors[fieldName];
      return nextErrors;
    });
  }

  async function handleContactSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();

    const nextErrors = validateContactForm(contactForm);
    setContactErrors(nextErrors);

    if (Object.keys(nextErrors).length > 0) {
      return;
    }

    setIsContactSubmitting(true);

    try {
      await sendContactMessage({
        name: contactForm.name.trim(),
        email: contactForm.email.trim(),
        message: contactForm.message.trim(),
      });

      setContactForm(CONTACT_FORM_INITIAL_STATE);
      setContactErrors({});

      showToast({
        type: "success",
        message: "Pesan berhasil dikirim",
      });
    } catch (error) {
      showToast({
        type: "error",
        message: getContactErrorMessage(
          error,
          "Pesan gagal dikirim. Coba lagi nanti.",
        ),
      });
    } finally {
      setIsContactSubmitting(false);
    }
  }

  return (
    <main className="landing-page">
      <nav className="landing-navbar" aria-label="Landing navigation">
        <Link className="landing-brand" to="/">
          <span className="landing-brand-mark" aria-hidden="true">
            <Grid2X2 size={22} />
          </span>
          <span>Bjir E-Commerce</span>
        </Link>

        <button
          className="landing-nav-toggle"
          type="button"
          onClick={() => setIsMobileNavOpen((current) => !current)}
          aria-controls="landing-mobile-menu"
          aria-expanded={isMobileNavOpen}
          aria-label={
            isMobileNavOpen
              ? "Tutup menu navigasi landing"
              : "Buka menu navigasi landing"
          }
        >
          {isMobileNavOpen ? (
            <X size={20} aria-hidden="true" />
          ) : (
            <Menu size={20} aria-hidden="true" />
          )}
        </button>

        <div className="landing-nav-actions landing-nav-actions-desktop">
          {isAuthLoading ? (
            <>
              <span className="landing-nav-placeholder" aria-hidden="true" />
              <span
                className="landing-nav-placeholder is-short"
                aria-hidden="true"
              />
            </>
          ) : isAuthenticated ? (
            <>
              <Link className="landing-nav-link" to="/products">
                Products
              </Link>
              <Link className="landing-nav-cta" to={dashboardPath}>
                Dashboard
              </Link>
            </>
          ) : (
            <>
              <Link className="landing-nav-link" to="/login">
                Login
              </Link>
              <Link className="landing-nav-cta" to="/register">
                Daftar
              </Link>
            </>
          )}
        </div>

        {isMobileNavOpen && (
          <nav
            className="landing-mobile-nav"
            id="landing-mobile-menu"
            aria-label="Mobile landing navigation"
          >
            {isAuthLoading ? (
              <>
                <span className="landing-nav-placeholder" aria-hidden="true" />
                <span
                  className="landing-nav-placeholder is-short"
                  aria-hidden="true"
                />
              </>
            ) : isAuthenticated ? (
              <>
                <Link
                  className="landing-nav-link"
                  to="/products"
                  onClick={() => setIsMobileNavOpen(false)}
                >
                  Produk
                </Link>
                <Link
                  className="landing-nav-link"
                  to="/cart"
                  onClick={() => setIsMobileNavOpen(false)}
                >
                  Keranjang
                </Link>
                <Link
                  className="landing-nav-cta"
                  to={dashboardPath}
                  onClick={() => setIsMobileNavOpen(false)}
                >
                  Dashboard
                </Link>
              </>
            ) : (
              <>
                <Link
                  className="landing-nav-link"
                  to="/products"
                  onClick={() => setIsMobileNavOpen(false)}
                >
                  Produk
                </Link>
                <Link
                  className="landing-nav-link"
                  to="/cart"
                  onClick={() => setIsMobileNavOpen(false)}
                >
                  Keranjang
                </Link>
                <Link
                  className="landing-nav-link"
                  to="/login"
                  onClick={() => setIsMobileNavOpen(false)}
                >
                  Login
                </Link>
                <Link
                  className="landing-nav-cta"
                  to="/register"
                  onClick={() => setIsMobileNavOpen(false)}
                >
                  Daftar
                </Link>
              </>
            )}
          </nav>
        )}
      </nav>

      <section
        className={`landing-hero${isAuthLoading ? " is-loading" : ""}`}
        aria-labelledby="landing-title"
      >
        <div className="landing-hero-copy">
          <span className="landing-eyebrow">
            <Sparkles size={15} aria-hidden="true" />
            Pratama Enterprise
          </span>

          <h1 id="landing-title">
            Commerce with elbows. Loud shelves, fast checkout.
          </h1>

          <p>
            Bjir E-Commerce brings the catalog, cart, checkout, and admin shelf
            into one sharp storefront built for products that refuse to look
            generic.
          </p>

          {isAuthLoading ? (
            <div
              className="landing-hero-actions is-loading"
              aria-label="Loading hero actions"
            >
              <span
                className="landing-hero-cta-placeholder is-wide"
                aria-hidden="true"
              />
              <span
                className="landing-hero-cta-placeholder"
                aria-hidden="true"
              />
            </div>
          ) : (
            <div className="landing-hero-actions">
              <Link className="landing-primary-button" to={primaryHeroPath}>
                {primaryHeroLabel}
                <ArrowRight size={17} aria-hidden="true" />
              </Link>
              <Link className="landing-secondary-button" to={secondaryHeroPath}>
                {secondaryHeroLabel}
              </Link>
            </div>
          )}
        </div>

        <aside className="landing-hero-panel" aria-label="Store highlights">
          <div className="landing-hero-ticket">
            <span>01</span>
            <strong>Catalog</strong>
            <small>Search, categories, product details.</small>
          </div>
          <div className="landing-hero-ticket accent">
            <span>02</span>
            <strong>Checkout</strong>
            <small>Cart flow, shipping, orders.</small>
          </div>
          <div className="landing-hero-ticket blue">
            <span>03</span>
            <strong>Admin</strong>
            <small>Products, categories, order control.</small>
          </div>
        </aside>
      </section>

      <section
        className="landing-section"
        aria-labelledby="landing-featured-title"
      >
        <div className="landing-section-heading">
          <span className="landing-eyebrow">Featured Products</span>
          <h2 id="landing-featured-title">Produk unggulan terbaru.</h2>
          <p>
            Delapan produk terbaru dari katalog, lengkap dengan gambar, harga,
            dan status stok.
          </p>
        </div>

        {featuredStatus === "loading" ? (
          <LandingFeaturedSkeleton />
        ) : featuredStatus === "error" ? (
          <div className="landing-featured-state" role="status">
            {featuredError}
          </div>
        ) : featuredProducts.length === 0 ? (
          <div className="landing-featured-state" role="status">
            Produk belum tersedia
          </div>
        ) : (
          <div className="landing-featured-grid">
            {featuredProducts.map((product) => {
              const stockState = getStockState(product.stock);

              return (
                <Link
                  className="landing-product-card"
                  to={`/products/${product.slug}`}
                  key={product.id}
                  aria-label={`Lihat produk ${product.name}`}
                >
                  <ProductImage
                    className="landing-product-image"
                    src={getProductImage(product)}
                    alt={product.name}
                    width={640}
                    height={480}
                    sizes="(max-width: 720px) 100vw, (max-width: 1180px) 50vw, 25vw"
                  />

                  <div className="landing-product-body">
                    <span>{product.category?.name || "Uncategorized"}</span>
                    <h3>{product.name}</h3>

                    <div className="landing-product-meta">
                      <strong>{formatRupiah(product.price)}</strong>
                      <span
                        className={`landing-product-stock ${stockState.className}`}
                      >
                        {stockState.label}
                      </span>
                    </div>
                  </div>
                </Link>
              );
            })}
          </div>
        )}
      </section>

      <section className="landing-about" aria-labelledby="landing-about-title">
        <div className="landing-about-story">
          <span className="landing-eyebrow">Tentang Bjir</span>
          <h2 id="landing-about-title">Belanja cepat, aman, tanpa drama.</h2>
          <p>
            Bjir E-Commerce dibuat sebagai etalase digital yang rapi, tegas, dan
            langsung ke inti: pembeli bisa menemukan produk, melihat stok,
            checkout, lalu memantau pesanan tanpa alur yang bertele-tele.
          </p>
        </div>

        <div
          className="landing-about-grid"
          aria-label="Keunggulan Bjir E-Commerce"
        >
          <article>
            <Truck aria-hidden="true" />
            <h3>Pengiriman Cepat</h3>
            <p>
              Pesanan diproses dengan alur yang jelas supaya barang lebih cepat
              masuk antrean pengiriman.
            </p>
          </article>

          <article>
            <ShieldCheck aria-hidden="true" />
            <h3>Pembayaran Aman</h3>
            <p>
              Flow checkout dibuat terstruktur, dengan status pembayaran dan
              pesanan yang mudah dipantau.
            </p>
          </article>

          <article>
            <BadgeCheck aria-hidden="true" />
            <h3>Produk Original</h3>
            <p>
              Katalog menampilkan informasi produk, harga, gambar, dan stok
              secara transparan.
            </p>
          </article>

          <article>
            <RotateCcw aria-hidden="true" />
            <h3>Garansi Uang Kembali</h3>
            <p>
              Belanja lebih tenang dengan dukungan pengembalian dana sesuai
              kebijakan toko.
            </p>
          </article>
        </div>
      </section>

      <section
        className="landing-contact"
        aria-labelledby="landing-contact-title"
      >
        <div className="landing-contact-copy">
          <span className="landing-eyebrow">Kontak Admin</span>
          <h2 id="landing-contact-title">Ada pertanyaan untuk toko?</h2>
          <p>
            Kirim pesan ke admin untuk bantuan pesanan, katalog, atau akses
            toko. Tim Bjir akan membaca pesanmu dari jalur support.
          </p>

          <div className="landing-contact-list" aria-label="Kontak cepat">
            <a href="mailto:yudisbaek@gmail.com">
              <Mail size={18} aria-hidden="true" />
              yudisbaek@gmail.com
            </a>
            <a href="tel:+62000000000">
              <Phone size={18} aria-hidden="true" />
              +62 000 000 000
            </a>
            <span>
              <MapPin size={18} aria-hidden="true" />
              Malang, Jawa Timur, Indonesia
            </span>
            <span>
              <Truck size={18} aria-hidden="true" />
                Pengiriman cepat, aman
            </span>
          </div>
        </div>

        <form
          className="landing-contact-form"
          onSubmit={handleContactSubmit}
          noValidate
        >
          <div className="landing-form-field">
            <label htmlFor="landing-contact-name">Nama</label>
            <input
              id="landing-contact-name"
              type="text"
              value={contactForm.name}
              onChange={(event) =>
                handleContactInputChange("name", event.target.value)
              }
              aria-invalid={Boolean(contactErrors.name)}
              aria-describedby={
                contactErrors.name ? "landing-contact-name-error" : undefined
              }
              disabled={isContactSubmitting}
            />
            {contactErrors.name && (
              <p className="landing-form-error" id="landing-contact-name-error">
                {contactErrors.name}
              </p>
            )}
          </div>

          <div className="landing-form-field">
            <label htmlFor="landing-contact-email">Email</label>
            <input
              id="landing-contact-email"
              type="email"
              value={contactForm.email}
              onChange={(event) =>
                handleContactInputChange("email", event.target.value)
              }
              aria-invalid={Boolean(contactErrors.email)}
              aria-describedby={
                contactErrors.email ? "landing-contact-email-error" : undefined
              }
              disabled={isContactSubmitting}
            />
            {contactErrors.email && (
              <p
                className="landing-form-error"
                id="landing-contact-email-error"
              >
                {contactErrors.email}
              </p>
            )}
          </div>

          <div className="landing-form-field">
            <label htmlFor="landing-contact-message">Pesan</label>
            <textarea
              id="landing-contact-message"
              rows={5}
              value={contactForm.message}
              onChange={(event) =>
                handleContactInputChange("message", event.target.value)
              }
              aria-invalid={Boolean(contactErrors.message)}
              aria-describedby={
                contactErrors.message
                  ? "landing-contact-message-error"
                  : undefined
              }
              disabled={isContactSubmitting}
            />
            {contactErrors.message && (
              <p
                className="landing-form-error"
                id="landing-contact-message-error"
              >
                {contactErrors.message}
              </p>
            )}
          </div>

          <button
            className="landing-contact-submit"
            type="submit"
            disabled={isContactSubmitting}
          >
            {isContactSubmitting ? "Mengirim..." : "Kirim Pesan"}
          </button>
        </form>
      </section>

      <footer className="landing-footer">
        <div className="landing-footer-separator" aria-hidden="true" />

        <div className="landing-footer-main">
          <div className="landing-footer-brand">
            <span className="landing-footer-mark" aria-hidden="true">
              <ShoppingCart size={22} />
            </span>
            <div>
              <strong>Bjir E-Commerce</strong>
              <p>Belanja cepat, stok jelas, checkout tanpa drama.</p>
            </div>
          </div>

          <nav className="landing-footer-nav" aria-label="Footer quick links">
            <h2>Quick Links</h2>
            <div>
              {FOOTER_QUICK_LINKS.map((link) => (
                <Link to={link.to} key={link.to}>
                  {link.label}
                </Link>
              ))}
            </div>
          </nav>

          <nav className="landing-footer-nav" aria-label="Footer social links">
            <h2>Social</h2>
            <div>
              {FOOTER_SOCIAL_LINKS.map((link) => {
                const Icon = link.icon;

                return (
                  <a
                    href={link.href}
                    key={link.href}
                    target={link.href.startsWith("http") ? "_blank" : undefined}
                    rel={
                      link.href.startsWith("http")
                        ? "noreferrer noopener"
                        : undefined
                    }
                  >
                    <Icon size={16} aria-hidden="true" />
                    {link.label}
                  </a>
                );
              })}
            </div>
          </nav>
        </div>

        <div className="landing-footer-bottom">
          <span>© {currentYear} Pratama Enterprise</span>
          <span>Built with love, made with passion.</span>
        </div>
      </footer>
    </main>
  );
}

export default LandingPage;
