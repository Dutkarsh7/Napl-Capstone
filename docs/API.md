# Smart Pharmacy REST API Documentation

This document specifies all endpoints, request/response formats, authorization requirements, and error codes for the **AI-Powered Smart Pharmacy and Prescription Assistant**.

> [!IMPORTANT]
> **Safety Notice**: This platform is designed strictly for pharmacy management and prescription verification assistance. It **NEVER** diagnoses conditions, prescribes medication, or alters doctor-prescribed dosages. Every prescription requires human review and explicit approval by a licensed pharmacist (staff user) before it can fulfill an order.

---

## Authentication & Headers

Protected endpoints require a JSON Web Token (JWT) supplied in the HTTP `Authorization` header:

```http
Authorization: Bearer <jwt-token>
```

Tokens have a 24-hour expiration and encode the user's `user_id`, `email`, and `role` (`customer`, `staff`, or `admin`).

### Standard Error Envelope

When an error occurs, the server responds with a consistent JSON envelope:

```json
{
  "error": "Descriptive error message"
}
```

Common status codes:
- `400 Bad Request` - Malformed payload, invalid transitions, or missing fields.
- `401 Unauthorized` - Missing or expired JWT token.
- `403 Forbidden` - Insufficient role permissions.
- `404 Not Found` - Requested entity does not exist.
- `409 Conflict` - Email collision or insufficient medicine stock.
- `500 Internal Server Error` - Unrecoverable server failure.

---

## 1. Authentication Endpoints

### Register Customer
- **Endpoint**: `POST /api/auth/register`
- **Auth**: Public
- **Request Body**:
  ```json
  {
    "name": "Jane Doe",
    "email": "jane@example.com",
    "password": "SecurePassword123!"
  }
  ```
- **Response** (`201 Created`):
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6...",
    "user": {
      "user_id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
      "name": "Jane Doe",
      "email": "jane@example.com",
      "role": "customer",
      "created_at": "2026-10-22T08:00:00Z"
    }
  }
  ```

### Login
- **Endpoint**: `POST /api/auth/login`
- **Auth**: Public
- **Request Body**:
  ```json
  {
    "email": "customer@pharmacy.local",
    "password": "CustomerPass123!"
  }
  ```
- **Response** (`200 OK`): Same as Register.

### Get Current Profile
- **Endpoint**: `GET /api/me`
- **Auth**: Bearer token
- **Response** (`200 OK`):
  ```json
  {
    "user_id": "33333333-3333-3333-3333-333333333333",
    "name": "John Doe (Customer)",
    "email": "customer@pharmacy.local",
    "role": "customer",
    "created_at": "2026-10-22T00:00:00Z"
  }
  ```

---

## 2. Medicine Catalogue & Inventory

### List & Search Medicines
- **Endpoint**: `GET /api/medicines?q={query}&page={page}&limit={limit}`
- **Auth**: Public
- **Query Parameters**:
  - `q` (optional): Case-insensitive search on name or generic name (e.g. `amox`).
  - `page` (optional, default `1`): Page number.
  - `limit` (optional, default `20`): Page size (max 100).
- **Response** (`200 OK`):
  ```json
  {
    "data": [
      {
        "medicine_id": "b1a2c3d4-...",
        "name": "Amoxicillin 500mg Capsule",
        "generic_name": "Amoxicillin",
        "description": "Broad-spectrum beta-lactam antibiotic capsule 500mg",
        "price": 14.00,
        "stock_qty": 150,
        "requires_rx": true,
        "created_at": "2026-10-22T00:00:00Z"
      }
    ],
    "total": 105,
    "page": 1,
    "limit": 20
  }
  ```

### Get Single Medicine
- **Endpoint**: `GET /api/medicines/{id}`
- **Auth**: Public
- **Response** (`200 OK`): Single `Medicine` object.

### Create Medicine (Staff/Admin)
- **Endpoint**: `POST /api/medicines`
- **Auth**: Staff or Admin
- **Request Body**:
  ```json
  {
    "name": "New Medicine 250mg",
    "generic_name": "Active Ingredient",
    "description": "Neutral non-prescriptive description",
    "price": 12.50,
    "stock_qty": 100,
    "requires_rx": true
  }
  ```
- **Response** (`201 Created`): Created `Medicine` record.

### Update Stock & Price (Staff/Admin)
- **Endpoint**: `PATCH /api/medicines/{id}/stock`
- **Auth**: Staff or Admin
- **Request Body**:
  ```json
  {
    "stock_qty": 180,
    "price": 13.00
  }
  ```
- **Response** (`200 OK`): Updated `Medicine` record.

---

## 3. Prescriptions

### Upload Prescription
- **Endpoint**: `POST /api/prescriptions`
- **Auth**: Bearer token (Customer/Staff/Admin)
- **Content-Type**: `multipart/form-data`
- **Form Field**: `file` (JPEG, PNG, or PDF up to 5 MB)
- **Behavior**: Sniffs content bytes, stores file securely outside web root with a UUID filename, records prescription with `validation_status: "PENDING"`, and dispatches background worker job.
- **Response** (`201 Created`):
  ```json
  {
    "prescription_id": "9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d",
    "user_id": "33333333-3333-3333-3333-333333333333",
    "file_path": "uploads/...",
    "extracted_text": "",
    "validation_status": "PENDING",
    "review_status": "PENDING",
    "review_note": "",
    "created_at": "2026-10-22T08:15:00Z"
  }
  ```

### Get My Prescriptions
- **Endpoint**: `GET /api/prescriptions/mine`
- **Auth**: Bearer token
- **Response** (`200 OK`): Array of user's prescriptions.

### Prescription Review Queue (Staff)
- **Endpoint**: `GET /api/prescriptions/queue`
- **Auth**: Staff or Admin
- **Response** (`200 OK`): Array of prescriptions where `review_status == "PENDING"`.

### Get Prescription Details
- **Endpoint**: `GET /api/prescriptions/{id}`
- **Auth**: Prescription owner, Staff, or Admin
- **Response** (`200 OK`): Single `Prescription` record with live `validation_status` and `review_status`.

### View/Download Prescription File
- **Endpoint**: `GET /api/prescriptions/{id}/file`
- **Auth**: Prescription owner, Staff, or Admin
- **Response** (`200 OK`): File binary stream with detected MIME type.

### Review Prescription (Staff/Admin)
- **Endpoint**: `POST /api/prescriptions/{id}/review`
- **Auth**: Staff or Admin
- **Request Body**:
  ```json
  {
    "action": "APPROVE",
    "note": "Verified doctor license #MED-88492. Dosages and duration within clinical parameters."
  }
  ```
- **Response** (`200 OK`): Updated `Prescription` record.

---

## 4. Orders & Stock Transactions

### Create Order
- **Endpoint**: `POST /api/orders`
- **Auth**: Bearer token
- **Request Body**:
  ```json
  {
    "prescription_id": "9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d",
    "items": [
      {
        "medicine_id": "b1a2c3d4-...",
        "quantity": 2
      }
    ]
  }
  ```
- **Rules Enforced**:
  - If any item requires Rx, `prescription_id` must reference an `APPROVED` prescription owned by the calling user.
  - Stock is locked via `SELECT ... FOR UPDATE` inside a PostgreSQL transaction to prevent race conditions and overselling.
  - Stock is decremented and initial status is set to `PLACED`.
- **Response** (`201 Created`): Created `Order` object with items and initial history entry.

### List Customer's Orders
- **Endpoint**: `GET /api/orders/mine`
- **Auth**: Bearer token
- **Response** (`200 OK`): Array of orders with items.

### Get Order by ID
- **Endpoint**: `GET /api/orders/{id}`
- **Auth**: Order owner, Staff, or Admin
- **Response** (`200 OK`): Order details including timeline history.

### List All Orders (Staff/Admin)
- **Endpoint**: `GET /api/orders`
- **Auth**: Staff or Admin
- **Response** (`200 OK`): Complete list of orders across all customers.

### Update Order Status (Staff/Admin)
- **Endpoint**: `PATCH /api/orders/{id}/status`
- **Auth**: Staff or Admin
- **Request Body**:
  ```json
  {
    "status": "CONFIRMED"
  }
  ```
- **Valid Transitions**:
  - `PLACED` &rarr; `UNDER_REVIEW`, `CANCELLED`
  - `UNDER_REVIEW` &rarr; `CONFIRMED`, `CANCELLED`
  - `CONFIRMED` &rarr; `READY`, `CANCELLED`
  - `READY` &rarr; `COMPLETED`, `CANCELLED`
  - Transitioning to `CANCELLED` automatically restores stock to each medicine in the order.
- **Response** (`200 OK`): Updated `Order` object.

---

## 5. Administration

### List Users
- **Endpoint**: `GET /api/users`
- **Auth**: Admin only
- **Response** (`200 OK`): Array of user profiles.

### Update User Role
- **Endpoint**: `PATCH /api/users/{id}/role`
- **Auth**: Admin only
- **Request Body**:
  ```json
  {
    "role": "staff"
  }
  ```
- **Response** (`200 OK`): Updated `User` profile.

---

## 6. Internal Validation Microservice (F# Giraffe)

- **Endpoint**: `POST /validate/prescription`
- **Host**: Port `5005` (Internal)
- **Request Body**:
  ```json
  {
    "patientName": "Alice Smith",
    "doctorName": "Dr. John Watson",
    "doctorRegNo": "MED-99482",
    "date": "2026-10-15",
    "medicines": [
      {
        "name": "Amoxicillin 500mg",
        "dosage": "1 cap TID",
        "quantity": 21
      }
    ]
  }
  ```
- **Response** (`200 OK`):
  ```json
  {
    "status": "valid",
    "missing": [],
    "issues": []
  }
  ```

---

## 7. Internal OCR Microservice (Python FastAPI)

- **Endpoint**: `POST /ocr`
- **Host**: Port `8000` (Internal)
- **Content-Type**: `multipart/form-data`
- **Form Field**: `file`
- **Response** (`200 OK`):
  ```json
  {
    "text": "St. Jude Medical Clinic\nDoctor: Dr. Elizabeth Adams, MD\n..."
  }
  ```
