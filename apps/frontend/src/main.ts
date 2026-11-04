import { api } from './api';
import { cart } from './cart';
import { renderAuthPage } from './pages/auth';
import { renderMedicinesPage } from './pages/medicines';
import { renderCartPage } from './pages/cart';
import { renderPrescriptionsPage } from './pages/prescriptions';
import { renderOrdersPage } from './pages/orders';
import { renderStaffPage } from './pages/staff';
import { renderAdminPage } from './pages/admin';

const app = document.getElementById('app') as HTMLElement;
const navLinksContainer = document.getElementById('nav-links') as HTMLElement;
const userMenuContainer = document.getElementById('user-menu') as HTMLElement;

function updateNavigation() {
  const user = api.getCurrentUser();
  const cartCount = cart.getItemCount();
  const currentHash = window.location.hash || '#/medicines';

  let linksHtml = `
    <li><a href="#/medicines" class="${currentHash.startsWith('#/medicines') ? 'active' : ''}">Medicines</a></li>
  `;

  if (user) {
    linksHtml += `
      <li><a href="#/prescriptions" class="${currentHash.startsWith('#/prescriptions') ? 'active' : ''}">Prescriptions</a></li>
      <li><a href="#/orders" class="${currentHash.startsWith('#/orders') ? 'active' : ''}">My Orders</a></li>
    `;

    if (user.role === 'staff' || user.role === 'admin') {
      linksHtml += `
        <li><a href="#/staff" class="${currentHash.startsWith('#/staff') ? 'active' : ''}" style="color: #0284c7; font-weight: 700;">Staff Dashboard</a></li>
      `;
    }

    if (user.role === 'admin') {
      linksHtml += `
        <li><a href="#/admin" class="${currentHash.startsWith('#/admin') ? 'active' : ''}" style="color: #9d174d; font-weight: 700;">Admin</a></li>
      `;
    }
  }

  linksHtml += `
    <li>
      <a href="#/cart" class="${currentHash.startsWith('#/cart') ? 'active' : ''}">
        Cart <span class="badge-cart">${cartCount}</span>
      </a>
    </li>
  `;

  navLinksContainer.innerHTML = linksHtml;

  if (user) {
    userMenuContainer.innerHTML = `
      <div style="display: flex; align-items: center; gap: 8px;">
        <span style="font-size: 0.9rem; font-weight: 600;">${escapeHtml(user.name)}</span>
        <span class="role-badge ${user.role}">${user.role}</span>
        <button id="logout-btn" class="btn btn-outline btn-sm" style="margin-left: 6px;">Sign Out</button>
      </div>
    `;

    document.getElementById('logout-btn')?.addEventListener('click', () => {
      api.logout();
      window.dispatchEvent(new CustomEvent('auth-changed'));
      window.location.hash = '#/medicines';
    });
  } else {
    userMenuContainer.innerHTML = `
      <div style="display: flex; gap: 8px;">
        <a href="#/login" class="btn btn-outline btn-sm">Sign In</a>
        <a href="#/register" class="btn btn-primary btn-sm">Register</a>
      </div>
    `;
  }
}

function handleRoute() {
  const hash = window.location.hash || '#/medicines';
  updateNavigation();

  if (hash.startsWith('#/login')) {
    renderAuthPage(app, false);
  } else if (hash.startsWith('#/register')) {
    renderAuthPage(app, true);
  } else if (hash.startsWith('#/cart')) {
    renderCartPage(app);
  } else if (hash.startsWith('#/prescriptions')) {
    renderPrescriptionsPage(app);
  } else if (hash.startsWith('#/orders')) {
    renderOrdersPage(app);
  } else if (hash.startsWith('#/staff')) {
    renderStaffPage(app);
  } else if (hash.startsWith('#/admin')) {
    renderAdminPage(app);
  } else {
    renderMedicinesPage(app);
  }
}

function escapeHtml(str: string): string {
  const div = document.createElement('div');
  div.textContent = str;
  return div.innerHTML;
}

// Event Listeners
window.addEventListener('hashchange', handleRoute);
window.addEventListener('auth-changed', () => {
  updateNavigation();
  handleRoute();
});
cart.subscribe(updateNavigation);

// Initial start
handleRoute();
