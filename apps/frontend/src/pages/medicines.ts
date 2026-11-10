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

  const showMedicineModal = (med: Medicine) => {
    const isOutOfStock = med.stock_qty <= 0;
    const isLowStock = med.stock_qty > 0 && med.stock_qty <= 10;
    const stockClass = isOutOfStock ? 'out-of-stock' : isLowStock ? 'low-stock' : 'in-stock';
    const stockLabel = isOutOfStock ? 'Out of Stock' : isLowStock ? `Low Stock (${med.stock_qty} left)` : `In Stock (${med.stock_qty} units)`;

    const modalOverlay = document.createElement('div');
    modalOverlay.className = 'modal-overlay';
    modalOverlay.id = 'med-details-modal';
    modalOverlay.innerHTML = `
      <div class="modal-content" style="max-width: 520px; border-top: 4px solid ${med.requires_rx ? 'var(--warning)' : 'var(--primary)'};">
        <div style="display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 16px;">
          <div>
            <div style="display: flex; gap: 8px; align-items: center; margin-bottom: 6px;">
              <span class="badge ${med.requires_rx ? 'badge-rx' : 'badge-otc'}">
                ${med.requires_rx ? 'Rx Required' : 'OTC'}
              </span>
              <span class="badge-stock ${stockClass}">
                ${stockLabel}
              </span>
            </div>
            <h2 style="font-size: 1.4rem; font-weight: 800; color: var(--text-main); margin: 0;">${escapeHtml(med.name)}</h2>
            <div style="color: var(--text-muted); font-size: 0.9rem; margin-top: 2px;">
              Active Generic: <strong>${escapeHtml(med.generic_name)}</strong>
            </div>
          </div>
          <button id="close-med-modal" class="btn btn-outline btn-sm" style="padding: 4px 10px; font-size: 1.1rem; line-height: 1;">&times;</button>
        </div>

        <div style="margin-bottom: 20px;">
          <h4 style="font-size: 0.85rem; text-transform: uppercase; letter-spacing: 0.05em; color: var(--text-muted); margin-bottom: 6px;">Clinical Description</h4>
          <p style="font-size: 0.95rem; line-height: 1.6; color: var(--text-main); background: #f8fafc; padding: 12px; border-radius: var(--radius-sm); border: 1px solid var(--border);">
            ${escapeHtml(med.description || 'Standard pharmaceutical grade preparation. Store at controlled room temperature.')}
          </p>
        </div>

        ${
          med.requires_rx
            ? `
          <div class="alert alert-warning" style="margin-bottom: 20px; font-size: 0.85rem; line-height: 1.4;">
            <strong>⚠️ Prescription Verification Required</strong><br/>
            This is a regulated prescription medicine. To complete purchase, you must upload an official prescription which will be reviewed by a licensed pharmacist before dispatch.
          </div>
        `
            : `
          <div class="alert alert-info" style="margin-bottom: 20px; font-size: 0.85rem; line-height: 1.4;">
            <strong>ℹ️ Over-the-Counter Medicine</strong><br/>
            This item does not require a prescription. Always consult a healthcare professional for persistent symptoms.
          </div>
        `
        }

        <div style="display: flex; justify-content: space-between; align-items: center; padding-top: 16px; border-top: 1px solid var(--border);">
          <div>
            <div style="font-size: 0.8rem; color: var(--text-muted);">Unit Price</div>
            <div style="font-size: 1.5rem; font-weight: 800; color: var(--primary);">₹${med.price.toFixed(2)}</div>
          </div>
          <div style="display: flex; gap: 8px;">
            <button id="modal-close-btn" class="btn btn-outline">Close</button>
            <button
              id="modal-add-cart-btn"
              class="btn btn-primary"
              ${isOutOfStock ? 'disabled' : ''}
            >
              ${isOutOfStock ? 'Out of Stock' : 'Add to Cart'}
            </button>
          </div>
        </div>
      </div>
    `;

    document.body.appendChild(modalOverlay);

    const closeModal = () => {
      modalOverlay.remove();
    };

    modalOverlay.querySelector('#close-med-modal')?.addEventListener('click', closeModal);
    modalOverlay.querySelector('#modal-close-btn')?.addEventListener('click', closeModal);
    modalOverlay.addEventListener('click', (e) => {
      if (e.target === modalOverlay) closeModal();
    });

    modalOverlay.querySelector('#modal-add-cart-btn')?.addEventListener('click', () => {
      cart.addItem(med, 1);
      const btn = modalOverlay.querySelector('#modal-add-cart-btn') as HTMLButtonElement;
      if (btn) {
        btn.textContent = '✓ Added to Cart!';
        btn.style.background = 'var(--success)';
        setTimeout(() => {
          closeModal();
        }, 800);
      }
    });
  };

  const loadMedicines = async () => {
    grid.innerHTML = `
      <div style="grid-column: 1 / -1; text-align: center; padding: 40px; color: var(--text-muted);">
        Fetching medicines...
      </div>
    `;

    try {
      const rxParam = currentFilter === 'rx' ? true : currentFilter === 'otc' ? false : undefined;
      const res = await api.getMedicines(currentSearch, currentPage, pageSize, rxParam);
      const items = res.data;

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
        pageIndicator.textContent = `Page 1 of 1 (0 items)`;
        return;
      }

      grid.innerHTML = items
        .map((med: Medicine) => {
          const isOutOfStock = med.stock_qty <= 0;
          const isLowStock = med.stock_qty > 0 && med.stock_qty <= 10;
          const stockClass = isOutOfStock ? 'out-of-stock' : isLowStock ? 'low-stock' : 'in-stock';
          const stockLabel = isOutOfStock ? 'Out of Stock' : isLowStock ? `Low Stock (${med.stock_qty})` : `In Stock (${med.stock_qty})`;

          return `
            <div class="card medicine-card" data-med-id="${med.medicine_id}" style="cursor: pointer; transition: transform 0.15s ease, box-shadow 0.15s ease;">
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

              <div class="card-footer" style="margin-top: auto; padding-top: 12px; border-top: 1px solid var(--border);">
                <div class="price">₹${med.price.toFixed(2)}</div>
                <div style="display: flex; gap: 6px;">
                  <button
                    class="btn btn-outline btn-sm view-details-btn"
                    data-med='${JSON.stringify(med).replace(/'/g, "&#39;")}'
                    title="View clinical details"
                  >
                    Info
                  </button>
                  <button
                    class="btn btn-primary btn-sm add-cart-btn"
                    data-med='${JSON.stringify(med).replace(/'/g, "&#39;")}'
                    ${isOutOfStock ? 'disabled' : ''}
                  >
                    ${isOutOfStock ? 'Out' : 'Add'}
                  </button>
                </div>
              </div>
            </div>
          `;
        })
        .join('');

      // Attach add-to-cart handlers
      grid.querySelectorAll('.add-cart-btn').forEach((button) => {
        button.addEventListener('click', (e) => {
          e.stopPropagation();
          const target = e.currentTarget as HTMLButtonElement;
          const medData = JSON.parse(target.getAttribute('data-med') || '{}');
          cart.addItem(medData, 1);
          target.textContent = '✓';
          target.style.background = 'var(--success)';
          setTimeout(() => {
            target.textContent = 'Add';
            target.style.background = '';
          }, 1200);
        });
      });

      // Attach view details modal handlers
      grid.querySelectorAll('.view-details-btn').forEach((button) => {
        button.addEventListener('click', (e) => {
          e.stopPropagation();
          const target = e.currentTarget as HTMLButtonElement;
          const medData = JSON.parse(target.getAttribute('data-med') || '{}');
          showMedicineModal(medData);
        });
      });

      // Card click opens modal
      grid.querySelectorAll('.medicine-card').forEach((cardEl) => {
        cardEl.addEventListener('click', (e) => {
          // If click was on button, do nothing
          if ((e.target as HTMLElement).closest('button')) return;
          const btn = cardEl.querySelector('.view-details-btn') as HTMLButtonElement;
          if (btn) {
            const medData = JSON.parse(btn.getAttribute('data-med') || '{}');
            showMedicineModal(medData);
          }
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
