import { describe, it, expect, beforeEach } from 'vitest';
import { api } from '../api';
import { User } from '../types';

describe('ApiClient auth & helper methods', () => {
  beforeEach(() => {
    api.logout();
  });

  it('stores and retrieves token and currentUser', () => {
    expect(api.getToken()).toBeNull();
    expect(api.getCurrentUser()).toBeNull();

    const mockUser: User = {
      user_id: 'test-user-id',
      name: 'Test Pharmacist',
      email: 'pharmacist@test.local',
      role: 'staff',
    };

    api.setAuth('sample-jwt-token-xyz', mockUser);

    expect(api.getToken()).toBe('sample-jwt-token-xyz');
    expect(api.getCurrentUser()).toEqual(mockUser);
    expect(localStorage.getItem('token')).toBe('sample-jwt-token-xyz');

    api.logout();
    expect(api.getToken()).toBeNull();
    expect(api.getCurrentUser()).toBeNull();
    expect(localStorage.getItem('token')).toBeNull();
  });

  it('generates correct prescription file URLs', () => {
    const rxId = 'rx-12345-abcd';
    const url = api.getPrescriptionFileUrl(rxId);
    expect(url).toBe('http://localhost:8080/api/prescriptions/rx-12345-abcd/file');
  });
});
