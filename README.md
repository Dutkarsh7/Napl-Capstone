# AI-Powered Smart Pharmacy and Prescription Assistant

A full-stack, distributed pharmacy management and prescription verification assistance platform developed as a University Capstone Project.

---

## ⚕️ Critical Safety Framing

> [!IMPORTANT]
> **Safety Regulations & Human-in-the-Loop Mandate**:
> This platform is strictly designed to assist licensed pharmacy personnel with digital workflows.
> - The platform **NEVER** diagnoses illnesses, prescribes medications, recommends treatments, or adjusts dosages.
> - Optical Character Recognition (OCR) and automated syntax validation are strictly clinical suggestions for human verification.
> - **Mandatory Review**: Every prescription must be examined, verified, and explicitly approved by a licensed pharmacist (staff user) before it can be used to fulfill any prescription-controlled order.
> - **Downtime Resilience**: If automated OCR or validation microservices are unreachable or time out, the system defaults safely to `MANUAL_REVIEW`. It **never** auto-approves.

---

## Tech Stack

| Layer | Technology | Key Capabilities |
|---|---|---|
| **Backend** | Go 1.22+ (`net/http`) | Standard library method routing (`GET /api/medicines`), `pgx/v5` PostgreSQL connection pool, JWT authentication (`jwt/v5`), `bcrypt` password hashing, buffered-channel asynchronous worker pool, `SELECT ... FOR UPDATE` row locks, graceful shutdown. |
| **Validation Microservice** | F# on .NET 8 / 10 + Giraffe | Domain modeling with discriminated unions (`Valid`, `Incomplete`, `Invalid`), pure functions, immutable rule evaluations, controlled substance flags, Expecto unit test suite. |
| **OCR Microservice** | Python 3 + FastAPI | `pytesseract` and `Pillow` image processing with automatic fallback transcription and health checks. |
| **Database** | PostgreSQL | Foreign key relationships, `CHECK` constraints, UUID primary keys, indices, automated migration scripts, and 105 realistic seed medicines. |
| **Frontend** | Vanilla TypeScript + Vite | HTML5 + CSS3 (CSS custom properties), responsive card layouts, hash routing, drag-and-drop file upload with live status polling, fulfillment timeline stepper, cart state management. |
| **Testing** | Go `testing`, Expecto, Vitest | 37 Go tests, 23 F# tests, and 8 Vitest tests (68+ total automated tests covering all state machines, stock concurrency, upload sniffing, and rule cases). |

---

## Repository Layout

```
smart-pharmacy/
├── backend-go/                # Go REST API & worker pool
│   ├── cmd/server/main.go     # Server entry point, ServeMux 1.22 routing
│   ├── internal/
│   │   ├── config/            # Environment configuration loader
│   │   ├── db/                # PostgreSQL pgxpool & in-memory stores
│   │   ├── fsclient/          # HTTP client for F# validation service
│   │   ├── handlers/          # Auth, Medicines, Prescriptions, Orders, Admin
│   │   ├── middleware/        # JWT auth, role enforcement, CORS, logger
│   │   ├── models/            # Domain models, transitions, and OCR structures
│   │   ├── ocrclient/         # HTTP client for Python OCR service
│   │   └── services/          # Regex text parser, worker pool, transition helpers
│   ├── tests/                 # Unit & httptest test suite
│   ├── Dockerfile
│   └── go.mod
├── validation-fsharp/         # F# Giraffe domain validation microservice
│   ├── Domain.fs              # Discriminated unions & records
│   ├── Rules.fs               # Pure immutable validation rules
│   ├── Api.fs                 # Giraffe HTTP handlers
│   ├── Program.fs             # ASP.NET Core host
│   ├── Tests/                 # Expecto unit test project (23 tests)
│   ├── Dockerfile
│   └── PharmacyValidation.fsproj
├── ocr-python/                # Python FastAPI OCR microservice
│   ├── main.py                # FastAPI app with pytesseract & fallback
│   ├── requirements.txt
│   └── Dockerfile
├── frontend/                  # TypeScript + Vite SPA
│   ├── index.html             # Entry point with safety disclaimer banner
│   ├── src/
│   │   ├── api.ts             # API client with token storage
│   │   ├── cart.ts            # Cart store with Rx tracking
│   │   ├── types.ts           # Complete TypeScript interfaces
│   │   ├── main.ts            # Hash router & navbar controller
│   │   ├── pages/             # Auth, Medicines, Cart, Prescriptions, Orders, Staff, Admin
│   │   └── tests/             # Vitest test suite
│   ├── styles.css             # Responsive CSS with custom variables
│   ├── vite.config.ts
│   └── Dockerfile
├── db/migrations/             # PostgreSQL migrations
│   ├── 001_init.sql           # Schema DDL (users, medicines, prescriptions, orders, history)
│   └── 002_seed.sql           # 105 realistic medicines + 3 demo user accounts
├── docs/                      # Technical documentation
│   ├── API.md                 # Full REST API specification
│   └── TEST_CASES.md          # 53-item test case matrix
├── docker-compose.yml         # Multi-service container orchestration
├── .env.example               # Environment variables template
└── README.md                  # This comprehensive guide
```

---

## Demo User Accounts

The database seed (`002_seed.sql`) pre-populates three realistic demo users:

| Role | Email | Password | Intended Capabilities |
|---|---|---|---|
| **Pharmacist (Staff)** | `staff@pharmacy.local` | `StaffPass123!` | Review prescription queue, approve/reject with notes, update order fulfillment stages, modify inventory stock & prices. |
| **Customer** | `customer@pharmacy.local` | `CustomerPass123!` | Search medicines, upload prescriptions, add items to cart, attach approved prescriptions, place orders, view timeline. |
| **System Admin** | `admin@pharmacy.local` | `AdminPass123!` | Full privileges plus user administration and role promotion/demotion. |

*(Note: The login page also features one-click quick login buttons to immediately sign in as any of these roles without typing).*

---

## Quick Start with Docker Compose

To build and run all services (PostgreSQL, Go backend, F# validation, Python OCR, Frontend Nginx):

```bash
docker compose up --build
```

Access the application:
- **Frontend SPA**: `http://localhost:5173`
- **Go Backend API**: `http://localhost:8080` (Health: `http://localhost:8080/health`)
- **F# Validation Microservice**: `http://localhost:5005` (Health: `http://localhost:5005/health`)
- **Python OCR Microservice**: `http://localhost:8000` (Health: `http://localhost:8000/health`)

---

## Local Run Instructions for Windows (Non-Docker)

You can run each service independently in standard PowerShell terminals.

### 1. F# Validation Service
```powershell
cd smart-pharmacy\validation-fsharp
dotnet run
```
*Listens on `http://localhost:5005`*.

### 2. Python OCR Microservice
```powershell
cd smart-pharmacy\ocr-python
pip install -r requirements.txt
python main.py
```
*Listens on `http://localhost:8000`*.

### 3. Go Backend
If running without a local PostgreSQL instance, the Go backend will automatically detect this and initialize a high-performance, thread-safe in-memory store pre-seeded with all 105 realistic medicines and demo users:
```powershell
cd smart-pharmacy\backend-go
go run .\cmd\server
```
*Listens on `http://localhost:8080`*.

### 4. Frontend (TypeScript + Vite)
```powershell
cd smart-pharmacy\frontend
npm run dev
```
*Opens at `http://localhost:5173`*.

---

## Running the Automated Test Suites

### 1. Go Backend Tests (37 Unit & Handler Tests)
Tests state machine transitions, stock decrement/cancellation concurrency logic, file upload content sniffing, JWT authentication middleware, and REST handlers:
```powershell
cd smart-pharmacy\backend-go
go test -v ./...
```

### 2. F# Expecto Tests (23 Domain Validation Tests)
Tests pure immutable domain rules covering `Valid`, `Incomplete`, and `Invalid` cases, boundary quantities (1 and 30), date ranges, and controlled opioid keyword detection:
```powershell
cd smart-pharmacy\validation-fsharp\Tests
dotnet test
```

### 3. Frontend Vitest Tests (8 Client & Cart Tests)
Tests cart calculations, Rx prescription requirement tracking, and API auth token persistence:
```powershell
cd smart-pharmacy\frontend
npm test
```

---

## End-to-End Walkthrough Script

Follow these steps to demonstrate the full application workflow:

1. **Browse Medicines**:
   - Navigate to `http://localhost:5173/#/medicines`.
   - Use the debounced search bar to filter by active ingredient (e.g. type `amox` or `paracetamol`).
   - Notice the distinct badges: green `OTC` for over-the-counter items and red `Rx Required` for prescription drugs.

2. **Upload a Prescription**:
   - Click **Sign In** and use the quick button **"Login as Demo Customer"** (`customer@pharmacy.local`).
   - Go to **Prescriptions** (`#/prescriptions`).
   - Drag and drop or browse for an image or PDF prescription document (a sample image or PDF file up to 5 MB).
   - Click **Submit Prescription for Verification**.
   - Watch the live progress indicator as the Go backend dispatches the job to the asynchronous worker pool.
   - The worker calls the OCR microservice and passes the parsed result to the F# validation service.
   - The result displays the extracted text and validation status badge (e.g. `VALID`). Notice the review status remains `PENDING` awaiting human approval.

3. **Pharmacist Clinical Review**:
   - In the top-right user menu, click **Sign Out**.
   - Click **Sign In** and select **"Login as Demo Pharmacist"** (`staff@pharmacy.local`).
   - Click **Staff Dashboard** (`#/staff`).
   - In the **Prescription Queue** tab, see the customer's uploaded prescription.
   - Click **View Document** to inspect the original file.
   - Enter a note (e.g. *"Verified doctor license and recommended dosage. Approved."*) and click **✓ Approve Prescription**.

4. **Order Placement with Approved Prescription**:
   - Sign back in as the customer (`customer@pharmacy.local`).
   - In **Medicines**, add `Amoxicillin 500mg Capsule` (an Rx item) to your cart.
   - Go to **Cart** (`#/cart`). Notice the warning: *Prescription-Only Medication in Cart*.
   - Select your approved prescription from the mandatory dropdown.
   - Click **Place Order**.
   - The Go backend verifies the prescription is approved and owned by the user, decrements the medicine's stock quantity inside a database transaction (`SELECT ... FOR UPDATE`), and initializes the order status as `PLACED`.

5. **Track Fulfillment Timeline**:
   - Go to **My Orders** (`#/orders`).
   - View the active order showing the fulfillment stepper: `PLACED` &rarr; `UNDER_REVIEW` &rarr; `CONFIRMED` &rarr; `READY` &rarr; `COMPLETED`.

6. **Advance Order Fulfillment**:
   - Sign in as the pharmacist (`staff@pharmacy.local`).
   - Go to **Staff Dashboard** &rarr; **Order Fulfillment** tab.
   - Click **Mark Under Review**, then **Confirm Order**, then **Mark Ready**, and finally **Complete Order**.
   - If an order is **Cancelled**, the Go backend executes a transaction that automatically restores the reserved quantities back into inventory stock.

---

## License & University Capstone Information

Developed for academic presentation and demonstration of resilient microservice architectures, functional domain validation, and human-in-the-loop safety in healthcare operations.