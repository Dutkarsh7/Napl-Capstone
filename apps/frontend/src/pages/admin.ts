import { api } from '../api';
import { User, UserRole } from '../types';

export async function renderAdminPage(container: HTMLElement) {
  const user = api.getCurrentUser();
  if (!user || user.role !== 'admin') {
    container.innerHTML = `
      <div style="max-width: 500px; margin: 40px auto; text-align: center;" class="alert alert-danger">
        <h3>Admin Access Only</h3>
        <p>This panel is restricted to system administrators.</p>
        <a href="#/medicines" class="btn btn-outline" style="margin-top: 12px;">Return to Catalogue</a>
      </div>
    `;
    return;
  }

  container.innerHTML = `
    <div style="max-width: 900px; margin: 0 auto;">
      <div style="margin-bottom: 24px;">
        <h1 style="font-size: 1.8rem; font-weight: 800; margin-bottom: 4px;">User &amp; Role Administration</h1>
        <p style="color: var(--text-muted);">
          Manage user accounts and assign staff or administrator privileges.
        </p>
      </div>

      <div id="users-table-container">
        <div style="text-align: center; padding: 40px; color: var(--text-muted);">Loading users...</div>
      </div>
    </div>
  `;

  const tableContainer = container.querySelector('#users-table-container') as HTMLDivElement;

  try {
    const users = await api.getUsers();

    tableContainer.innerHTML = `
      <div class="table-container">
        <table>
          <thead>
            <tr>
              <th>User ID</th>
              <th>Full Name</th>
              <th>Email Address</th>
              <th>Current Role</th>
              <th>Update Role</th>
            </tr>
          </thead>
          <tbody>
            ${users.map((u: User) => `
              <tr data-user-id="${u.user_id}">
                <td><code>${u.user_id.substring(0, 8)}...</code></td>
                <td><strong>${escapeHtml(u.name)}</strong></td>
                <td>${escapeHtml(u.email)}</td>
                <td>
                  <span class="role-badge ${u.role}">
                    ${u.role}
                  </span>
                </td>
                <td>
                  <div style="display: flex; gap: 8px; align-items: center;">
                    <select class="form-select role-select" style="padding: 4px 8px; width: 130px; font-size: 0.85rem;">
                      <option value="customer" ${u.role === 'customer' ? 'selected' : ''}>Customer</option>
                      <option value="staff" ${u.role === 'staff' ? 'selected' : ''}>Staff (Pharmacist)</option>
                      <option value="admin" ${u.role === 'admin' ? 'selected' : ''}>Admin</option>
                    </select>
                    <button class="btn btn-outline btn-sm apply-role-btn" data-id="${u.user_id}">
                      Apply
                    </button>
                  </div>
                </td>
              </tr>
            `).join('')}
          </tbody>
        </table>
      </div>
    `;

    tableContainer.querySelectorAll('.apply-role-btn').forEach((btn) => {
      btn.addEventListener('click', async (e) => {
        const button = e.currentTarget as HTMLButtonElement;
        const id = button.getAttribute('data-id')!;
        const row = tableContainer.querySelector(`[data-user-id="${id}"]`)!;
        const select = row.querySelector('.role-select') as HTMLSelectElement;
        const targetRole = select.value as UserRole;

        try {
          button.disabled = true;
          button.textContent = 'Updating...';
          await api.updateUserRole(id, targetRole);
          button.textContent = '✓ Updated';
          button.style.color = 'var(--success)';
          setTimeout(() => {
            renderAdminPage(container);
          }, 1000);
        } catch (err: any) {
          alert('Failed to update role: ' + err.message);
          button.disabled = false;
          button.textContent = 'Apply';
        }
      });
    });

  } catch (err: any) {
    tableContainer.innerHTML = `<div class="alert alert-danger">Failed to load users: ${err.message}</div>`;
  }
}

function escapeHtml(str: string): string {
  const div = document.createElement('div');
  div.textContent = str;
  return div.innerHTML;
}
