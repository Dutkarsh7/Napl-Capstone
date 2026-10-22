export type UserRole = 'customer' | 'staff' | 'admin';

export interface User {
  user_id: string;
  name: string;
  email: string;
  role: UserRole;
  created_at?: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface Medicine {
  medicine_id: string;
  name: string;
  generic_name: string;
  description: string;
  price: number;
  stock_qty: number;
  requires_rx: boolean;
  created_at?: string;
}

export interface MedicinesListResponse {
  data: Medicine[];
  total: number;
  page: number;
  limit: number;
}

export type ValidationStatus = 'PENDING' | 'VALID' | 'INCOMPLETE' | 'INVALID' | 'MANUAL_REVIEW';
export type ReviewStatus = 'PENDING' | 'APPROVED' | 'REJECTED';

export interface Prescription {
  prescription_id: string;
  user_id: string;
  file_path: string;
  extracted_text: string;
  validation_status: ValidationStatus;
  review_status: ReviewStatus;
  reviewed_by?: string;
  review_note: string;
  created_at: string;
}

export type OrderStatus = 'PLACED' | 'UNDER_REVIEW' | 'CONFIRMED' | 'READY' | 'COMPLETED' | 'CANCELLED';

export interface OrderItem {
  item_id: string;
  order_id: string;
  medicine_id: string;
  medicine_name?: string;
  quantity: number;
  unit_price: number;
  requires_rx?: boolean;
}

export interface OrderStatusHistory {
  id: string;
  order_id: string;
  status: OrderStatus;
  changed_by?: string;
  changed_at: string;
}

export interface Order {
  order_id: string;
  user_id: string;
  customer_name?: string;
  customer_email?: string;
  prescription_id?: string;
  status: OrderStatus;
  total: number;
  created_at: string;
  updated_at: string;
  items?: OrderItem[];
  history?: OrderStatusHistory[];
}

export interface CartItem {
  medicine: Medicine;
  quantity: number;
}

export interface ApiError {
  error: string;
}
