## Environment Setup

Before running the application, ensure the following environment variables are set:

- `SERVER_PORT`: Port on which the server will run (e.g., `":8080"`).
- `SERVER_URL`: Base URL of the server (e.g., `"127.0.0.1"`).
- `MONGODB_URL`: MongoDB connection string (e.g., `"mongodb+srv://<username>:<password>@cluster.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"`).
- `JWT_SECRET`: Secret key for signing JWT tokens (e.g., `"Group3"`).
- `JWT_REFRESH_TOKEN_SECRET`: Secret key for signing JWT refresh tokens (e.g., `"REFG55"`).
- `ACCESS_TOKEN_EXPIRY_HOUR`: Expiry time for access tokens in hours (e.g., `2`).
- `REFRESH_TOKEN_EXPIRY_HOUR`: Expiry time for refresh tokens in hours (e.g., `168`).
- `RATE_LIMIT_MAX_REQUEST`: Maximum number of requests allowed within the specified time window (e.g., `10`).
- `RATE_LIMIT_EXPIRATION_MINUTE`: Expiration time for the rate limit window in minutes (e.g., `1`).

## Endpoints Documentation

### 1. **Login**

**Endpoint:** `/login`  
**Method:** `POST`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
- **200 OK**
  ```json
  {
    "tokens": {
      "access_token": "string",
      "refresh_token": "string"
    }
  }
  ```
- **400 Bad Request**
  ```json
  {
    "error": "Invalid request data"
  }
  ```
- **401 Unauthorized**
  ```json
  {
    "error": "Invalid email or password"
  }
  ```

**Example:**
```shell
curl -X POST http://127.0.0.1:8080/login -H "Content-Type: application/json" -d '{"email": "user@example.com", "password": "password123"}'
```

---

### 2. **Refresh Token**

**Endpoint:** `/refresh-token`  
**Method:** `POST`

**Request Body:**
```json
{
  "user_id": "1234567890",
  "token": "refresh_token_string"
}
```

**Response:**
- **200 OK**
  ```json
  {
    "tokens": {
      "access_token": "new_access_token",
      "refresh_token": "new_refresh_token"
    }
  }
  ```
- **400 Bad Request**
  ```json
  {
    "error": "Invalid request data"
  }
  ```
- **401 Unauthorized**
  ```json
  {
    "error": "Invalid or expired refresh token"
  }
  ```

**Example:**
```shell
curl -X POST http://127.0.0.1:8080/refresh-token -H "Content-Type: application/json" -d '{"user_id": "1234567890", "token": "refresh_token_string"}'
```

---

### 3. **Register**

**Endpoint:** `/register`  
**Method:** `POST`

**Request Body:**
```json
{
  "username": "newuser",
  "email": "newuser@example.com",
  "password": "password123"
}
```

**Response:**
- **200 OK**
  ```json
  {
    "message": "Registered successfully. Please check your email for account activation."
  }
  ```
- **400 Bad Request**
  ```json
  {
    "error": "Invalid request data"
  }
  ```
- **409 Conflict**
  ```json
  {
    "error": "Email already exists"
  }
  ```

**Example:**
```shell
curl -X POST http://127.0.0.1:8080/register -H "Content-Type: application/json" -d '{"username": "newuser", "email": "newuser@example.com", "password": "password123"}'
```

---

### 4. **Activate Account**

**Endpoint:** `/activate/:email/:token`  
**Method:** `GET`

**Response:**
- **200 OK**
  ```json
  {
    "message": "Account activated successfully"
  }
  ```
- **400 Bad Request**
  ```json
  {
    "error": "Invalid activation token"
  }
  ```
- **404 Not Found**
  ```json
  {
    "error": "Account not found"
  }
  ```

**Example:**
```shell
curl -X GET http://127.0.0.1:8080/activate/user@example.com/sometoken
```

---

### 5. **Get My Profile**

**Endpoint:** `/profile`  
**Method:** `GET`

**Headers:**  
- `Authorization: Bearer <access_token>`

**Response:**
- **200 OK**
  ```json
  {
    "id": "1234567890",
    "username": "user",
    "email": "user@example.com",
    "name": "User Name",
    "bio": "This is my bio",
    "role": "user",
    "is_active": true
  }
  ```
- **401 Unauthorized**
  ```json
  {
    "error": "Invalid or expired token"
  }
  ```

**Example:**
```shell
curl -X GET http://127.0.0.1:8080/profile -H "Authorization: Bearer access_token"
```

---

### 6. **Password Reset**

**Endpoint:** `/password-reset`  
**Method:** `POST`

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Response:**
- **200 OK**
  ```json
  {
    "status": 200,
    "message": "Successfully sent password reset link to your email"
  }
  ```
- **400 Bad Request**
  ```json
  {
    "error": "Invalid input"
  }
  ```
- **404 Not Found**
  ```json
  {
    "error": "Email not found"
  }
  ```

**Example:**
```shell
curl -X POST http://127.0.0.1:8080/password-reset -H "Content-Type: application/json" -d '{"email": "user@example.com"}'
```

---

### 7. **Update Password**

**Endpoint:** `/update-password`  
**Method:** `POST`

**Request Body:**
```json
{
  "user_id": "1234567890",
  "new_password": "newpassword123"
}
```

**Response:**
- **200 OK**
  ```json
  {
    "message": "Password has been reset"
  }
  ```
- **400 Bad Request**
  ```json
  {
    "error": "Invalid input"
  }
  ```
- **404 Not Found**
  ```json
  {
    "error": "User not found"
  }
  ```

**Example:**
```shell
curl -X POST http://127.0.0.1:8080/update-password -H "Content-Type: application/json" -d '{"user_id": "1234567890", "new_password": "newpassword123"}'
```

---

### 8. **Get Users (Admin Only)**

**Endpoint:** `/users`  
**Method:** `GET`

**Headers:**  
- `Authorization: Bearer <admin_access_token>`

**Response:**
- **200 OK**
  ```json
  {
    "users": [
      {
        "id": "1234567890",
        "username": "user",
        "email": "user@example.com",
        "role": "user"
      },
      ...
    ]
  }
  ```
- **401 Unauthorized**
  ```json
  {
    "error": "Unauthorized"
  }
  ```

**Example:**
```shell
curl -X GET http://127.0.0.1:8080/users -H "Authorization: Bearer admin_access_token"
```

---

### 9. **Delete User (Admin Only)**

**Endpoint:** `/users/:id`  
**Method:** `DELETE`

**Headers:**  
- `Authorization: Bearer <admin_access_token>`

**Response:**
- **200 OK**
  ```json
  {
    "message": "User deleted successfully",
    "user": {
      "id": "1234567890",
      "username": "user",
      "email": "user@example.com"
    }
  }
  ```
- **401 Unauthorized**
  ```json
  {
    "error": "Unauthorized"
  }
  ```
- **404 Not Found**
  ```json
  {
    "error": "User not found"
  }
  ```

**Example:**
```shell
curl -X DELETE http://127.0.0.1:8080/users/1234567890 -H "Authorization: Bearer admin_access_token"
```