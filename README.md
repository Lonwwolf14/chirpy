# Chirpy API

The Chirpy API is a simple API for user authentication, chirp management, and request metrics tracking. It includes endpoints for creating and managing users, creating and deleting chirps, and viewing metrics as an admin. The API also implements key features like JWT-based authentication, profanity filtering, and rate limiting.

## Table of Contents

- [API Endpoints](#api-endpoints)
  - [Users](#users)
  - [Chirps](#chirps)
  - [Health Check](#health-check)
  - [Metrics](#metrics)
- [Key Features](#key-features)

## API Endpoints

### Users

- **POST** `/api/users`  
  Create a new user.

- **PUT** `/api/users`  
  Update an existing user.

- **POST** `/api/login`  
  Login a user. Returns a JWT token for authenticated access.

- **POST** `/api/refresh`  
  Refresh the JWT token when expired.

- **POST** `/api/revoke`  
  Revoke the refresh token.

### Chirps

- **GET** `/api/chirps`  
  Retrieve all chirps.

- **GET** `/api/chirps/{id}`  
  Retrieve a specific chirp by ID.

- **POST** `/api/chirps`  
  Create a new chirp. All chirps are filtered for profanity.

- **DELETE** `/api/chirps/{id}`  
  Delete a chirp by ID.

### Health Check

- **GET** `/api/healthz`  
  Check the health of the API.

### Metrics

- **GET** `/admin/metrics`  
  View request metrics (admin only). Requires admin authentication.

## Key Features

- **User Authentication**: Implements JWT authentication for secure access to user data and chirps.
- **Profanity Filtering**: All chirps are filtered for profanity before being saved.
- **Request Metrics Tracking**: Tracks and provides insights into the API's usage (admin only).
- **Database Persistence**: Data is stored in a local file-based database (`database.json`).
- **Admin Dashboard**: Provides administrative functionality such as viewing request metrics.
- **Rate Limiting**: Protects the API from abuse by limiting the number of requests a user can make in a given time.
- **Token Management**: Supports token refreshing and revocation for secure session management.
