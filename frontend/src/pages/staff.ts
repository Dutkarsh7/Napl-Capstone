import { api } from '../api';
import { Order, OrderStatus, Prescription, Medicine } from '../types';

export function renderStaffPage(container: HTMLElement) {
  const user = api.getCurrentUser();
  if (!user || (user.role !== 'staff' && user.role !== 'admin')) {
    container.innerHTML = `
      <div style="max-width: 500px; margin: 40px auto; text-align: center;" class="alert alert-danger">
        <h3>Access Restricted</h3>
        <p>This dashboard is reserved for licensed pharmacists and pharmacy staff.</p>
        <a href="#/medicines" class="btn btn-outline" style="margin-top: 12px;">Return to Catalogue</a>
      </div>
    `;
    return;
  }

  let activeTab: 'prescriptions' | 'orders' | 'inventory' = 'prescriptions';

  container.innerHTML = `
    <div>
      <div style="margin-bottom: 24px; display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 12px;">
        <div>
          <h1 style="font-size: 1.8rem; font-weight: 800; margin-bottom: 4px;">Pharmacist &amp; Staff Operations</h1>
          <p style="color: var(--text-muted);">
            Manage clinical prescription approvals, order fulfillment workflow, and medicine inventory.
          </p>
        </div>
        <div style="display: flex; gap: 8px;">
          <button id="tab-rx" class="btn btn-primary btn-sm">Prescription Queue</button>
          <button id="tab-orders" class="btn btn-outline btn-sm">Order Fulfillment</button>
          <button id="tab-inventory" class="btn btn-outline btn-sm">Inventory Control</button>
        </div>
      </div>

      <div id="staff-tab-content"></div>
    </div>
  `;

  const tabRx = container.querySelector('#tab-rx') as HTMLButtonElement;
  const tabOrders = container.querySelector('#tab-orders') as HTMLButtonElement;
  const tabInventory = container.querySelector('#tab-inventory') as HTMLButtonElement;
  const tabContent = container.querySelector('#staff-tab-content') as HTMLDivElement;

  const updateTabButtons = () => {
    [tabRx, tabOrders, tabInventory].forEach((b) => {
      b.classList.remove('btn-primary');
      b.classList.add('btn-outline');
    });
    if (activeTab === 'prescriptions') {
      tabRx.classList.add('btn-primary');
      tabRx.classList.remove('btn-outline');
    } else if (activeTab === 'orders') {
      tabOrders.classList.add('btn-primary');
      tabOrders.classList.remove('btn-outline');
    } else if (activeTab === 'inventory') {
      tabInventory.classList.add('btn-primary');
      tabInventory.classList.remove('btn-outline');
    }
  };

  const renderTab = () => {
    updateTabButtons();
    if (activeTab === 'prescriptions') {
      renderPrescriptionQueue(tabContent);
    } else if (activeTab === 'orders') {
      renderOrderFulfillment(tabContent);
    } else if (activeTab === 'inventory') {
      renderInventoryControl(tabContent);
    }
  };

  tabRx.addEventListener('click', () => { activeTab = 'prescriptions'; renderTab(); });
  tabOrders.addEventListener('click', () => { activeTab = 'orders'; renderTab(); });
  tabInventory.addEventListener('click', () => { activeTab = 'inventory'; renderTab(); });

  renderTab();
}

async function renderPrescriptionQueue(container: HTMLElement) {
  container.innerHTML = `<div style="text-align: center; padding: 40px; color: var(--text-muted);">Loading pending prescription review queue...</div>`;

  try {
    const queue = await api.getPendingPrescriptions();

    if (queue.length === 0) {
      container.innerHTML = `
        <div style="text-align: center; padding: 60px 20px; background: white; border-radius: var(--radius-md); border: 1px dashed var(--border);">
          <div style="font-size: 2.5rem; margin-bottom: 8px;">✅</div>
          <h3>Queue is Empty</h3>
          <p style="color: var(--text-muted); font-size: 0.9rem;">There are no prescriptions currently awaiting pharmacist review.</p>
        </div>
      `;
      return;
    }

    container.innerHTML = `
      <div style="display: flex; flex-direction: column; gap: 20px;">
        ${queue.map((rx: Prescription) => `
          <div class="card" style="padding: 24px;" data-rx-id="${rx.prescription_id}">
            <div style="display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 12px; flex-wrap: wrap; gap: 8px;">
              <div>
                <h3 style="font-size: 1.15rem; margin-bottom: 4px;">Prescription #${rx.prescription_id.substring(0, 8)}</h3>
                <div style="font-size: 0.85rem; color: var(--text-muted);">
                  Uploaded: ${new Date(rx.created_at).toLocaleString()} &bull; User ID: <code>${rx.user_id.substring(0, 8)}</code>
                </div>
              </div>
              <div style="display: flex; gap: 8px; align-items: center;">
                <span class="badge ${rx.validation_status === 'VALID' ? 'badge-otc' : 'badge-warning'}">
                  F# Validation: ${rx.validation_status}
                </span>
                <a href="${api.getPrescriptionFileUrl(rx.prescription_id)}" target="_blank" class="btn btn-outline btn-sm">
                  🔍 View Document
                </a>
              </div>
            </div>

            <div style="margin-bottom: 14px;">
              <div style="font-size: 0.85rem; font-weight: 700; color: var(--text-muted); margin-bottom: 6px;">
                AI OCR Extracted Text:
              </div>
              <div class="ocr-box">${escapeHtml(rx.extracted_text || '(No text extracted or manual review required)')}</div>
            </div>

            <div class="form-group" style="margin-bottom: 14px;">
              <label style="font-size: 0.85rem; font-weight: 600;">Pharmacist Review Note (Required for Reject, Optional for Approve):</label>
              <textarea class="form-textarea rx-note-input" rows="2" placeholder="e.g. Verified doctor registration, valid dosage, clear instructions..."></textarea>
            </div>

            <div style="display: flex; gap: 12px; justify-content: flex-end;">
              <button class="btn btn-danger btn-sm reject-rx-btn" data-id="${rx.prescription_id}">
                ❌ Reject Prescription
              </button>
              <button class="btn btn-primary btn-sm approve-rx-btn" data-id="${rx.prescription_id}">
                ✓ Approve Prescription
              </button>
            </div>
          </div>
        `).join('')}
      </div>
    `;

    // Attach approve/reject handlers
    container.querySelectorAll('.approve-rx-btn').forEach((btn) => {
      btn.addEventListener('click', async (e) => {
        const id = (e.currentTarget as HTMLElement).getAttribute('data-id')!;
        const card = container.querySelector(`[data-rx-id="${id}"]`)!;
        const note = (card.querySelector('.rx-note-input') as HTMLTextAreaElement).value.trim();
        try {
          (e.currentTarget as HTMLButtonElement).disabled = true;
          await api.reviewPrescription(id, 'APPROVE', note || 'Approved by pharmacist');
          renderPrescriptionQueue(container);
        } catch (err: any) {
          alert('Error approving prescription: ' + err.message);
          (e.currentTarget as HTMLButtonElement).disabled = false;
        }
      });
    });

    container.querySelectorAll('.reject-rx-btn').forEach((btn) => {
      btn.addEventListener('click', async (e) => {
        const id = (e.currentTarget as HTMLElement).getAttribute('data-id')!;
        const card = container.querySelector(`[data-rx-id="${id}"]`)!;
        const note = (card.querySelector('.rx-note-input') as HTMLTextAreaElement).value.trim();
        if (!note) {
          alert('Please enter a clinical reason in the Pharmacist Review Note when rejecting a prescription.');
          return;
        }
        try {
          (e.currentTarget as HTMLButtonElement).disabled = true;
          await api.reviewPrescription(id, 'REJECT', note);
          renderPrescriptionQueue(container);
        } catch (err: any) {
          alert('Error rejecting prescription: ' + err.message);
          (e.currentTarget as HTMLButtonElement).disabled = false;
        }
      });
    });

  } catch (err: any) {
    container.innerHTML = `<div class="alert alert-danger">Failed to load prescription queue: ${err.message}</div>`;
  }
}

async function renderOrderFulfillment(container: HTMLElement) {
  container.innerHTML = `<div style="text-align: center; padding: 40px; color: var(--text-muted);">Loading orders...</div>`;

  try {
    const orders = await api.getAllOrders();

    if (orders.length === 0) {
      container.innerHTML = `
        <div style="text-align: center; padding: 40px; background: white; border-radius: var(--radius-md); border: 1px dashed var(--border); color: var(--text-muted);">
          No customer orders found.
        </div>
      `;
      return;
    }

    container.innerHTML = `
      <div class="table-container">
        <table>
          <thead>
            <tr>
              <th>Order ID</th>
              <th>Customer</th>
              <th>Date</th>
              <th>Total</th>
              <th>Current Status</th>
              <th>Allowed Actions</th>
            </tr>
          </thead>
          <tbody>
            ${orders.map((o: Order) => {
              const actions = getAllowedStatusButtons(o.status, o.order_id);
              return `
                <tr>
                  <td><code>#${o.order_id.substring(0, 8)}</code></td>
                  <td>
                    <strong>${escapeHtml(o.customer_name || 'Customer')}</strong>
                    <div style="font-size: 0.78rem; color: var(--text-muted);">${escapeHtml(o.customer_email || '')}</div>
                  </td>
                  <td>${new Date(o.created_at).toLocaleDateString()}</td>
                  <td style="font-weight: 700;">$${o.total.toFixed(2)}</td>
                  <td>
                    <span class="badge ${o.status === 'COMPLETED' ? 'badge-otc' : o.status === 'CANCELLED' ? 'badge-rx' : 'badge-warning'}">
                      ${o.status}
                    </span>
                  </td>
                  <td>
                    <div style="display: flex; gap: 6px; flex-wrap: wrap;">
                      ${actions}
                    </div>
                  </td>
                </tr>
              `;
            }).join('')}
          </tbody>
        </table>
      </div>
    `;

    container.querySelectorAll('.status-action-btn').forEach((btn) => {
      btn.addEventListener('click', async (e) => {
        const target = e.currentTarget as HTMLButtonElement;
        const id = target.getAttribute('data-id')!;
        const nextStatus = target.getAttribute('data-status') as OrderStatus;

        try {
          target.disabled = true;
          await api.updateOrderStatus(id, nextStatus);
          renderOrderFulfillment(container);
        } catch (err: any) {
          alert('Error transitioning order: ' + err.message);
          target.disabled = false;
        }
      });
    });

  } catch (err: any) {
    container.innerHTML = `<div class="alert alert-danger">Failed to load orders: ${err.message}</div>`;
  }
}

function getAllowedStatusButtons(current: OrderStatus, orderId: string): string {
  switch (current) {
    case 'PLACED':
      return `
        <button class="btn btn-outline btn-sm status-action-btn" data-id="${orderId}" data-status="UNDER_REVIEW">Mark Under Review</button>
        <button class="btn btn-danger btn-sm status-action-btn" data-id="${orderId}" data-status="CANCELLED">Cancel</button>
      `;
    case 'UNDER_REVIEW':
      return `
        <button class="btn btn-primary btn-sm status-action-btn" data-id="${orderId}" data-status="CONFIRMED">Confirm Order</button>
        <button class="btn btn-danger btn-sm status-action-btn" data-id="${orderId}" data-status="CANCELLED">Cancel</button>
      `;
    case 'CONFIRMED':
      return `
        <button class="btn btn-primary btn-sm status-action-btn" data-id="${orderId}" data-status="READY">Mark Ready</button>
        <button class="btn btn-danger btn-sm status-action-btn" data-id="${orderId}" data-status="CANCELLED">Cancel</button>
      `;
    case 'READY':
      return `
        <button class="btn btn-primary btn-sm status-action-btn" data-id="${orderId}" data-status="COMPLETED" style="background: var(--success);">Complete Order</button>
        <button class="btn btn-danger btn-sm status-action-btn" data-id="${orderId}" data-status="CANCELLED">Cancel</button>
      `;
    case 'COMPLETED':
      return `<span style="font-size: 0.8rem; color: var(--success); font-weight: 600;">✓ Completed</span>`;
    case 'CANCELLED':
      return `<span style="font-size: 0.8rem; color: var(--danger); font-weight: 600;">Cancelled (Stock Restored)</span>`;
    default:
      return '';
  }
}

async function renderInventoryControl(container: HTMLElement) {
  container.innerHTML = `<div style="text-align: center; padding: 40px; color: var(--text-muted);">Loading inventory...</div>`;

  try {
    const res = await api.getMedicines('', 1, 100);
    const medicines = res.data;

    container.innerHTML = `
      <div class="table-container">
        <table>
          <thead>
            <tr>
              <th>Medicine Name</th>
              <th>Generic Name</th>
              <th>Type</th>
              <th>Price ($)</th>
              <th>Stock Qty</th>
              <th>Save</th>
            </tr>
          </thead>
          <tbody>
            ${medicines.map((m: Medicine) => `
              <tr data-med-id="${m.medicine_id}">
                <td><strong>${escapeHtml(m.name)}</strong></td>
                <td>${escapeHtml(m.generic_name)}</td>
                <td>
                  <span class="badge ${m.requires_rx ? 'badge-rx' : 'badge-otc'}">
                    ${m.requires_rx ? 'Rx' : 'OTC'}
                  </span>
                </td>
                <td>
                  <input type="number" step="0.25" min="0" class="form-input price-input" value="${m.price.toFixed(2)}" style="width: 100px; padding: 4px 8px;" />
                </td>
                <td>
                  <input type="number" min="0" class="form-input stock-input" value="${m.stock_qty}" style="width: 80px; padding: 4px 8px;" />
                </td>
                <td>
                  <button class="btn btn-outline btn-sm save-inventory-btn" data-id="${m.medicine_id}">Save</button>
                </td>
              </tr>
            `).join('')}
          </tbody>
        </table>
      </div>
    `;

    container.querySelectorAll('.save-inventory-btn').forEach((btn) => {
      btn.addEventListener('click', async (e) => {
        const button = e.currentTarget as HTMLButtonElement;
        const id = button.getAttribute('data-id')!;
        const row = container.querySelector(`[data-med-id="${id}"]`)!;
        const priceInput = row.querySelector('.price-input') as HTMLInputElement;
        const stockInput = row.querySelector('.stock-input') as HTMLInputElement;

        const price = parseFloat(priceInput.value);
        const stock = parseInt(stockInput.value, 10);

        if (isNaN(price) || price < 0 || isNaN(stock) || stock < 0) {
          alert('Please enter valid non-negative numbers for price and stock.');
          return;
        }

        try {
          button.disabled = true;
          button.textContent = 'Saving...';
          await api.updateStock(id, stock, price);
          button.textContent = '✓ Saved';
          button.style.color = 'var(--success)';
          setTimeout(() => {
            button.disabled = false;
            button.textContent = 'Save';
            button.style.color = '';
          }, 1500);
        } catch (err: any) {
          alert('Failed to update inventory: ' + err.message);
          button.disabled = false;
          button.textContent = 'Save';
        }
      });
    });

  } catch (err: any) {
    container.innerHTML = `<div class="alert alert-danger">Failed to load inventory: ${err.message}</div>`;
  }
}

function escapeHtml(str: string): string {
  const div = document.createElement('div');
  div.textContent = str;
  return div.innerHTML;
}
