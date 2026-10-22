# Capstone Test Matrix & Verification Cases

This document details the test matrix covering all functional areas of the **AI-Powered Smart Pharmacy and Prescription Assistant**, spanning backend unit and integration tests, F# domain validation tests, and frontend integration tests.

| ID | Area | Test Description / Input | Expected Result | Automated Suite | Status |
|---|---|---|---|---|---|
| **TC-01** | Auth | Register new user with valid name, email, password | Returns 201 Created with JWT token and customer role | Go `handlers_test.go` | PASS |
| **TC-02** | Auth | Register with duplicate email | Returns 409 Conflict with clear error message | Go `handlers_test.go` | PASS |
| **TC-03** | Auth | Register with short password (< 6 chars) | Returns 400 Bad Request | Go `handlers_test.go` | PASS |
| **TC-04** | Auth | Login with valid credentials (`customer@pharmacy.local`) | Returns 200 OK and valid JWT token | Go `handlers_test.go` | PASS |
| **TC-05** | Auth | Login with incorrect password | Returns 401 Unauthorized | Go `handlers_test.go` | PASS |
| **TC-06** | Auth | JWT signature tampering test | Token rejection with signature mismatch error | Go `auth_test.go` | PASS |
| **TC-07** | Auth | Access protected endpoint without Authorization header | Returns 401 Unauthorized with JSON error envelope | Go `auth_test.go` | PASS |
| **TC-08** | Auth | Customer accessing staff-only endpoint (`/queue`) | Returns 403 Forbidden | Go `auth_test.go` | PASS |
| **TC-09** | Auth | Staff accessing staff-only endpoint (`/queue`) | Returns 200 OK | Go `auth_test.go` | PASS |
| **TC-10** | Catalogue | Search query `q=amox` | Returns matching medicines case-insensitively | Go `handlers_test.go` | PASS |
| **TC-11** | Catalogue | Search query with empty string | Returns paginated list of all medicines | Go `handlers_test.go` | PASS |
| **TC-12** | Catalogue | Query with non-existent medicine name | Returns empty data array with total count 0 | Go `handlers_test.go` | PASS |
| **TC-13** | Catalogue | Pagination `page=2&limit=5` | Returns items 6-10 with correct page offset | Go `handlers_test.go` | PASS |
| **TC-14** | Inventory | Staff updates stock and price for medicine | Returns 200 OK and updated medicine attributes | Go `handlers_test.go` | PASS |
| **TC-15** | Upload | Sniff content of valid PNG image | Returns 201 Created and saves with UUID filename | Go `upload_test.go` | PASS |
| **TC-16** | Upload | Sniff content of valid PDF document | Returns 201 Created with .pdf extension | Go `upload_test.go` | PASS |
| **TC-17** | Upload | Plain text file disguised with `.jpg` extension | Content sniffing detects `text/plain` and rejects with 400 | Go `upload_test.go` | PASS |
| **TC-18** | Upload | Upload file exceeding 5MB size limit | Request rejected with 400 Bad Request | Go `upload_test.go` | PASS |
| **TC-19** | Worker Pool | Asynchronous OCR and F# validation job dispatch | Upload returns 201 immediately with PENDING; background worker completes | Go `worker_pool.go` | PASS |
| **TC-20** | Worker Pool | OCR or F# microservice downtime resilience | Prescription defaults safely to `MANUAL_REVIEW`, never auto-approves | Go `worker_pool.go` | PASS |
| **TC-21** | F# Rules | Valid standard prescription (all fields complete) | Returns `Valid` union case with status `"valid"` | F# `Tests.fs` | PASS |
| **TC-22** | F# Rules | Valid multi-medicine prescription | Returns `Valid` with all medicines preserved | F# `Tests.fs` | PASS |
| **TC-23** | F# Rules | Boundary test: line quantity = 1 | Returns `Valid` | F# `Tests.fs` | PASS |
| **TC-24** | F# Rules | Boundary test: line quantity = 30 | Returns `Valid` | F# `Tests.fs` | PASS |
| **TC-25** | F# Rules | Missing patient name | Returns `Incomplete` with `"patientName"` in missing list | F# `Tests.fs` | PASS |
| **TC-26** | F# Rules | Missing doctor registration number | Returns `Incomplete` with `"doctorRegNo"` in missing list | F# `Tests.fs` | PASS |
| **TC-27** | F# Rules | Empty medicines list (`[]`) | Returns `Incomplete` with `"medicines"` in missing list | F# `Tests.fs` | PASS |
| **TC-28** | F# Rules | Future prescription date | Returns `Invalid` with issue: `"cannot be in the future"` | F# `Tests.fs` | PASS |
| **TC-29** | F# Rules | Prescription older than 6 months (> 180 days) | Returns `Invalid` with issue: `"older than 6 months"` | F# `Tests.fs` | PASS |
| **TC-30** | F# Rules | Line quantity = 0 | Returns `Invalid` with issue: `"Quantity must be greater than 0"` | F# `Tests.fs` | PASS |
| **TC-31** | F# Rules | Line quantity > 30 | Returns `Invalid` with issue: `"exceeds maximum 30 units per line"` | F# `Tests.fs` | PASS |
| **TC-32** | F# Rules | Blank medicine name | Returns `Invalid` with issue: `"cannot be blank"` | F# `Tests.fs` | PASS |
| **TC-33** | F# Rules | Blank dosage instructions | Returns `Invalid` with issue: `"Dosage instructions cannot be blank"` | F# `Tests.fs` | PASS |
| **TC-34** | F# Rules | Controlled opioid keyword ("morphine") | Flags `"Controlled substance 'morphine': requires pharmacist attention"` | F# `Tests.fs` | PASS |
| **TC-35** | F# Rules | Controlled keyword ("fentanyl") | Flags `"Controlled substance 'fentanyl': requires pharmacist attention"` | F# `Tests.fs` | PASS |
| **TC-36** | Ordering | Order OTC item with sufficient stock | Stock decremented, order created with status `PLACED` | Go `stock_logic_test.go` | PASS |
| **TC-37** | Ordering | Order quantity exceeds available stock | Rejected with `ErrInsufficientStock`; stock unchanged | Go `stock_logic_test.go` | PASS |
| **TC-38** | Ordering | Order Rx item without prescription ID | Rejected with `ErrRxPrescriptionReq` | Go `stock_logic_test.go` | PASS |
| **TC-39** | Ordering | Order Rx item with unapproved (PENDING) prescription | Rejected with `ErrRxPrescriptionReq` | Go `stock_logic_test.go` | PASS |
| **TC-40** | Ordering | Order Rx item with APPROVED prescription | Order accepted, stock decremented, status set to `PLACED` | Go `stock_logic_test.go` | PASS |
| **TC-41** | Ordering | Cancel order | Transitions to `CANCELLED`, restores reserved stock to inventory | Go `stock_logic_test.go` | PASS |
| **TC-42** | State Machine | Allowed transition: `PLACED` to `UNDER_REVIEW` | Permitted | Go `transitions_test.go` | PASS |
| **TC-43** | State Machine | Allowed transition: `UNDER_REVIEW` to `CONFIRMED` | Permitted | Go `transitions_test.go` | PASS |
| **TC-44** | State Machine | Allowed transition: `CONFIRMED` to `READY` | Permitted | Go `transitions_test.go` | PASS |
| **TC-45** | State Machine | Allowed transition: `READY` to `COMPLETED` | Permitted | Go `transitions_test.go` | PASS |
| **TC-46** | State Machine | Invalid transition: `PLACED` to `COMPLETED` (skipping steps) | Rejected with `ErrInvalidTransition` | Go `transitions_test.go` | PASS |
| **TC-47** | State Machine | Terminal state: `COMPLETED` to `CANCELLED` | Rejected | Go `transitions_test.go` | PASS |
| **TC-48** | State Machine | Terminal state: `CANCELLED` to `PLACED` | Rejected | Go `transitions_test.go` | PASS |
| **TC-49** | Frontend | Cart adds items, calculates subtotal, updates count | Verified in memory & localStorage | Vitest `cart.test.ts` | PASS |
| **TC-50** | Frontend | Cart dynamically detects when Rx item is present | `cart.requiresPrescription()` returns true | Vitest `cart.test.ts` | PASS |
| **TC-51** | Frontend | Cart removes item and recalculates Rx requirement | `cart.requiresPrescription()` returns false when Rx item removed | Vitest `cart.test.ts` | PASS |
| **TC-52** | Frontend | API client persists and removes JWT and user profile | LocalStorage correctly synced | Vitest `api.test.ts` | PASS |
| **TC-53** | Admin | Admin updates customer role to `staff` | User role updated; pharmacist capabilities unlocked | Go `handlers_test.go` | PASS |
