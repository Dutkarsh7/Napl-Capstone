import { api } from '../api';
import { Order, OrderStatus } from '../types';

export async function renderOrdersPage(container: HTMLElement) {
  const user = api.getCurrentUser();
  if (!user) {
    window.location.hash = '#/login';
    return;
  }

  container.innerHTML = `
    <div style="max-width: 900px; margin: 0 auto;">
      <div style="margin-bottom: 24px;">
        <h1 style="font-size: 1.8rem; font-weight: 800; margin-bottom: 6px;">My Order History</h1>
        <p style="color: var(--text-muted);">
          Track the live fulfillment status and pharmacist verification stages of your pharmacy orders.
        </p>
      </div>

      <div id="orders-list">
        <div style="text-align: center; padding: 40px; color: var(--text-muted);">Loading orders...</div>
      </div>
    </div>
  `;

  const listContainer = container.querySelector('#orders-list') as HTMLDivElement;

  try {
    const orders = await api.getMyOrders();

    if (orders.length === 0) {
      listContainer.innerHTML = `
        <div style="text-align: center; padding: 60px 20px; background: white; border-radius: var(--radius-md); border: 1px dashed var(--border);">
          <div style="font-size: 3rem; margin-bottom: 12px;">📦</div>
          <h3>No Orders Placed Yet</h3>
          <p style="color: var(--text-muted); margin: 8px 0 20px;">Browse our catalog to add items to your cart.</p>
          <a href="#/medicines" class="btn btn-primary">Browse Catalogue</a>
        </div>
      `;
      return;
    }

    listContainer.innerHTML = orders.map((order: Order) => renderOrderCard(order)).join('');
  } catch (err: any) {
    listContainer.innerHTML = `
      <div class="alert alert-danger">
        Failed to load orders: ${err.message}
      </div>
    `;
  }
}

function renderOrderCard(order: Order): string {
  const isCancelled = order.status === 'CANCELLED';
  const steps: OrderStatus[] = ['PLACED', 'UNDER_REVIEW', 'CONFIRMED', 'READY', 'COMPLETED'];
  const currentIndex = steps.indexOf(order.status);

  return `
    <div class="card" style="padding: 24px; margin-bottom: 24px;">
      <div style="display: flex; justify-content: space-between; align-items: flex-start; flex-wrap: wrap; gap: 8px; border-bottom: 1px solid var(--border); padding-bottom: 14px; margin-bottom: 16px;">
        <div>
          <h3 style="font-size: 1.15rem; margin-bottom: 4px;">
            Order #${order.order_id.substring(0, 8)}
          </h3>
          <div style="font-size: 0.85rem; color: var(--text-muted);">
            Placed on ${new Date(order.created_at).toLocaleString()}
          </div>
        </div>
        <div style="text-align: right;">
          <div class="price" style="margin-bottom: 4px;">$${order.total.toFixed(2)}</div>
          <span class="badge ${getStatusBadgeClass(order.status)}">${order.status}</span>
        </div>
      </div>

      <!-- Status Timeline Stepper -->
      ${isCancelled ? `
        <div class="alert alert-danger" style="margin: 16px 0; text-align: center; font-weight: 700;">
          ❌ This order was CANCELLED. All reserved stock has been safely restored.
        </div>
      ` : `
        <div class="timeline">
          ${steps.map((st, idx) => {
            const isCompleted = idx < currentIndex;
            const isActive = idx === currentIndex;
            const cls = isCompleted ? 'completed' : isActive ? 'active' : '';
            return `
              <div class="timeline-step ${cls}">
                <div class="step-node">${isCompleted ? '✓' : idx + 1}</div>
                <div class="step-label">${formatStatusLabel(st)}</div>
              </div>
            `;
          }).join('')}
        </div>
      `}

      <!-- Items Table -->
      <div style="margin-top: 16px;">
        <h4 style="font-size: 0.95rem; margin-bottom: 10px;">Order Items:</h4>
        <div class="table-container">
          <table>
            <thead>
              <tr>
                <th>Medicine</th>
                <th>Qty</th>
                <th>Unit Price</th>
                <th>Total</th>
              </tr>
            </thead>
            <tbody>
              ${(order.items || []).map((it) => `
                <tr>
                  <td>
                    <strong>${escapeHtml(it.medicine_name || it.medicine_id)}</strong>
                    ${it.requires_rx ? '<span class="badge badge-rx" style="margin-left: 6px;">Rx</span>' : ''}
                  </td>
                  <td>${it.quantity}</td>
                  <td>$${it.unit_price.toFixed(2)}</td>
                  <td>$${(it.unit_price * it.quantity).toFixed(2)}</td>
                </tr>
              `).join('')}
            </tbody>
          </table>
        </div>
      </div>

      ${order.prescription_id ? `
        <div style="margin-top: 12px; font-size: 0.85rem; color: var(--text-muted);">
          📎 Attached Prescription: <code>${order.prescription_id.substring(0, 8)}...</code>
        </div>
      ` : ''}
    </div>
  `;
}

function getStatusBadgeClass(status: OrderStatus): string {
  switch (status) {
    case 'COMPLETED':
      return 'badge-otc';
    case 'CANCELLED':
      return 'badge-rx';
    case 'READY':
      return 'badge-otc';
    default:
      return 'badge-warning';
  }
}

function formatStatusLabel(status: OrderStatus): string {
  switch (status) {
    case 'PLACED':
      return 'Placed';
    case 'UNDER_REVIEW':
      return 'Under Review';
    case 'CONFIRMED':
      return 'Confirmed';
    case 'READY':
      return 'Ready';
    case 'COMPLETED':
      return 'Completed';
    default:
      return status;
  }
}

function escapeHtml(str: string): string {
  const div = document.createElement('div');
  div.textContent = str;
  return div.innerHTML;
}
