import { api } from '../api';
import { cart } from '../cart';
import { Medicine } from '../types';

export function renderMedicinesPage(container: HTMLElement) {
  let currentSearch = '';
  let currentFilter = 'all'; // 'all', 'otc', 'rx'
  let currentPage = 1;
  const pageSize = 12;
  let debounceTimer: any = null;

  container.innerHTML = `
    <div>
      <div style="margin-bottom: 24px;">
        <h1 style="font-size: 1.8rem; font-weight: 800; margin-bottom: 6px;">Medicine Catalogue</h1>
        <p style="color: var(--text-muted);">
          Search verified pharmaceuticals and over-the-counter remedies. Prescription items require pharmacist authorization.
        </p>
      </div>

      <!-- Controls: Search & Category Filter -->
      <div class="search-container">
        <input
          type="text"
          id="med-search-input"
          class="search-input"
          placeholder="Search by trade name or active generic substance (e.g. Amoxicillin, Paracetamol, Ibuprofen)..."
        />

        <div style="display: flex; gap: 8px; align-items: center; flex-wrap: wrap;">
          <button id="filter-all" class="btn btn-primary btn-sm">All Medicines</button>
          <button id="filter-otc" class="btn btn-outline btn-sm">OTC Only</button>
          <button id="filter-rx" class="btn btn-outline btn-sm">Rx Only</button>
        </div>
      </div>

      <!-- Medicine Cards Grid Container -->
      <div id="medicines-grid" class="grid">
        <div style="grid-column: 1 / -1; text-align: center; padding: 40px; color: var(--text-muted);">
          Loading catalogue...
        </div>
      </div>

      <!-- Pagination -->
      <div id="pagination-controls" style="display: flex; justify-content: center; align-items: center; gap: 16px; margin-top: 32px;">
        <button id="prev-page-btn" class="btn btn-outline btn-sm" disabled>&larr; Previous</button>
        <span id="page-indicator" style="font-size: 0.9rem; font-weight: 600; color: var(--text-muted);">Page 1</span>
        <button id="next-page-btn" class="btn btn-outline btn-sm" disabled>Next &rarr;</button>
      </div>
    </div>
  `;

  const searchInput = container.querySelector('#med-search-input') as HTMLInputElement;
  const grid = container.querySelector('#medicines-grid') as HTMLDivElement;
  const prevBtn = container.querySelector('#prev-page-btn') as HTMLButtonElement;
  const nextBtn = container.querySelector('#next-page-btn') as HTMLButtonElement;
  const pageIndicator = container.querySelector('#page-indicator') as HTMLSpanElement;

  const btnAll = container.querySelector('#filter-all') as HTMLButtonElement;
  const btnOtc = container.querySelector('#filter-otc') as HTMLButtonElement;
  const btnRx = container.querySelector('#filter-rx') as HTMLButtonElement;

  const updateFilterStyles = () => {
    [btnAll, btnOtc, btnRx].forEach((b) => {
      b.classList.remove('btn-primary');
      b.classList.add('btn-outline');
    });
    if (currentFilter === 'all') {
      btnAll.classList.add('btn-primary');
      btnAll.classList.remove('btn-outline');
    } else if (currentFilter === 'otc') {
      btnOtc.classList.add('btn-primary');
      btnOtc.classList.remove('btn-outline');
    } else if (currentFilter === 'rx') {
      btnRx.classList.add('btn-primary');
      btnRx.classList.remove('btn-outline');
    }
  };

  const loadMedicines = async () => {
    grid.innerHTML = `
      <div style="grid-column: 1 / -1; text-align: center; padding: 40px; color: var(--text-muted);">
        Fetching medicines...
      </div>
    `;

    try {
      const res = await api.getMedicines(currentSearch, currentPage, pageSize);
      let items = res.data;

      // Filter client-side if tab is selected
      if (currentFilter === 'otc') {
        items = items.filter((m) => !m.requires_rx);
      } else if (currentFilter === 'rx') {
        items = items.filter((m) => m.requires_rx);
      }

      if (items.length === 0) {
        grid.innerHTML = `
          <div style="grid-column: 1 / -1; text-align: center; padding: 60px 20px; background: white; border-radius: var(--radius-md); border: 1px dashed var(--border);">
            <div style="font-size: 2rem; margin-bottom: 8px;">🔍</div>
            <h3 style="margin-bottom: 6px;">No medicines found</h3>
            <p style="color: var(--text-muted); font-size: 0.9rem;">Try adjusting your search terms or filter selection.</p>
          </div>
        `;
        prevBtn.disabled = true;
        nextBtn.disabled = true;
        return;
      }

      grid.innerHTML = items
        .map((med: Medicine) => {
          const isOutOfStock = med.stock_qty <= 0;
          const isLowStock = med.stock_qty > 0 && med.stock_qty <= 10;
          const stockClass = isOutOfStock ? 'out-of-stock' : isLowStock ? 'low-stock' : 'in-stock';
          const stockLabel = isOutOfStock ? 'Out of Stock' : isLowStock ? `Low Stock (${med.stock_qty})` : `In Stock (${med.stock_qty})`;

          return `
            <div class="card" data-med-id="${med.medicine_id}">
              <div class="card-header">
                <span class="badge ${med.requires_rx ? 'badge-rx' : 'badge-otc'}">
                  ${med.requires_rx ? 'Rx Required' : 'OTC'}
                </span>
                <span class="badge-stock ${stockClass}">
                  ${stockLabel}
                </span>
              </div>

              <div class="card-title">${escapeHtml(med.name)}</div>
              <div class="card-subtitle">Active: ${escapeHtml(med.generic_name)}</div>
              <div class="card-body">
                ${escapeHtml(med.description || 'Standard pharmaceutical preparation.')}
              </div>

              <div class="card-footer">
                <div class="price">$${med.price.toFixed(2)}</div>
                <button
                  class="btn btn-primary btn-sm add-cart-btn"
                  data-med='${JSON.stringify(med).replace(/'/g, "&#39;")}'
                  ${isOutOfStock ? 'disabled' : ''}
                >
                  ${isOutOfStock ? 'Unavailable' : 'Add to Cart'}
                </button>
              </div>
            </div>
          `;
        })
        .join('');

      // Attach add-to-cart handlers
      grid.querySelectorAll('.add-cart-btn').forEach((button) => {
        button.addEventListener('click', (e) => {
          const target = e.currentTarget as HTMLButtonElement;
          const medData = JSON.parse(target.getAttribute('data-med') || '{}');
          cart.addItem(medData, 1);
          target.textContent = '✓ Added!';
          target.style.background = 'var(--success)';
          setTimeout(() => {
            target.textContent = 'Add to Cart';
            target.style.background = '';
          }, 1200);
        });
      });

      // Update pagination
      const totalPages = Math.ceil(res.total / pageSize) || 1;
      pageIndicator.textContent = `Page ${currentPage} of ${totalPages} (${res.total} total items)`;
      prevBtn.disabled = currentPage <= 1;
      nextBtn.disabled = currentPage >= totalPages;
    } catch (err: any) {
      grid.innerHTML = `
        <div style="grid-column: 1 / -1; padding: 20px;" class="alert alert-danger">
          Failed to load medicines: ${err.message}
        </div>
      `;
    }
  };

  // Search input with debounce
  searchInput.addEventListener('input', () => {
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => {
      currentSearch = searchInput.value.trim();
      currentPage = 1;
      loadMedicines();
    }, 300);
  });

  // Filter chips
  btnAll.addEventListener('click', () => {
    currentFilter = 'all';
    updateFilterStyles();
    currentPage = 1;
    loadMedicines();
  });

  btnOtc.addEventListener('click', () => {
    currentFilter = 'otc';
    updateFilterStyles();
    currentPage = 1;
    loadMedicines();
  });

  btnRx.addEventListener('click', () => {
    currentFilter = 'rx';
    updateFilterStyles();
    currentPage = 1;
    loadMedicines();
  });

  prevBtn.addEventListener('click', () => {
    if (currentPage > 1) {
      currentPage--;
      loadMedicines();
    }
  });

  nextBtn.addEventListener('click', () => {
    currentPage++;
    loadMedicines();
  });

  // Initial load
  loadMedicines();
}

function escapeHtml(str: string): string {
  const div = document.createElement('div');
  div.textContent = str;
  return div.innerHTML;
}
