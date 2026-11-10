import { api } from '../api';
import { cart } from '../cart';
import { Order, Prescription } from '../types';

export function renderDashboardPage(container: HTMLElement) {
  const user = api.getCurrentUser();
  if (!user) {
    window.location.hash = '#/login';
    return;
  }

  container.innerHTML = `
    <div>
      <!-- Welcome Header -->
      <div style="background: white; border-radius: var(--radius-lg); padding: 24px 28px; border: 1px solid var(--border); box-shadow: var(--shadow); margin-bottom: 28px; display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 16px;">
        <div>
          <div style="display: flex; align-items: center; gap: 10px; margin-bottom: 6px;">
            <h1 style="font-size: 1.8rem; font-weight: 800; color: var(--text-main); margin: 0;">Welcome back, ${escapeHtml(user.name)}!</h1>
            <span class="role-badge ${user.role}">${user.role}</span>
          </div>
          <p style="color: var(--text-muted); font-size: 0.95rem; margin: 0;">
            ${escapeHtml(user.email)} &bull; Manage your prescriptions, track active orders, and quickly refill medications.
          </p>
        </div>
        <div style="display: flex; gap: 10px;">
          <a href="#/medicines" class="btn btn-primary btn-sm">Browse Catalog</a>
          <a href="#/prescriptions" class="btn btn-outline btn-sm">Upload Prescription</a>
        </div>
      </div>

      <!-- Quick Stats Row -->
      <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 16px; margin-bottom: 28px;" id="stats-row">
        <div class="card" style="padding: 16px 20px; text-align: center;">
          <div style="font-size: 0.85rem; font-weight: 600; color: var(--text-muted); text-transform: uppercase;">Active Orders</div>
          <div id="stat-active-orders" style="font-size: 2rem; font-weight: 800; color: var(--primary); margin: 6px 0;">--</div>
          <div style="font-size: 0.8rem; color: var(--text-muted);">In progress or placed</div>
        </div>
        <div class="card" style="padding: 16px 20px; text-align: center;">
          <div style="font-size: 0.85rem; font-weight: 600; color: var(--text-muted); text-transform: uppercase;">Approved Prescriptions</div>
          <div id="stat-approved-rx" style="font-size: 2rem; font-weight: 800; color: var(--success); margin: 6px 0;">--</div>
          <div style="font-size: 0.8rem; color: var(--text-muted);">Verified by pharmacist</div>
        </div>
        <div class="card" style="padding: 16px 20px; text-align: center;">
          <div style="font-size: 0.85rem; font-weight: 600; color: var(--text-muted); text-transform: uppercase;">Pending Prescriptions</div>
          <div id="stat-pending-rx" style="font-size: 2rem; font-weight: 800; color: var(--warning); margin: 6px 0;">--</div>
          <div style="font-size: 0.8rem; color: var(--text-muted);">In clinical review queue</div>
        </div>
        <div class="card" style="padding: 16px 20px; text-align: center;">
          <div style="font-size: 0.85rem; font-weight: 600; color: var(--text-muted); text-transform: uppercase;">Cart Items</div>
          <div id="stat-cart-items" style="font-size: 2rem; font-weight: 800; color: var(--secondary); margin: 6px 0;">${cart.getItemCount()}</div>
          <div style="font-size: 0.8rem; color: var(--text-muted);"><a href="#/cart" style="color: var(--secondary); text-decoration: none; font-weight: 600;">Go to Cart &rarr;</a></div>
        </div>
      </div>

      <!-- Main Dashboard Grid -->
      <div style="display: grid; grid-template-columns: 3fr 2fr; gap: 24px; align-items: start;" id="dashboard-grid">
        <!-- Left: Active Orders Tracker -->
        <div class="card" style="padding: 24px;">
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px;">
            <h3 style="font-size: 1.25rem; font-weight: 700;">Recent & Active Orders</h3>
            <a href="#/orders" style="font-size: 0.88rem; color: var(--primary); text-decoration: none; font-weight: 600;">View All Orders &rarr;</a>
          </div>

          <div id="dashboard-orders-list">
            <div style="text-align: center; padding: 30px; color: var(--text-muted);">Loading orders...</div>
          </div>
        </div>

        <!-- Right: Prescriptions Summary -->
        <div class="card" style="padding: 24px;">
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px;">
            <h3 style="font-size: 1.25rem; font-weight: 700;">Prescription Status</h3>
            <a href="#/prescriptions" style="font-size: 0.88rem; color: var(--primary); text-decoration: none; font-weight: 600;">Manage &rarr;</a>
          </div>

          <div id="dashboard-rx-list">
            <div style="text-align: center; padding: 30px; color: var(--text-muted);">Loading prescriptions...</div>
          </div>
        </div>
      </div>
    </div>
  `;

  // Fetch data
  loadDashboardData(container);
}

async function loadDashboardData(container: HTMLElement) {
  const ordersContainer = container.querySelector('#dashboard-orders-list') as HTMLDivElement;
  const rxContainer = container.querySelector('#dashboard-rx-list') as HTMLDivElement;
  const statActive = container.querySelector('#stat-active-orders') as HTMLDivElement;
  const statApproved = container.querySelector('#stat-approved-rx') as HTMLDivElement;
  const statPending = container.querySelector('#stat-pending-rx') as HTMLDivElement;

  try {
    const [orders, prescriptions] = await Promise.all([
      api.getMyOrders().catch(() => [] as Order[]),
      api.getMyPrescriptions().catch(() => [] as Prescription[]),
    ]);

    // Calculate stats
    const activeOrders = orders.filter(
      (o) => o.status === 'PLACED' || o.status === 'UNDER_REVIEW' || o.status === 'CONFIRMED' || o.status === 'READY'
    );
    const approvedRx = prescriptions.filter((p) => p.review_status === 'APPROVED');
    const pendingRx = prescriptions.filter((p) => p.review_status === 'PENDING');

    if (statActive) statActive.textContent = activeOrders.length.toString();
    if (statApproved) statApproved.textContent = approvedRx.length.toString();
    if (statPending) statPending.textContent = pendingRx.length.toString();

    // Render Orders
    if (orders.length === 0) {
      ordersContainer.innerHTML = `
        <div style="text-align: center; padding: 40px 20px; background: #f8fafc; border-radius: var(--radius-md); border: 1px dashed var(--border);">
          <div style="font-size: 2.2rem; margin-bottom: 8px;">📦</div>
          <h4 style="margin-bottom: 6px;">No orders placed yet</h4>
          <p style="color: var(--text-muted); font-size: 0.88rem; margin-bottom: 16px;">
            Explore our medicine catalogue and order over-the-counter or prescription items.
          </p>
          <a href="#/medicines" class="btn btn-primary btn-sm">Start Shopping</a>
        </div>
      `;
    } else {
      const recentOrders = orders.slice(0, 3);
      ordersContainer.innerHTML = recentOrders
        .map((order) => {
          const statusClass =
            order.status === 'COMPLETED'
              ? 'badge-success'
              : order.status === 'CANCELLED'
              ? 'badge-danger'
              : order.status === 'UNDER_REVIEW'
              ? 'badge-warning'
              : 'badge-primary';

          const dateStr = new Date(order.created_at).toLocaleDateString(undefined, {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
          });

          const itemsSummary = (order.items || [])
            .map((item) => `${escapeHtml(item.medicine_name || 'Item')} &times; ${item.quantity}`)
            .join(', ');

          return `
            <div style="background: #f8fafc; border-radius: var(--radius-md); border: 1px solid var(--border); padding: 16px; margin-bottom: 14px;">
              <div style="display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 10px;">
                <div>
                  <div style="font-weight: 700; font-size: 0.95rem; color: var(--text-main);">
                    Order #${order.order_id.substring(0, 8)}
                  </div>
                  <div style="font-size: 0.8rem; color: var(--text-muted);">${dateStr}</div>
                </div>
                <div style="text-align: right;">
                  <span class="badge ${statusClass}">${order.status.replace('_', ' ')}</span>
                  <div style="font-weight: 700; font-size: 1.1rem; color: var(--primary); margin-top: 4px;">
                    ₹${order.total.toFixed(2)}
                  </div>
                </div>
              </div>

              <div style="font-size: 0.85rem; color: var(--text-main); margin-bottom: 12px; line-height: 1.4;">
                <span style="color: var(--text-muted); font-weight: 600;">Items:</span> ${itemsSummary || 'Standard order items'}
              </div>

              <!-- Mini Timeline -->
              ${renderMiniStepper(order.status)}

              <div style="display: flex; justify-content: flex-end; gap: 8px; margin-top: 12px;">
                <button
                  class="btn btn-outline btn-sm reorder-btn"
                  data-order='${JSON.stringify(order).replace(/'/g, "&#39;")}'
                >
                  ↻ Reorder Items
                </button>
                <a href="#/orders" class="btn btn-primary btn-sm">Full Tracker</a>
              </div>
            </div>
          `;
        })
        .join('');

      // Wire reorder handlers
      ordersContainer.querySelectorAll('.reorder-btn').forEach((btn) => {
        btn.addEventListener('click', (e) => {
          const target = e.currentTarget as HTMLButtonElement;
          const orderData: Order = JSON.parse(target.getAttribute('data-order') || '{}');
          if (orderData.items && orderData.items.length > 0) {
            orderData.items.forEach((item) => {
              cart.addItem(
                {
                  medicine_id: item.medicine_id,
                  name: item.medicine_name || 'Medicine',
                  generic_name: '',
                  description: '',
                  price: item.unit_price,
                  stock_qty: 99,
                  requires_rx: !!item.requires_rx,
                },
                item.quantity
              );
            });
            target.textContent = '✓ Added to Cart!';
            target.style.background = 'var(--success)';
            target.style.color = 'white';
            setTimeout(() => {
              window.location.hash = '#/cart';
            }, 700);
          }
        });
      });
    }

    // Render Prescriptions
    if (prescriptions.length === 0) {
      rxContainer.innerHTML = `
        <div style="text-align: center; padding: 40px 20px; background: #f8fafc; border-radius: var(--radius-md); border: 1px dashed var(--border);">
          <div style="font-size: 2.2rem; margin-bottom: 8px;">📋</div>
          <h4 style="margin-bottom: 6px;">No prescriptions uploaded</h4>
          <p style="color: var(--text-muted); font-size: 0.88rem; margin-bottom: 16px;">
            Upload your prescription image or PDF to purchase Rx-restricted medicines.
          </p>
          <a href="#/prescriptions" class="btn btn-primary btn-sm">Upload Rx</a>
        </div>
      `;
    } else {
      const recentRx = prescriptions.slice(0, 4);
      rxContainer.innerHTML = recentRx
        .map((rx) => {
          const isApproved = rx.review_status === 'APPROVED';
          const isRejected = rx.review_status === 'REJECTED';
          const badgeClass = isApproved ? 'badge-success' : isRejected ? 'badge-danger' : 'badge-warning';

          const dateStr = new Date(rx.created_at).toLocaleDateString(undefined, {
            month: 'short',
            day: 'numeric',
            year: 'numeric',
          });

          return `
            <div style="background: #f8fafc; border-radius: var(--radius-md); border: 1px solid var(--border); padding: 14px; margin-bottom: 12px;">
              <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 6px;">
                <div style="font-weight: 700; font-size: 0.9rem;">
                  Prescription #${rx.prescription_id.substring(0, 8)}
                </div>
                <span class="badge ${badgeClass}">${rx.review_status}</span>
              </div>
              <div style="font-size: 0.8rem; color: var(--text-muted); margin-bottom: 8px;">
                Uploaded on ${dateStr} &bull; OCR: <strong>${rx.validation_status}</strong>
              </div>
              ${
                rx.review_note
                  ? `
                <div style="font-size: 0.8rem; background: white; padding: 6px 10px; border-radius: var(--radius-sm); border: 1px solid var(--border); color: var(--text-main); margin-bottom: 6px;">
                  <span style="font-weight: 600; color: var(--text-muted);">Pharmacist Note:</span> ${escapeHtml(rx.review_note)}
                </div>
              `
                  : ''
              }
              <div style="display: flex; justify-content: space-between; align-items: center; margin-top: 8px;">
                <span style="font-size: 0.78rem; color: ${isApproved ? 'var(--success)' : 'var(--warning)'}; font-weight: 600;">
                  ${isApproved ? '✓ Ready to use for checkout' : '⏳ Awaiting pharmacist approval'}
                </span>
                <a href="#/prescriptions" style="font-size: 0.8rem; color: var(--primary); font-weight: 600; text-decoration: none;">View &rarr;</a>
              </div>
            </div>
          `;
        })
        .join('');
    }
  } catch (err: any) {
    ordersContainer.innerHTML = `<div class="alert alert-danger">Failed to load dashboard data: ${err.message}</div>`;
  }
}

function renderMiniStepper(status: string): string {
  const steps = ['PLACED', 'UNDER_REVIEW', 'CONFIRMED', 'READY', 'COMPLETED'];
  const stepLabels: Record<string, string> = {
    PLACED: 'Placed',
    UNDER_REVIEW: 'Review',
    CONFIRMED: 'Approved',
    READY: 'Ready',
    COMPLETED: 'Delivered',
  };

  if (status === 'CANCELLED') {
    return `
      <div style="font-size: 0.8rem; color: var(--danger); font-weight: 600; background: var(--danger-light); padding: 4px 8px; border-radius: var(--radius-sm); text-align: center;">
        ✕ This order was cancelled
      </div>
    `;
  }

  const currentIdx = steps.indexOf(status);

  return `
    <div style="display: flex; justify-content: space-between; align-items: center; position: relative; margin: 12px 0 6px;">
      ${steps
        .map((step, idx) => {
          const isDone = currentIdx >= idx;
          const isCurrent = currentIdx === idx;
          const color = isDone ? 'var(--success)' : 'var(--border)';
          const textColor = isCurrent ? 'var(--primary)' : isDone ? 'var(--text-main)' : 'var(--text-muted)';
          return `
            <div style="text-align: center; flex: 1;">
              <div style="width: 14px; height: 14px; border-radius: 50%; background: ${color}; margin: 0 auto 4px; border: 2px solid white; box-shadow: 0 0 0 1px ${color};"></div>
              <div style="font-size: 0.7rem; font-weight: ${isCurrent ? '700' : '500'}; color: ${textColor};">${stepLabels[step]}</div>
            </div>
          `;
        })
        .join('')}
    </div>
  `;
}

function escapeHtml(str: string): string {
  const div = document.createElement('div');
  div.textContent = str;
  return div.innerHTML;
}
