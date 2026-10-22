import { api } from '../api';

export function renderAuthPage(container: HTMLElement, isRegister: boolean = false) {
  container.innerHTML = `
    <div style="max-width: 480px; margin: 40px auto;">
      <div class="card" style="padding: 30px;">
        <h2 style="text-align: center; margin-bottom: 8px;">
          ${isRegister ? 'Create Your Account' : 'Sign In to Smart Pharmacy'}
        </h2>
        <p style="text-align: center; color: var(--text-muted); margin-bottom: 24px; font-size: 0.9rem;">
          ${isRegister ? 'Register as a customer to browse medicines and upload prescriptions' : 'Access your prescription history and pharmacy orders'}
        </p>

        <div id="auth-error" class="alert alert-danger" style="display: none;"></div>

        <form id="auth-form">
          ${isRegister ? `
            <div class="form-group">
              <label for="reg-name">Full Name</label>
              <input type="text" id="reg-name" class="form-input" placeholder="e.g. John Doe" required />
            </div>
          ` : ''}

          <div class="form-group">
            <label for="auth-email">Email Address</label>
            <input type="email" id="auth-email" class="form-input" placeholder="you@domain.com" required />
          </div>

          <div class="form-group">
            <label for="auth-password">Password</label>
            <input type="password" id="auth-password" class="form-input" placeholder="••••••••" required />
          </div>

          <button type="submit" id="auth-submit-btn" class="btn btn-primary" style="width: 100%; margin-top: 10px;">
            ${isRegister ? 'Register Account' : 'Sign In'}
          </button>
        </form>

        <div style="text-align: center; margin-top: 18px; font-size: 0.88rem;">
          ${isRegister ? `
            Already have an account? <a href="#/login" style="color: var(--primary); font-weight: 600;">Sign in here</a>
          ` : `
            Don't have an account? <a href="#/register" style="color: var(--primary); font-weight: 600;">Create one here</a>
          `}
        </div>

        <hr style="margin: 24px 0; border: none; border-top: 1px solid var(--border);" />

        <div>
          <p style="font-size: 0.82rem; font-weight: 700; text-transform: uppercase; color: var(--text-muted); margin-bottom: 10px; text-align: center;">
            ⚡ Quick Demo Logins
          </p>
          <div style="display: flex; flex-direction: column; gap: 8px;">
            <button id="demo-customer-btn" class="btn btn-outline btn-sm" type="button">
              👤 Login as Demo Customer (John Doe)
            </button>
            <button id="demo-staff-btn" class="btn btn-outline btn-sm" type="button" style="border-color: #93c5fd; color: #1e40af;">
              🩺 Login as Demo Pharmacist (Sarah Jenkins)
            </button>
            <button id="demo-admin-btn" class="btn btn-outline btn-sm" type="button" style="border-color: #f472b6; color: #9d174d;">
              🔑 Login as Demo Admin Supervisor
            </button>
          </div>
        </div>
      </div>
    </div>
  `;

  const form = container.querySelector('#auth-form') as HTMLFormElement;
  const errorDiv = container.querySelector('#auth-error') as HTMLDivElement;
  const submitBtn = container.querySelector('#auth-submit-btn') as HTMLButtonElement;

  const showError = (msg: string) => {
    errorDiv.textContent = msg;
    errorDiv.style.display = 'block';
  };

  const handleLogin = async (email: string, pass: string) => {
    try {
      submitBtn.disabled = true;
      submitBtn.textContent = 'Authenticating...';
      errorDiv.style.display = 'none';

      const res = await api.login(email, pass);
      window.dispatchEvent(new CustomEvent('auth-changed'));

      if (res.user.role === 'staff' || res.user.role === 'admin') {
        window.location.hash = '#/staff';
      } else {
        window.location.hash = '#/medicines';
      }
    } catch (err: any) {
      showError(err.message || 'Login failed');
    } finally {
      submitBtn.disabled = false;
      submitBtn.textContent = isRegister ? 'Register Account' : 'Sign In';
    }
  };

  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    const email = (container.querySelector('#auth-email') as HTMLInputElement).value;
    const password = (container.querySelector('#auth-password') as HTMLInputElement).value;

    if (isRegister) {
      const name = (container.querySelector('#reg-name') as HTMLInputElement).value;
      try {
        submitBtn.disabled = true;
        submitBtn.textContent = 'Registering...';
        errorDiv.style.display = 'none';

        await api.register(name, email, password);
        window.dispatchEvent(new CustomEvent('auth-changed'));
        window.location.hash = '#/medicines';
      } catch (err: any) {
        showError(err.message || 'Registration failed');
      } finally {
        submitBtn.disabled = false;
        submitBtn.textContent = 'Register Account';
      }
    } else {
      await handleLogin(email, password);
    }
  });

  // Demo auto-fill buttons
  container.querySelector('#demo-customer-btn')?.addEventListener('click', () => {
    handleLogin('customer@pharmacy.local', 'CustomerPass123!');
  });
  container.querySelector('#demo-staff-btn')?.addEventListener('click', () => {
    handleLogin('staff@pharmacy.local', 'StaffPass123!');
  });
  container.querySelector('#demo-admin-btn')?.addEventListener('click', () => {
    handleLogin('admin@pharmacy.local', 'AdminPass123!');
  });
}
