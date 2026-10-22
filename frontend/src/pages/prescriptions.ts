import { api } from '../api';
import { Prescription } from '../types';

export function renderPrescriptionsPage(container: HTMLElement) {
  const user = api.getCurrentUser();
  if (!user) {
    window.location.hash = '#/login';
    return;
  }

  container.innerHTML = `
    <div style="max-width: 900px; margin: 0 auto;">
      <div style="margin-bottom: 24px;">
        <h1 style="font-size: 1.8rem; font-weight: 800; margin-bottom: 6px;">Upload &amp; Manage Prescriptions</h1>
        <p style="color: var(--text-muted);">
          Upload your medical prescription for automated OCR transcription and F# domain validation assistance.
        </p>
      </div>

      <!-- Safety Alert Box -->
      <div class="alert alert-info" style="display: flex; gap: 12px; align-items: flex-start;">
        <span style="font-size: 1.3rem;">ℹ️</span>
        <div>
          <strong>Human-in-the-Loop Pharmacy Safety</strong>
          <p style="font-size: 0.88rem; margin-top: 4px;">
            AI and OCR tools only assist by extracting text and checking field completeness. Every prescription must be thoroughly examined and approved by a licensed pharmacist before any medication can be dispensed.
          </p>
        </div>
      </div>

      <!-- Upload Form Card -->
      <div class="card" style="padding: 28px; margin-bottom: 30px;">
        <h3 style="margin-bottom: 16px;">Upload Prescription Document</h3>
        <p style="font-size: 0.88rem; color: var(--text-muted); margin-bottom: 16px;">
          Accepted formats: JPEG, PNG, PDF (Maximum 5 MB).
        </p>

        <div id="upload-error" class="alert alert-danger" style="display: none;"></div>

        <div id="dropzone" class="upload-dropzone">
          <div class="upload-icon">📄</div>
          <div style="font-weight: 700; margin-bottom: 6px;">Drag &amp; Drop prescription image or PDF here</div>
          <div style="font-size: 0.85rem; color: var(--text-muted); margin-bottom: 16px;">or click to browse from your device</div>
          <button id="browse-btn" type="button" class="btn btn-outline btn-sm">Select File</button>
          <input type="file" id="file-input" accept="image/jpeg,image/png,application/pdf" style="display: none;" />
        </div>

        <div id="preview-area" style="display: none; margin-top: 20px;">
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px;">
            <strong id="preview-filename" style="font-size: 0.9rem;"></strong>
            <button id="clear-file-btn" class="btn btn-outline btn-sm" style="color: var(--danger);">Change File</button>
          </div>
          <div id="preview-container" style="text-align: center; max-height: 250px; overflow: hidden; border-radius: var(--radius-sm); border: 1px solid var(--border); padding: 10px; background: #fafafa;"></div>
          <button id="start-upload-btn" class="btn btn-primary" style="width: 100%; margin-top: 14px;">
            Submit Prescription for Verification
          </button>
        </div>

        <!-- Progress Indicator -->
        <div id="progress-container" style="display: none; margin-top: 20px;">
          <div style="display: flex; justify-content: space-between; font-size: 0.85rem; font-weight: 600; margin-bottom: 6px;">
            <span id="progress-status-text">Uploading &amp; processing...</span>
            <span id="progress-percent">0%</span>
          </div>
          <div style="height: 8px; background: #e2e8f0; border-radius: 9999px; overflow: hidden;">
            <div id="progress-bar" style="height: 100%; width: 0%; background: var(--primary); transition: width 0.3s;"></div>
          </div>
        </div>

        <!-- Live Status / Result Panel -->
        <div id="live-result" style="display: none; margin-top: 24px; border-top: 1px solid var(--border); padding-top: 20px;">
          <h4 style="margin-bottom: 12px;">Prescription Processing Status</h4>
          <div id="result-badges" style="display: flex; gap: 10px; margin-bottom: 16px;"></div>

          <div style="margin-bottom: 10px; font-size: 0.88rem; font-weight: 600;">Extracted Clinical Transcription (OCR Suggestion):</div>
          <div id="ocr-text-view" class="ocr-box"></div>
          <div style="margin-top: 10px; font-size: 0.8rem; color: var(--text-muted); font-style: italic;">
            * This transcription is forwarded to the staff pharmacist review queue.
          </div>
        </div>
      </div>

      <!-- Prescriptions History -->
      <div>
        <h3 style="margin-bottom: 16px;">My Uploaded Prescriptions</h3>
        <div id="prescriptions-table-container">
          <div style="text-align: center; padding: 30px; color: var(--text-muted);">Loading prescriptions history...</div>
        </div>
      </div>
    </div>
  `;

  const dropzone = container.querySelector('#dropzone') as HTMLDivElement;
  const fileInput = container.querySelector('#file-input') as HTMLInputElement;
  const browseBtn = container.querySelector('#browse-btn') as HTMLButtonElement;
  const previewArea = container.querySelector('#preview-area') as HTMLDivElement;
  const previewFilename = container.querySelector('#preview-filename') as HTMLElement;
  const previewContainer = container.querySelector('#preview-container') as HTMLDivElement;
  const clearFileBtn = container.querySelector('#clear-file-btn') as HTMLButtonElement;
  const startUploadBtn = container.querySelector('#start-upload-btn') as HTMLButtonElement;
  const errorDiv = container.querySelector('#upload-error') as HTMLDivElement;
  const progressContainer = container.querySelector('#progress-container') as HTMLDivElement;
  const progressBar = container.querySelector('#progress-bar') as HTMLDivElement;
  const progressText = container.querySelector('#progress-status-text') as HTMLSpanElement;
  const progressPercent = container.querySelector('#progress-percent') as HTMLSpanElement;
  const liveResult = container.querySelector('#live-result') as HTMLDivElement;
  const resultBadges = container.querySelector('#result-badges') as HTMLDivElement;
  const ocrTextView = container.querySelector('#ocr-text-view') as HTMLDivElement;
  const historyContainer = container.querySelector('#prescriptions-table-container') as HTMLDivElement;

  let selectedFile: File | null = null;
  let pollInterval: any = null;

  const showError = (msg: string) => {
    errorDiv.textContent = msg;
    errorDiv.style.display = 'block';
  };

  const clearError = () => {
    errorDiv.style.display = 'none';
  };

  const handleSelectedFile = (file: File) => {
    clearError();
    if (file.size > 5 * 1024 * 1024) {
      showError('File size exceeds 5MB limit. Please choose a smaller image or compressed PDF.');
      return;
    }

    selectedFile = file;
    previewFilename.textContent = `${file.name} (${(file.size / 1024).toFixed(1)} KB)`;
    dropzone.style.display = 'none';
    previewArea.style.display = 'block';

    previewContainer.innerHTML = '';
    if (file.type.startsWith('image/')) {
      const img = document.createElement('img');
      img.src = URL.createObjectURL(file);
      img.style.maxHeight = '200px';
      img.style.maxWidth = '100%';
      img.style.objectFit = 'contain';
      previewContainer.appendChild(img);
    } else {
      previewContainer.innerHTML = `<div style="padding: 20px; font-size: 1.1rem;">📄 PDF Document: ${escapeHtml(file.name)}</div>`;
    }
  };

  browseBtn.addEventListener('click', () => fileInput.click());
  dropzone.addEventListener('click', () => fileInput.click());

  fileInput.addEventListener('change', () => {
    if (fileInput.files && fileInput.files[0]) {
      handleSelectedFile(fileInput.files[0]);
    }
  });

  dropzone.addEventListener('dragover', (e) => {
    e.preventDefault();
    dropzone.classList.add('dragover');
  });

  dropzone.addEventListener('dragleave', () => {
    dropzone.classList.remove('dragover');
  });

  dropzone.addEventListener('drop', (e) => {
    e.preventDefault();
    dropzone.classList.remove('dragover');
    if (e.dataTransfer && e.dataTransfer.files && e.dataTransfer.files[0]) {
      handleSelectedFile(e.dataTransfer.files[0]);
    }
  });

  clearFileBtn.addEventListener('click', () => {
    selectedFile = null;
    fileInput.value = '';
    previewArea.style.display = 'none';
    dropzone.style.display = 'block';
    clearError();
  });

  startUploadBtn.addEventListener('click', async () => {
    if (!selectedFile) return;

    try {
      startUploadBtn.disabled = true;
      progressContainer.style.display = 'block';
      progressBar.style.width = '30%';
      progressPercent.textContent = '30%';
      progressText.textContent = 'Uploading prescription file...';

      const created = await api.uploadPrescription(selectedFile);

      progressBar.style.width = '60%';
      progressPercent.textContent = '60%';
      progressText.textContent = 'Prescription uploaded. Asynchronous OCR & F# validation in progress...';

      // Start polling status
      pollPrescription(created.prescription_id);
    } catch (err: any) {
      showError(err.message || 'Upload failed');
      startUploadBtn.disabled = false;
      progressContainer.style.display = 'none';
    }
  });

  const pollPrescription = (prescriptionId: string) => {
    let attempts = 0;
    const maxAttempts = 15;

    pollInterval = setInterval(async () => {
      attempts++;
      try {
        const p = await api.getPrescription(prescriptionId);
        if (p.validation_status !== 'PENDING' || attempts >= maxAttempts) {
          clearInterval(pollInterval);
          progressBar.style.width = '100%';
          progressPercent.textContent = '100%';
          progressText.textContent = 'Analysis complete!';
          showPrescriptionResult(p);
          loadHistory();
        }
      } catch (err) {
        console.error('Error polling prescription status:', err);
      }
    }, 2000);
  };

  const showPrescriptionResult = (p: Prescription) => {
    liveResult.style.display = 'block';

    let valClass = 'badge-otc';
    if (p.validation_status === 'VALID') valClass = 'badge-otc';
    else if (p.validation_status === 'INCOMPLETE' || p.validation_status === 'MANUAL_REVIEW') valClass = 'badge-warning';
    else if (p.validation_status === 'INVALID') valClass = 'badge-rx';

    resultBadges.innerHTML = `
      <div><strong>Validation Status:</strong> <span class="badge ${valClass}">${p.validation_status}</span></div>
      <div><strong>Review Status:</strong> <span class="badge ${p.review_status === 'APPROVED' ? 'badge-otc' : 'badge-warning'}">${p.review_status}</span></div>
    `;

    ocrTextView.textContent = p.extracted_text || '(Processing or no text detected)';
  };

  const loadHistory = async () => {
    try {
      const list = await api.getMyPrescriptions();
      if (list.length === 0) {
        historyContainer.innerHTML = `
          <div style="text-align: center; padding: 40px; background: white; border-radius: var(--radius-md); border: 1px dashed var(--border); color: var(--text-muted);">
            No prescriptions uploaded yet.
          </div>
        `;
        return;
      }

      historyContainer.innerHTML = `
        <div class="table-container">
          <table>
            <thead>
              <tr>
                <th>Prescription ID</th>
                <th>Date Uploaded</th>
                <th>Validation</th>
                <th>Review Status</th>
                <th>Pharmacist Notes</th>
                <th>File</th>
              </tr>
            </thead>
            <tbody>
              ${list.map((p) => `
                <tr>
                  <td><code>${p.prescription_id.substring(0, 8)}...</code></td>
                  <td>${new Date(p.created_at).toLocaleString()}</td>
                  <td><span class="badge ${p.validation_status === 'VALID' ? 'badge-otc' : 'badge-warning'}">${p.validation_status}</span></td>
                  <td><span class="badge ${p.review_status === 'APPROVED' ? 'badge-otc' : p.review_status === 'REJECTED' ? 'badge-rx' : 'badge-warning'}">${p.review_status}</span></td>
                  <td>${escapeHtml(p.review_note || 'Awaiting pharmacist review')}</td>
                  <td>
                    <a href="${api.getPrescriptionFileUrl(p.prescription_id)}" target="_blank" class="btn btn-outline btn-sm">
                      View File
                    </a>
                  </td>
                </tr>
              `).join('')}
            </tbody>
          </table>
        </div>
      `;
    } catch (err: any) {
      historyContainer.innerHTML = `<div class="alert alert-danger">Failed to load prescription history: ${err.message}</div>`;
    }
  };

  loadHistory();
}

function escapeHtml(str: string): string {
  const div = document.createElement('div');
  div.textContent = str;
  return div.innerHTML;
}
