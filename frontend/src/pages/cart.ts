import { api } from '../api';
import { cart } from '../cart';
import { Prescription } from '../types';

export async function renderCartPage(container: HTMLElement) {
  const items = cart.getItems();
  const requiresRx = cart.requiresPrescription();
  const user = api.getCurrentUser();

  let approvedPrescriptions: Prescription[] = [];
  if (user && requiresRx) {
    try {
      const myRx = await api.getMyPrescriptions();
      approvedPrescriptions = myRx.filter((p) => p.review_status === 'APPROVED');
    } catch (err) {
      console.error('Failed to load prescriptions:', err);
    }
  }

  const renderContent = () => {
    if (items.length === 0) {
      container.innerHTML = `
        <div style="max-width: 600px; margin: 40px auto; text-align: center; background: white; padding: 40px; border-radius: var(--radius-lg); border: 1px solid var(--border);">
          <div style="font-size: 3rem; margin-bottom: 12px;">🛒</div>
          <h2>Your Cart is Empty</h2>
          <p style="color: var(--text-muted); margin: 12px 0 24px;">Browse our comprehensive catalog of verified OTC and prescription medications.</p>
          <a href="#/medicines" class="btn btn-primary">Browse Catalogue</a>
        </div>
      `;
      return;
    }

    const total = cart.getTotal();

    container.innerHTML = `
      <div style="max-width: 900px; margin: 0 auto;">
        <h1 style="font-size: 1.8rem; font-weight: 800; margin-bottom: 20px;">Shopping Cart & Checkout</h1>

        <div id="checkout-error" class="alert alert-danger" style="display: none;"></div>

        ${requiresRx ? `
          <div class="alert alert-warning" style="display: flex; gap: 12px; align-items: flex-start;">
            <span style="font-size: 1.3rem;">⚠️</span>
            <div>
              <strong>Prescription-Only Medication in Cart</strong>
              <p style="font-size: 0.88rem; margin-top: 4px;">
                Safety Rule: Items marked with <span class="badge badge-rx">Rx Required</span> require an approved prescription reviewed by our licensed pharmacist before an order can be fulfilled.
              </p>
            </div>
          </div>
        ` : ''}

        <div class="table-container" style="margin-bottom: 24px;">
          <table>
            <thead>
              <tr>
                <th>Medicine</th>
                <th>Type</th>
                <th>Price</th>
                <th>Quantity</th>
                <th>Subtotal</th>
                <th>Action</th>
              </tr>
            </thead>
            <tbody>
              ${items.map((item) => `
                <tr>
                  <td>
                    <strong>${escapeHtml(item.medicine.name)}</strong>
                    <div style="font-size: 0.8rem; color: var(--text-muted);">${escapeHtml(item.medicine.generic_name)}</div>
                  </td>
                  <td>
                    <span class="badge ${item.medicine.requires_rx ? 'badge-rx' : 'badge-otc'}">
                      ${item.medicine.requires_rx ? 'Rx Required' : 'OTC'}
                    </span>
                  </td>
                  <td>$${item.medicine.price.toFixed(2)}</td>
                  <td>
                    <div style="display: flex; align-items: center; gap: 6px;">
                      <button class="btn btn-outline btn-sm dec-btn" data-id="${item.medicine.medicine_id}">-</button>
                      <span style="font-weight: 700; min-width: 20px; text-align: center;">${item.quantity}</span>
                      <button class="btn btn-outline btn-sm inc-btn" data-id="${item.medicine.medicine_id}">+</button>
                    </div>
                  </td>
                  <td style="font-weight: 700;">$${(item.medicine.price * item.quantity).toFixed(2)}</td>
                  <td>
                    <button class="btn btn-outline btn-sm remove-btn" data-id="${item.medicine.medicine_id}" style="color: var(--danger); border-color: #fca5a5;">
                      Remove
                    </button>
                  </td>
                </tr>
              `).join('')}
            </tbody>
          </table>
        </div>

        <!-- Checkout Summary Card -->
        <div class="card" style="padding: 24px;">
          <div style="display: flex; justify-content: space-between; font-size: 1.3rem; font-weight: 800; margin-bottom: 16px;">
            <span>Order Total:</span>
            <span class="price">$${total.toFixed(2)}</span>
          </div>

          ${requiresRx ? `
            <div class="form-group" style="margin-top: 10px;">
              <label for="prescription-select">Attach Approved Prescription (Mandatory for Rx Items):</label>
              ${approvedPrescriptions.length > 0 ? `
                <select id="prescription-select" class="form-select">
                  <option value="">-- Select an approved prescription --</option>
                  ${approvedPrescriptions.map((p) => `
                    <option value="${p.prescription_id}">
                      Prescription #${p.prescription_id.substring(0, 8)} (Approved - ${new Date(p.created_at).toLocaleDateString()})
                    </option>
                  `).join('')}
                </select>
              ` : `
                <div class="alert alert-danger" style="margin-top: 6px;">
                  You do not have any <strong>Approved</strong> prescriptions on file.
                  <div style="margin-top: 8px;">
                    <a href="#/prescriptions" class="btn btn-primary btn-sm">Upload Prescription Now</a>
                  </div>
                </div>
              `}
            </div>
          ` : ''}

          <div style="display: flex; justify-content: space-between; align-items: center; margin-top: 20px; flex-wrap: wrap; gap: 12px;">
            <a href="#/medicines" class="btn btn-outline">&larr; Continue Shopping</a>
            <button id="checkout-btn" class="btn btn-primary" style="padding: 10px 24px;">
              Place Order ($${total.toFixed(2)})
            </button>
          </div>
        </div>
      </div>
    `;

    // Attach event listeners
    container.querySelectorAll('.inc-btn').forEach((btn) => {
      btn.addEventListener('click', (e) => {
        const id = (e.currentTarget as HTMLElement).getAttribute('data-id')!;
        const it = items.find((i) => i.medicine.medicine_id === id);
        if (it) {
          cart.updateQuantity(id, it.quantity + 1);
          renderCartPage(container);
        }
      });
    });

    container.querySelectorAll('.dec-btn').forEach((btn) => {
      btn.addEventListener('click', (e) => {
        const id = (e.currentTarget as HTMLElement).getAttribute('data-id')!;
        const it = items.find((i) => i.medicine.medicine_id === id);
        if (it) {
          cart.updateQuantity(id, it.quantity - 1);
          renderCartPage(container);
        }
      });
    });

    container.querySelectorAll('.remove-btn').forEach((btn) => {
      btn.addEventListener('click', (e) => {
        const id = (e.currentTarget as HTMLElement).getAttribute('data-id')!;
        cart.removeItem(id);
        renderCartPage(container);
      });
    });

    const checkoutBtn = container.querySelector('#checkout-btn') as HTMLButtonElement | null;
    const errorDiv = container.querySelector('#checkout-error') as HTMLDivElement;

    if (checkoutBtn) {
      checkoutBtn.addEventListener('click', async () => {
        if (!api.getToken()) {
          window.location.hash = '#/login';
          return;
        }

        let prescriptionId: string | undefined = undefined;
        if (requiresRx) {
          const select = container.querySelector('#prescription-select') as HTMLSelectElement | null;
          prescriptionId = select?.value;
          if (!prescriptionId) {
            errorDiv.textContent = 'Please select an approved prescription to fulfill prescription-required items.';
            errorDiv.style.display = 'block';
            return;
          }
        }

        try {
          checkoutBtn.disabled = true;
          checkoutBtn.textContent = 'Placing Order...';
          errorDiv.style.display = 'none';

          const orderItems = items.map((i) => ({
            medicine_id: i.medicine.medicine_id,
            quantity: i.quantity,
          }));

          const created = await api.createOrder(orderItems, prescriptionId);
          cart.clear();
          window.location.hash = `#/orders?placed=${created.order_id}`;
        } catch (err: any) {
          errorDiv.textContent = err.message || 'Failed to place order';
          errorDiv.style.display = 'block';
          checkoutBtn.disabled = false;
          checkoutBtn.textContent = `Place Order ($${total.toFixed(2)})`;
        }
      });
    }
  };

  renderContent();
}

function escapeHtml(str: string): string {
  const div = document.createElement('div');
  div.textContent = str;
  return div.innerHTML;
}
