import {
  AuthResponse,
  User,
  Medicine,
  MedicinesListResponse,
  Prescription,
  Order,
  OrderStatus,
  UserRole
} from './types';

const API_BASE = 'http://localhost:8080/api';

class ApiClient {
  private token: string | null = null;
  private currentUser: User | null = null;

  constructor() {
    this.token = localStorage.getItem('token');
    const savedUser = localStorage.getItem('user');
    if (savedUser) {
      try {
        this.currentUser = JSON.parse(savedUser);
      } catch {
        this.currentUser = null;
      }
    }
  }

  public getToken(): string | null {
    return this.token;
  }

  public getCurrentUser(): User | null {
    return this.currentUser;
  }

  public setAuth(token: string, user: User) {
    this.token = token;
    this.currentUser = user;
    localStorage.setItem('token', token);
    localStorage.setItem('user', JSON.stringify(user));
  }

  public logout() {
    this.token = null;
    this.currentUser = null;
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  }

  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${API_BASE}${endpoint}`;
    const headers = new Headers(options.headers || {});

    if (this.token && !headers.has('Authorization')) {
      headers.set('Authorization', `Bearer ${this.token}`);
    }

    if (!(options.body instanceof FormData) && !headers.has('Content-Type')) {
      headers.set('Content-Type', 'application/json');
    }

    const response = await fetch(url, {
      ...options,
      headers,
    });

    if (response.status === 204) {
      return {} as T;
    }

    const data = await response.json().catch(() => ({}));

    if (!response.ok) {
      const errMsg = (data && data.error) ? data.error : `HTTP ${response.status}: Request failed`;
      throw new Error(errMsg);
    }

    return data as T;
  }

  // Auth
  public async login(email: string, password: string): Promise<AuthResponse> {
    const res = await this.request<AuthResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
    this.setAuth(res.token, res.user);
    return res;
  }

  public async register(name: string, email: string, password: string): Promise<AuthResponse> {
    const res = await this.request<AuthResponse>('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ name, email, password }),
    });
    this.setAuth(res.token, res.user);
    return res;
  }

  public async getMe(): Promise<User> {
    const user = await this.request<User>('/me');
    this.currentUser = user;
    localStorage.setItem('user', JSON.stringify(user));
    return user;
  }

  // Medicines
  public async getMedicines(query: string = '', page: number = 1, limit: number = 20): Promise<MedicinesListResponse> {
    const params = new URLSearchParams({
      q: query,
      page: page.toString(),
      limit: limit.toString(),
    });
    return this.request<MedicinesListResponse>(`/medicines?${params.toString()}`);
  }

  public async getMedicine(id: string): Promise<Medicine> {
    return this.request<Medicine>(`/medicines/${id}`);
  }

  public async createMedicine(med: Partial<Medicine>): Promise<Medicine> {
    return this.request<Medicine>('/medicines', {
      method: 'POST',
      body: JSON.stringify(med),
    });
  }

  public async updateMedicine(id: string, med: Partial<Medicine>): Promise<Medicine> {
    return this.request<Medicine>(`/medicines/${id}`, {
      method: 'PUT',
      body: JSON.stringify(med),
    });
  }

  public async updateStock(id: string, stockQty?: number, price?: number): Promise<Medicine> {
    return this.request<Medicine>(`/medicines/${id}/stock`, {
      method: 'PATCH',
      body: JSON.stringify({ stock_qty: stockQty, price }),
    });
  }

  // Prescriptions
  public async uploadPrescription(file: File): Promise<Prescription> {
    const formData = new FormData();
    formData.append('file', file);

    return this.request<Prescription>('/prescriptions', {
      method: 'POST',
      body: formData,
    });
  }

  public async getMyPrescriptions(): Promise<Prescription[]> {
    return this.request<Prescription[]>('/prescriptions/mine');
  }

  public async getPendingPrescriptions(): Promise<Prescription[]> {
    return this.request<Prescription[]>('/prescriptions/queue');
  }

  public async getPrescription(id: string): Promise<Prescription> {
    return this.request<Prescription>(`/prescriptions/${id}`);
  }

  public getPrescriptionFileUrl(id: string): string {
    return `${API_BASE}/prescriptions/${id}/file`;
  }

  public async reviewPrescription(id: string, action: 'APPROVE' | 'REJECT', note: string): Promise<Prescription> {
    return this.request<Prescription>(`/prescriptions/${id}/review`, {
      method: 'POST',
      body: JSON.stringify({ action, note }),
    });
  }

  // Orders
  public async createOrder(items: { medicine_id: string; quantity: number }[], prescriptionId?: string): Promise<Order> {
    return this.request<Order>('/orders', {
      method: 'POST',
      body: JSON.stringify({
        prescription_id: prescriptionId || null,
        items,
      }),
    });
  }

  public async getMyOrders(): Promise<Order[]> {
    return this.request<Order[]>('/orders/mine');
  }

  public async getOrder(id: string): Promise<Order> {
    return this.request<Order>(`/orders/${id}`);
  }

  public async getAllOrders(): Promise<Order[]> {
    return this.request<Order[]>('/orders');
  }

  public async updateOrderStatus(id: string, status: OrderStatus): Promise<Order> {
    return this.request<Order>(`/orders/${id}/status`, {
      method: 'PATCH',
      body: JSON.stringify({ status }),
    });
  }

  // Admin
  public async getUsers(): Promise<User[]> {
    return this.request<User[]>('/users');
  }

  public async updateUserRole(id: string, role: UserRole): Promise<User> {
    return this.request<User>(`/users/${id}/role`, {
      method: 'PATCH',
      body: JSON.stringify({ role }),
    });
  }
}

export const api = new ApiClient();
