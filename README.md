### Loan Management System Documentation

#### Overview

This documentation provides an overview of the Loan Management System, including configuration settings, routes, and examples for each API endpoint.

#### Environment Configuration

Ensure the following environment variables are set in your `.env` file:

```env
SERVER_PORT = ":8080"
SERVER_URL = "127.0.0.1"

MONGODB_URL = "mongodb+srv://<username>:<password>@cluster0.bd7iqjq.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

JWT_SECRET = "secret"
JWT_REFRESH_TOKEN_SECRET = "REFG55"
ACCESS_TOKEN_EXPIRY_HOUR = 2
REFRESH_TOKEN_EXPIRY_HOUR = 168
```

Replace `<username>` and `<password>` with your actual MongoDB credentials.

#### Routes Setup

The application routes are organized into different modules: User, Admin, Loan, and Log.

##### User Routes

**Path:** `/users`

**Endpoints:**

1. **Register User**
   - **Method:** `POST`
   - **Endpoint:** `/users/register`
   - **Example Request:**
     ```json
     {
       "name": "John Doe",
       "email": "john.doe@example.com",
       "password": "password123"
     }
     ```

2. **Login**
   - **Method:** `POST`
   - **Endpoint:** `/users/login`
   - **Example Request:**
     ```json
     {
       "email": "john.doe@example.com",
       "password": "password123"
     }
     ```

3. **Get My Profile**
   - **Method:** `GET`
   - **Endpoint:** `/users/profile`
   - **Headers:**
     - `Authorization: Bearer <access_token>`
   - **Example Response:**
     ```json
     {
       "id": "601c9e5e68c4e031c8b9bdeb",
       "name": "John Doe",
       "email": "john.doe@example.com"
     }
     ```

##### Admin Routes

**Path:** `/admin`

**Endpoints:**

1. **Get Users**
   - **Method:** `GET`
   - **Endpoint:** `/admin/users`
   - **Headers:**
     - `Authorization: Bearer <access_token>`
   - **Example Response:**
     ```json
     [
       {
         "id": "601c9e5e68c4e031c8b9bdeb",
         "name": "John Doe",
         "email": "john.doe@example.com"
       },
       {
         "id": "601c9e5e68c4e031c8b9bded",
         "name": "Jane Smith",
         "email": "jane.smith@example.com"
       }
     ]
     ```

2. **Delete User**
   - **Method:** `DELETE`
   - **Endpoint:** `/admin/users/:id`
   - **Headers:**
     - `Authorization: Bearer <access_token>`
   - **Example Request:**
     ```
     DELETE /admin/users/601c9e5e68c4e031c8b9bded
     ```

##### Loan Routes

**Path:** `/loans`

**Endpoints:**

1. **Create Loan**
   - **Method:** `POST`
   - **Endpoint:** `/loans`
   - **Headers:**
     - `Authorization: Bearer <access_token>`
   - **Example Request:**
     ```json
     {
       "user_id": "601c9e5e68c4e031c8b9bdeb",
       "amount": 5000,
       "interest_rate": 5,
       "term": 12,
       "reason": "Home renovation"
     }
     ```

2. **View Loan Status**
   - **Method:** `GET`
   - **Endpoint:** `/loans/:loanID`
   - **Headers:**
     - `Authorization: Bearer <access_token>`
   - **Example Request:**
     ```
     GET /loans/601c9e5e68c4e031c8b9bded
     ```

3. **View All Loans (Admin)**
   - **Method:** `GET`
   - **Endpoint:** `/admin/loans`
   - **Headers:**
     - `Authorization: Bearer <access_token>`
   - **Query Parameters:**
     - `status` (optional): Filter by loan status (`pending`, `approved`, `rejected`)
     - `order` (optional): Sort order (`asc`, `desc`)
     - `limit` (optional): Number of records per page (default: 10)
     - `offset` (optional): Pagination offset (default: 0)
   - **Example Request:**
     ```
     GET /admin/loans?status=pending&order=asc&limit=5&offset=0
     ```

4. **Update Loan Status (Admin)**
   - **Method:** `PATCH`
   - **Endpoint:** `/admin/loans/:loanID/:status`
   - **Headers:**
     - `Authorization: Bearer <access_token>`
   - **Example Request:**
     ```
     PATCH /admin/loans/601c9e5e68c4e031c8b9bded/approved
     ```

5. **Delete Loan (Admin)**
   - **Method:** `DELETE`
   - **Endpoint:** `/admin/loans/:loanID`
   - **Headers:**
     - `Authorization: Bearer <access_token>`
   - **Example Request:**
     ```
     DELETE /admin/loans/601c9e5e68c4e031c8b9bded
     ```

##### Log Routes

**Path:** `/admin/logs`

**Endpoints:**

1. **View Logs (Admin)**
   - **Method:** `GET`
   - **Endpoint:** `/admin/logs`
   - **Headers:**
     - `Authorization: Bearer <access_token>`
   - **Query Parameters:**
     - `event` (optional): Filter by event type (`create_loan`, `update_loan_status`, etc.)
     - `order` (optional): Sort order (`asc`, `desc`)
     - `limit` (optional): Number of records per page (default: 10)
     - `offset` (optional): Pagination offset (default: 0)
   - **Example Request:**
     ```
     GET /admin/logs?event=view_logs&order=asc&limit=5&offset=0
     ```

   - **Example Response:**
     ```json
     {
       "logs": [
         {
           "event": "view_logs",
           "details": "Logs viewed by Admin ID: 601c9e5e68c4e031c8b9bdeb",
           "user_id": "601c9e5e68c4e031c8b9bdeb",
           "timestamp": "2024-08-27T10:15:30Z"
         }
       ],
       "current_page": 0,
       "per_page": 5,
       "total": 1,
       "total_pages": 1
     }
     ```

