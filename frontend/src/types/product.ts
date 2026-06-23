export type SortOrder = "asc" | "desc";

export interface Category {
  id: string;
  name: string;
  slug: string;
  description?: string | null;
  image_url?: string | null;
  parent_id?: string | null;
  parent?: Category | null;
  children?: Category[];
  created_at?: string;
  updated_at?: string;
}

export interface ProductImage {
  id: string;
  product_id: string;
  image_url: string;
  sort_order: number;
  created_at?: string;
  updated_at?: string;
}

export interface Product {
  id: string;
  name: string;
  slug: string;
  description?: string | null;
  price: number;
  stock: number;
  category_id?: string | null;
  category?: Category | null;
  image_url?: string | null;
  images?: ProductImage[];
  created_at?: string;
  updated_at?: string;
}

export interface ProductListParams {
  category_id?: string;
  category?: string;
  search?: string;
  sort_by?: string;
  sort_order?: SortOrder;
  page?: number;
  limit?: number;
}

export interface ProductListMeta {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
  sort_by?: string;
  sort_order?: SortOrder | "";
  category_id?: string;
  category?: string;
  search?: string;
}

export interface ProductListResponse {
  data: Product[];
  meta: ProductListMeta;
}