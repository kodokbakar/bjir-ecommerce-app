import { useCallback, useState, type ReactNode } from "react";

import { useCartCount } from "../../hooks/useCartCount";
import Footer from "./Footer";
import Header from "./Header";
import Sidebar from "./Sidebar";

interface LayoutProps {
  children: ReactNode;
}

function Layout({ children }: LayoutProps) {
  const [isSidebarCollapsed, setIsSidebarCollapsed] = useState(false);
  const [isMobileSidebarOpen, setIsMobileSidebarOpen] = useState(false);
  const { count: cartCount } = useCartCount();

  const handleToggleCollapse = useCallback(() => {
    setIsSidebarCollapsed((current) => !current);
  }, []);

  const handleOpenMobile = useCallback(() => {
    setIsMobileSidebarOpen(true);
  }, []);

  const handleCloseMobile = useCallback(() => {
    setIsMobileSidebarOpen(false);
  }, []);

  return (
    <div className="min-h-screen bg-[#f8f6f4]">
      <Sidebar
        isCollapsed={isSidebarCollapsed}
        isMobileOpen={isMobileSidebarOpen}
        cartCount={cartCount}
        onCloseMobile={handleCloseMobile}
      />

      <div
        className={[
          "flex min-h-screen flex-col transition-[padding-left] duration-300",
          "md:pl-[76px]",
          isSidebarCollapsed ? "lg:pl-[76px]" : "lg:pl-[260px]",
        ].join(" ")}
      >
        <Header
          isCollapsed={isSidebarCollapsed}
          cartCount={cartCount}
          onToggleCollapse={handleToggleCollapse}
          onOpenMobile={handleOpenMobile}
        />

        <main className="flex flex-1 flex-col px-4 py-5 sm:px-6 lg:px-8 lg:py-8">
          {children}
        </main>

        <Footer />
      </div>
    </div>
  );
}

export default Layout;
