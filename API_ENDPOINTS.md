# Trimana Dashboard API Endpoints

## Base URL
Your Lambda API Gateway URL (e.g., `https://your-api-id.execute-api.us-west-2.amazonaws.com/prod`)

## Authentication

### Login
**POST** `/auth/login`

Login with username and password.

**Request Body:**
```json
{
  "username": "string",
  "password": "string"
}
```

**Response (200):**
```json
{
  "username": "string",
  "token": "string"
}
```

**Response (401):**
```json
{
  "statusCode": 401,
  "result": "Invalid credentials"
}
```

---

## Employees

### Get All Employees
**GET** `/employees`

Returns all employees in the system.

**Response (200):**
```json
[
  {
    "id": "string",
    "branchID": "string",
    "name": "string",
    "employeeID": "string",
    "hourlyRate": 20.0,
    "payFrequency": "B",
    "earningsCode": "REG",
    "rateCode": "BASE",
    "cardID": "string",
    "email": "string"
  }
]
```

### Get Employee by ID
**GET** `/employees/:id`

Returns a single employee by ID.

**Response (200):**
```json
{
  "id": "string",
  "branchID": "string",
  "name": "string",
  "employeeID": "string",
  "hourlyRate": 20.0,
  "payFrequency": "B",
  "earningsCode": "REG",
  "rateCode": "BASE",
  "cardID": "string",
  "email": "string"
}
```

**Response (404):**
```json
{
  "statusCode": 404,
  "result": "Employee not found"
}
```

### Create Employee
**POST** `/employees`

Creates a new employee.

**Request Body:**
```json
{
  "branchID": "string",
  "name": "string",
  "employeeID": "string",
  "hourlyRate": 20.0,
  "payFrequency": "B",
  "earningsCode": "REG",
  "rateCode": "BASE",
  "cardID": "string",
  "email": "string"
}
```

**Response (201):**
```json
{
  "id": "generated-id",
  "branchID": "string",
  "name": "string",
  "employeeID": "string",
  "hourlyRate": 20.0,
  "payFrequency": "B",
  "earningsCode": "REG",
  "rateCode": "BASE",
  "cardID": "string",
  "email": "string"
}
```

### Update Employee
**PUT** `/employees/:id`

Updates an existing employee.

**Request Body:**
```json
{
  "branchID": "string",
  "name": "string",
  "employeeID": "string",
  "hourlyRate": 20.0,
  "payFrequency": "B",
  "earningsCode": "REG",
  "rateCode": "BASE",
  "cardID": "string",
  "email": "string"
}
```

**Response (200):**
```json
{
  "id": "string",
  "branchID": "string",
  "name": "string",
  "employeeID": "string",
  "hourlyRate": 20.0,
  "payFrequency": "B",
  "earningsCode": "REG",
  "rateCode": "BASE",
  "cardID": "string",
  "email": "string"
}
```

### Delete Employee
**DELETE** `/employees/:id`

Deletes an employee by ID.

**Response (200):**
```json
{
  "result": "Employee deleted successfully"
}
```

---

## Timecards

### Get All Timecards
**GET** `/timecards`

Returns all timecards. Can filter by employee ID using query parameter.

**Query Parameters:**
- `employeeID` (optional): Filter by employee ID

**Response (200):**
```json
[
  {
    "id": "ObjectID",
    "employeeID": "string",
    "createdAt": "2025-08-27T07:00:00Z",
    "timeBlocks": [
      {
        "checkIn": "2025-08-27T13:06:35Z",
        "checkOut": "2025-08-27T16:53:08Z"
      }
    ]
  }
]
```

### Get Timecard by ID
**GET** `/timecards/:id`

Returns a single timecard by ID.

**Response (200):**
```json
{
  "id": "ObjectID",
  "employeeID": "string",
  "createdAt": "2025-08-27T07:00:00Z",
  "timeBlocks": [
    {
      "checkIn": "2025-08-27T13:06:35Z",
      "checkOut": "2025-08-27T16:53:08Z"
    }
  ]
}
```

**Response (404):**
```json
{
  "statusCode": 404,
  "result": "Timecard not found"
}
```

### Create Timecard
**POST** `/timecards`

Creates a new timecard manually.

**Request Body:**
```json
{
  "employeeID": "string",
  "createdAt": "2025-08-27T07:00:00Z",
  "timeBlocks": [
    {
      "checkIn": "2025-08-27T13:06:35Z",
      "checkOut": "2025-08-27T16:53:08Z"
    }
  ]
}
```

**Response (201):**
```json
{
  "id": "generated-ObjectID",
  "employeeID": "string",
  "createdAt": "2025-08-27T07:00:00Z",
  "timeBlocks": [
    {
      "checkIn": "2025-08-27T13:06:35Z",
      "checkOut": "2025-08-27T16:53:08Z"
    }
  ]
}
```

### Update Timecard
**PUT** `/timecards/:id`

Updates an existing timecard.

**Request Body:**
```json
{
  "employeeID": "string",
  "createdAt": "2025-08-27T07:00:00Z",
  "timeBlocks": [
    {
      "checkIn": "2025-08-27T13:06:35Z",
      "checkOut": "2025-08-27T16:53:08Z"
    }
  ]
}
```

**Response (200):**
```json
{
  "id": "ObjectID",
  "employeeID": "string",
  "createdAt": "2025-08-27T07:00:00Z",
  "timeBlocks": [
    {
      "checkIn": "2025-08-27T13:06:35Z",
      "checkOut": "2025-08-27T16:53:08Z"
    }
  ]
}
```

### Delete Timecard
**DELETE** `/timecards/:id`

Deletes a timecard by ID.

**Response (200):**
```json
{
  "result": "Timecard deleted successfully"
}
```

---

## Existing Payroll Endpoints

### Payroll Event (Card Scan)
**POST** `/payroll/event`

Handles employee card scan for check-in/check-out.

**Request Body:**
```json
{
  "cardID": "string"
}
```

### Payroll Report
**POST** `/payroll/report`

Generates and emails the payroll report for the current pay period.

**Response (200):**
```json
{
  "statusCode": 200,
  "result": "OK"
}
```

---

## Setup

### Database Setup

Before using the API, ensure you have a `users` collection in your MongoDB database with at least one user:

```json
{
  "_id": "user1",
  "username": "admin",
  "password": "your-password"
}
```

**Note:** In production, passwords should be hashed using bcrypt or similar.

### CORS

All endpoints include CORS headers:
```
Access-Control-Allow-Origin: *
```

### Error Responses

All error responses follow this format:
```json
{
  "statusCode": 400-500,
  "result": "Error message description"
}
```

