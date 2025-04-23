# Cars API with Authentication

This project implements a RESTful API for car management with user authentication using Gin and GORM.

## Authentication

The system uses JWT (JSON Web Token) for authentication. When a user registers or logs in, they receive a token that must be included in the Authorization header for protected routes.

## Getting Started

1. Run the application:
   ```
   go run cmd/main.go
   ```

2. The server will start on port 8081.

## Testing with Postman

A Postman collection is included in the project (`postman_collection.json`). Follow these steps to use it:

1. Open Postman
2. Import the collection file: `postman_collection.json`
3. Create a new environment in Postman and set a variable named `token` (it will be automatically populated when you login or register)

### Authentication Requests

#### Register a New User
- **URL**: `POST http://localhost:8081/auth/register`
- **Body**:
  ```json
  {
    "name": "Test User",
    "email": "test@example.com",
    "password": "password123"
  }
  ```
- This will return a token and user information

#### Login
- **URL**: `POST http://localhost:8081/auth/login`
- **Body**:
  ```json
  {
    "email": "test@example.com",
    "password": "password123"
  }
  ```
- This will return a token and user information

#### Get Current User
- **URL**: `GET http://localhost:8081/auth/me`
- **Headers**: `Authorization: Bearer {{token}}`
- Returns information about the currently authenticated user

### User CRUD Operations

#### Get All Users
- **URL**: `GET http://localhost:8081/users`
- **Headers**: `Authorization: Bearer {{token}}`

#### Create User
- **URL**: `POST http://localhost:8081/users`
- **Headers**: `Authorization: Bearer {{token}}`
- **Body**:
  ```json
  {
    "name": "New User",
    "email": "newuser@example.com",
    "password": "password123"
  }
  ```

#### Update User
- **URL**: `PUT http://localhost:8081/users/1` (replace 1 with the user ID)
- **Headers**: `Authorization: Bearer {{token}}`
- **Body**:
  ```json
  {
    "name": "Updated User Name",
    "email": "updated@example.com"
  }
  ```

#### Delete User
- **URL**: `DELETE http://localhost:8081/users/1` (replace 1 with the user ID)
- **Headers**: `Authorization: Bearer {{token}}`

### Car CRUD Operations

#### Get All Cars
- **URL**: `GET http://localhost:8081/cars`
- **Headers**: `Authorization: Bearer {{token}}`

#### Create Car
- **URL**: `POST http://localhost:8081/cars`
- **Headers**: `Authorization: Bearer {{token}}`
- **Body**:
  ```json
  {
    "make": "Toyota",
    "model": "Camry",
    "year": 2023
  }
  ```

#### Update Car
- **URL**: `PUT http://localhost:8081/cars/1` (replace 1 with the car ID)
- **Headers**: `Authorization: Bearer {{token}}`
- **Body**:
  ```json
  {
    "make": "Toyota",
    "model": "Corolla",
    "year": 2024
  }
  ```

#### Delete Car
- **URL**: `DELETE http://localhost:8081/cars/1` (replace 1 with the car ID)
- **Headers**: `Authorization: Bearer {{token}}`

## Notes

- The token is automatically saved to your Postman environment when you register or login
- All protected routes require the Authorization header with a valid token
- For security in production, make sure to store the JWT secret in environment variables
