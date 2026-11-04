import { describe, it, expect, beforeEach } from 'vitest';
import { cart } from '../cart';
import { Medicine } from '../types';

describe('CartStore Logic', () => {
  const otcMedicine: Medicine = {
    medicine_id: 'med-otc-1',
    name: 'Paracetamol 500mg',
    generic_name: 'Paracetamol',
    description: 'Pain relief',
    price: 5.0,
    stock_qty: 100,
    requires_rx: false,
  };

  const rxMedicine: Medicine = {
    medicine_id: 'med-rx-1',
    name: 'Amoxicillin 500mg',
    generic_name: 'Amoxicillin',
    description: 'Antibiotic',
    price: 15.0,
    stock_qty: 50,
    requires_rx: true,
  };

  beforeEach(() => {
    cart.clear();
  });

  it('starts with an empty cart', () => {
    expect(cart.getItems().length).toBe(0);
    expect(cart.getItemCount()).toBe(0);
    expect(cart.getTotal()).toBe(0);
    expect(cart.requiresPrescription()).toBe(false);
  });

  it('adds items and calculates total correctly', () => {
    cart.addItem(otcMedicine, 2);
    expect(cart.getItemCount()).toBe(2);
    expect(cart.getTotal()).toBe(10.0);
    expect(cart.requiresPrescription()).toBe(false);

    // Add more of the same item
    cart.addItem(otcMedicine, 1);
    expect(cart.getItemCount()).toBe(3);
    expect(cart.getTotal()).toBe(15.0);
  });

  it('detects prescription requirement when an Rx item is added', () => {
    cart.addItem(otcMedicine, 1);
    expect(cart.requiresPrescription()).toBe(false);

    cart.addItem(rxMedicine, 1);
    expect(cart.requiresPrescription()).toBe(true);
    expect(cart.getTotal()).toBe(20.0);
  });

  it('updates quantity and recalculates totals', () => {
    cart.addItem(otcMedicine, 1);
    cart.updateQuantity(otcMedicine.medicine_id, 4);

    expect(cart.getItemCount()).toBe(4);
    expect(cart.getTotal()).toBe(20.0);

    // Setting quantity to 0 removes item
    cart.updateQuantity(otcMedicine.medicine_id, 0);
    expect(cart.getItemCount()).toBe(0);
  });

  it('removes specific items cleanly', () => {
    cart.addItem(otcMedicine, 1);
    cart.addItem(rxMedicine, 2);

    expect(cart.getItems().length).toBe(2);

    cart.removeItem(rxMedicine.medicine_id);
    expect(cart.getItems().length).toBe(1);
    expect(cart.requiresPrescription()).toBe(false);
    expect(cart.getTotal()).toBe(5.0);
  });

  it('clears all items', () => {
    cart.addItem(otcMedicine, 2);
    cart.addItem(rxMedicine, 1);
    expect(cart.getItemCount()).toBe(3);

    cart.clear();
    expect(cart.getItemCount()).toBe(0);
    expect(cart.getTotal()).toBe(0);
  });
});
