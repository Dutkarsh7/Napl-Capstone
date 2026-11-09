# Phase-II Capstone Project Documentation

## Project Title
AI-Powered Smart Pharmacy and Prescription Assistant

## 1. Introduction
The AI-Powered Smart Pharmacy and Prescription Assistant is a full-stack pharmacy management platform designed to support medicine discovery, prescription submission and verification, and order fulfilment. It uses multiple services written in different programming languages. The system is designed to assist pharmacy staff with digital workflows while keeping the final prescription approval under the control of a licensed pharmacist.

## 2. Objectives
- Provide a digital platform for customers to search and browse medicines.
- Allow customers to upload prescription documents for processing.
- Use OCR to extract text from uploaded prescription images and documents.
- Apply rule-based validation to identify missing or potentially invalid prescription information.
- Provide a pharmacist with a review queue to approve or reject prescriptions.
- Support medicine inventory, cart operations, order placement and order-status tracking.
- Demonstrate the use of Go, F#, and other technologies in a distributed application.

## 3. Key Features
- User authentication and role-based access for customers, pharmacists and administrators.
- Medicine catalogue with search and prescription-required indicators.
- Prescription upload with OCR-based text extraction.
- Separate F# validation service for prescription rule checks.
- Pharmacist dashboard for manual prescription review and approval or rejection.
- Shopping cart and order placement with prescription checks for restricted medicines.
- Inventory management and stock updates.
- Order fulfilment workflow with status tracking.
- PostgreSQL database with schema constraints, migrations and sample seed data.
- Automated tests for backend logic, validation rules and frontend behaviour.

## 4. Technologies Used
- **Go:** Main REST API, authentication, request handling, asynchronous prescription processing and order operations.
- **F#:** Prescription validation service using domain models and rule-based checks.
- **Python:** OCR microservice built with FastAPI, pytesseract and Pillow.
- **TypeScript, HTML and CSS:** Frontend application, built and served with Vite.
- **PostgreSQL:** Persistent data storage.
- **Docker and Docker Compose:** Containerisation and running the application services together.
- **Go testing, Expecto and Vitest:** Automated testing for the backend, F# validation service and frontend.

## 5. Implementation
The application is divided into separate services. The frontend communicates with the Go REST API. Customers can browse medicines and submit prescription files. The Go backend handles the upload and dispatches the processing task to an asynchronous worker. The worker communicates with the Python OCR service to extract text and then sends the extracted information to the F# validation service for rule-based checks.

The processing result is stored and shown to the pharmacist for manual review. A prescription is not automatically approved by OCR or validation. The pharmacist must review and approve it before it can be used for an order involving prescription-required medicines.

The backend also manages user roles, medicine information, inventory and order processing. Order fulfilment is handled through defined status transitions. PostgreSQL stores application data, while Docker Compose provides a way to run the database, backend, frontend and microservices together.

## 6. Application Workflow
1. A customer signs in and searches the medicine catalogue.
2. The customer uploads a prescription when required.
3. The OCR service extracts text from the uploaded document.
4. The F# service performs rule-based validation on the extracted information.
5. The pharmacist checks the original document and the processing result, then approves or rejects the prescription.
6. The customer can use an approved prescription when placing an order for prescription-required medicines.
7. The pharmacist updates the order through the fulfilment stages, and the customer can track the order status.

## 7. Testing
The repository includes automated tests for the main application components:
- Go tests for authentication, HTTP handlers, order and stock logic, state transitions and file uploads.
- F# Expecto tests for prescription validation rules.
- Vitest tests for frontend API and cart behaviour.

The README documents commands for running these test suites. Actual test execution should be performed in the submission environment before claiming a fresh passing test run.

## 8. Running the Project
The project can be started using Docker Compose from the repository root:

```bash
docker compose up --build
```

The documented local endpoints are:
- Frontend: http://localhost:5173
- Go backend: http://localhost:8080
- F# validation service: http://localhost:5005
- Python OCR service: http://localhost:8000

Refer to README.md for environment configuration, individual service commands and the detailed demonstration steps.

## 9. Safety and Limitations
This project is intended to assist pharmacy personnel with digital workflows. It does not diagnose illnesses, prescribe medicines, recommend treatments or change dosages. OCR and automated validation provide supporting information only. A licensed pharmacist must review and explicitly approve each prescription before it is used to fulfil a prescription-controlled order. If automated services are unavailable, the system is designed to leave the prescription for manual review rather than approving it automatically.

## 10. Conclusion
The project demonstrates how Go, F# and Python can be used together in a distributed application. It combines prescription document processing, rule-based validation, pharmacist review, medicine inventory and order management in one platform. The project also demonstrates the use of automated testing, database-backed services and Docker-based deployment.

## Repository
https://github.com/Dutkarsh7/Napl-Capstone
