import { CartItem, Medicine } from './types';

class CartStore {
  private items: CartItem[] = [];
  private listeners: (() => void)[] = [];

  constructor() {
    const saved = localStorage.getItem('cart');
    if (saved) {
      try {
        this.items = JSON.parse(saved);
      } catch {
        this.items = [];
      }
    }
  }

  public getItems(): CartItem[] {
    return [...this.items];
  }

  public getItemCount(): number {
    return this.items.reduce((sum, item) => sum + item.quantity, 0);
  }

  public getTotal(): number {
    return this.items.reduce((sum, item) => sum + item.medicine.price * item.quantity, 0);
  }

  public requiresPrescription(): boolean {
    return this.items.some((item) => item.medicine.requires_rx);
  }

  public addItem(medicine: Medicine, quantity: number = 1) {
    if (quantity <= 0) return;
    const existing = this.items.find((i) => i.medicine.medicine_id === medicine.medicine_id);
    if (existing) {
      existing.quantity += quantity;
    } else {
      this.items.push({ medicine, quantity });
    }
    this.save();
  }

  public updateQuantity(medicineId: string, quantity: number) {
    if (quantity <= 0) {
      this.removeItem(medicineId);
      return;
    }
    const item = this.items.find((i) => i.medicine.medicine_id === medicineId);
    if (item) {
      item.quantity = quantity;
      this.save();
    }
  }

  public removeItem(medicineId: string) {
    this.items = this.items.filter((i) => i.medicine.medicine_id !== medicineId);
    this.save();
  }

  public clear() {
    this.items = [];
    this.save();
  }

  public subscribe(listener: () => void): () => void {
    this.listeners.push(listener);
    return () => {
      this.listeners = this.listeners.filter((l) => l !== listener);
    };
  }

  private save() {
    localStorage.setItem('cart', JSON.stringify(this.items));
    this.listeners.forEach((l) => l());
  }
}

export const cart = new CartStore();
