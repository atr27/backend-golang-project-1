# Hospital EMR System - API Documentation

## Overview

This document provides comprehensive API documentation for the Hospital EMR System.

**Base URL**: `http://localhost:8080/api/v1`

**Authentication**: All protected endpoints require a Bearer token in the Authorization header.

```
Authorization: Bearer <your_jwt_token>
```

## Table of Contents

1. [Authentication](#authentication)
2. [Patients](#patients)
3. [Encounters](#encounters)
4. [Appointments](#appointments)
5. [Orders](#orders)
6. [Results](#results)

---

## Authentication

### Login

Authenticate a user and receive JWT tokens.

**Endpoint**: `POST /auth/login`

**Request Body**:
```json
{
  "email": "doctor@hospital-emr.com",
  "password": "doctor123",
  "mfa_code": "123456"
}
```

**Response**: `200 OK`
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 86400,
  "user": {
    "id": "uuid",
    "email": "doctor@hospital-emr.com",
    "first_name": "John",
    "last_name": "Smith",
    "roles": [
      {
        "name": "Doctor",
        "code": "doctor"
      }
    ]
  }
}
```

### Logout

Invalidate the current session.

**Endpoint**: `POST /auth/logout`

**Headers**: `Authorization: Bearer <token>`

**Response**: `200 OK`
```json
{
  "message": "Logged out successfully"
}
```

### Refresh Token

Get a new access token using refresh token.

**Endpoint**: `POST /auth/refresh`

**Request Body**:
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response**: `200 OK`
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 86400
}
```

---

## Patients

### List Patients

Get a paginated list of patients.

**Endpoint**: `GET /patients`

**Query Parameters**:
- `page` (integer, optional): Page number (default: 1)
- `page_size` (integer, optional): Items per page (default: 20, max: 100)
- `search` (string, optional): Search by name, MRN, or email

**Headers**: `Authorization: Bearer <token>`

**Response**: `200 OK`
```json
{
  "data": [
    {
      "id": "uuid",
      "mrn": "MRN000001",
      "first_name": "Alice",
      "last_name": "Johnson",
      "date_of_birth": "1985-05-15T00:00:00Z",
      "gender": "female",
      "blood_type": "A+",
      "email": "alice.johnson@email.com",
      "phone_number": "+1234567892",
      "status": "active"
    }
  ],
  "total": 100,
  "page": 1,
  "page_size": 20,
  "total_pages": 5
}
```

### Create Patient

Register a new patient.

**Endpoint**: `POST /patients`

**Headers**: `Authorization: Bearer <token>`

**Request Body**:
```json
{
  "first_name": "Alice",
  "last_name": "Johnson",
  "middle_name": "Marie",
  "date_of_birth": "1985-05-15T00:00:00Z",
  "gender": "female",
  "blood_type": "A+",
  "marital_status": "married",
  "nationality": "American",
  "religion": "Christian",
  "ssn": "123-45-6789",
  "email": "alice.johnson@email.com",
  "phone_number": "+1234567892",
  "mobile_number": "+1234567892",
  "address": "123 Main St",
  "city": "New York",
  "state": "NY",
  "zip_code": "10001",
  "country": "USA",
  "emergency_contact": {
    "name": "Bob Johnson",
    "relationship": "Spouse",
    "phone_number": "+1234567893",
    "email": "bob@email.com"
  },
  "insurance": {
    "provider": "Blue Cross",
    "policy_number": "BC123456",
    "group_number": "GRP001",
    "expiry_date": "2025-12-31",
    "coverage_type": "Full"
  },
  "language": "English",
  "occupation": "Engineer"
}
```

**Response**: `201 Created`
```json
{
  "id": "uuid",
  "mrn": "MRN000123",
  "first_name": "Alice",
  "last_name": "Johnson",
  ...
}
```

### Get Patient

Get detailed information about a specific patient.

**Endpoint**: `GET /patients/:id`

**Headers**: `Authorization: Bearer <token>`

**Response**: `200 OK`
```json
{
  "id": "uuid",
  "mrn": "MRN000001",
  "first_name": "Alice",
  "last_name": "Johnson",
  "date_of_birth": "1985-05-15T00:00:00Z",
  "gender": "female",
  "blood_type": "A+",
  "allergies": [
    {
      "id": "uuid",
      "allergy_type": "drug",
      "allergen": "Penicillin",
      "reaction": "Rash",
      "severity": "moderate"
    }
  ],
  "medications": [
    {
      "id": "uuid",
      "medication_name": "Lisinopril",
      "dosage": "10mg",
      "frequency": "Once daily",
      "status": "active"
    }
  ]
}
```

### Update Patient

Update patient information.

**Endpoint**: `PUT /patients/:id`

**Headers**: `Authorization: Bearer <token>`

**Request Body**: Same as Create Patient

**Response**: `200 OK`

### Delete Patient

Soft delete a patient.

**Endpoint**: `DELETE /patients/:id`

**Headers**: `Authorization: Bearer <token>`

**Response**: `204 No Content`

### Get Patient Timeline

Get complete medical timeline for a patient.

**Endpoint**: `GET /patients/:id/timeline`

**Headers**: `Authorization: Bearer <token>`

**Response**: `200 OK`
```json
{
  "patient": { ... },
  "encounters": [ ... ],
  "appointments": [ ... ],
  "allergies": [ ... ],
  "medications": [ ... ]
}
```

---

## Error Responses

All endpoints may return the following error responses:

### 400 Bad Request
```json
{
  "code": "BAD_REQUEST",
  "message": "Bad request",
  "status_code": 400,
  "details": "Validation error details"
}
```

### 401 Unauthorized
```json
{
  "code": "UNAUTHORIZED",
  "message": "Unauthorized access",
  "status_code": 401
}
```

### 403 Forbidden
```json
{
  "code": "FORBIDDEN",
  "message": "Access forbidden",
  "status_code": 403
}
```

### 404 Not Found
```json
{
  "code": "PATIENT_NOT_FOUND",
  "message": "Patient with ID xyz not found",
  "status_code": 404
}
```

### 500 Internal Server Error
```json
{
  "code": "INTERNAL_ERROR",
  "message": "Internal server error",
  "status_code": 500
}
```

---

## Rate Limiting

API requests are limited to 100 requests per minute per user. If you exceed this limit, you will receive a `429 Too Many Requests` response.

---

## Pagination

List endpoints support pagination with the following parameters:
- `page`: Page number (starts at 1)
- `page_size`: Number of items per page (default: 20, max: 100)

Response includes:
- `data`: Array of items
- `total`: Total number of items
- `page`: Current page number
- `page_size`: Items per page
- `total_pages`: Total number of pages

---

## Filtering and Searching

List endpoints support searching with the `search` query parameter. The search is case-insensitive and searches across multiple fields.

Example:
```
GET /api/v1/patients?search=alice&page=1&page_size=20
```

---

## Timestamps

All timestamps are in UTC and follow ISO 8601 format:
```
2025-01-06T10:30:00Z
```

---

## FHIR Compatibility

The system is designed to be FHIR R4 compatible. FHIR endpoints will be available at:
```
/fhir/r4/Patient
/fhir/r4/Encounter
/fhir/r4/Observation
```

---

## Swagger Documentation

Interactive API documentation is available at:
```
http://localhost:8080/swagger/index.html
```

---

## Support

For API support, contact: support@hospital-emr.com
