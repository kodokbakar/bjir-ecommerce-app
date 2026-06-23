import {
  Grid2X2,
  Home,
  Package,
  Settings,
  ShoppingBag,
  ShoppingCart,
  type LucideIcon,
} from "lucide-react";

export const APP_VERSION = "1.0.0";

export interface NavigationItem {
  label: string;
  path: string;
  Icon: LucideIcon;
}

export interface BreadcrumbItem {
  label: string;
  path?: string;
}

export const MAIN_NAV_ITEMS: NavigationItem[] = [
  {
    label: "Dashboard",
    path: "/dashboard",
    Icon: Home,
  },
  {
    label: "Produk",
    path: "/products",
    Icon: Package,
  },
  {
    label: "Keranjang",
    path: "/cart",
    Icon: ShoppingCart,
  },
  {
    label: "Pesanan",
    path: "/orders",
    Icon: ShoppingBag,
  },
  {
    label: "Profil",
    path: "/profile",
    Icon: Settings,
  },
];

export function getPageTitle(pathname: string): string {
  if (pathname.startsWith("/categories")) {
    return "Categories";
  }

  if (pathname.startsWith("/products/")) {
    return "Product Detail";
  }

  if (pathname.startsWith("/checkout")) {
    return "Checkout";
  }

  if (pathname.startsWith("/orders/")) {
    return "Order Detail";
  }

  if (pathname.startsWith("/profile")) {
    return "Profile";
  }

  const activeItem = MAIN_NAV_ITEMS.find((item) =>
    pathname.startsWith(item.path),
  );

  return activeItem?.label ?? "Halaman";
}

export function getBreadcrumbs(pathname: string): BreadcrumbItem[] {
  const breadcrumbs: BreadcrumbItem[] = [
    {
      label: "Home",
      path: "/dashboard",
    },
  ];

  if (pathname.startsWith("/products/")) {
    breadcrumbs.push(
      {
        label: "Products",
        path: "/products",
      },
      {
        label: "Product Detail",
      },
    );

    return breadcrumbs;
  }

  if (pathname === "/products") {
    breadcrumbs.push({
      label: "Products",
    });

    return breadcrumbs;
  }

  if (pathname.startsWith("/categories")) {
    breadcrumbs.push(
      {
        label: "Categories",
        path: "/products",
      },
      {
        label: "Category",
      },
    );

    return breadcrumbs;
  }

  if (pathname.startsWith("/cart")) {
    breadcrumbs.push({
      label: "Cart",
    });

    return breadcrumbs;
  }

  if (pathname.startsWith("/checkout")) {
    breadcrumbs.push(
      {
        label: "Cart",
        path: "/cart",
      },
      {
        label: "Checkout",
      },
    );

    return breadcrumbs;
  }

  if (pathname.startsWith("/orders/")) {
    breadcrumbs.push(
      {
        label: "Orders",
        path: "/orders",
      },
      {
        label: "Order Detail",
      },
    );

    return breadcrumbs;
  }

  if (pathname.startsWith("/orders")) {
    breadcrumbs.push({
      label: "Orders",
    });

    return breadcrumbs;
  }

  if (pathname.startsWith("/profile") || pathname.startsWith("/settings")) {
    breadcrumbs.push({
      label: "Profile",
    });

    return breadcrumbs;
  }

  return breadcrumbs;
}

export { Grid2X2 };
