# User Authentication API

This is a robust user authentication API built with Go and Gin framework. It provides secure user management functionalities with features like rate limiting, account lockout, and email verification.

## Features

- User registration with email verification
- Secure login with account lockout protection
- Password reset functionality
- User profile management
- Role-based access control
- Rate limiting to prevent abuse

## API Endpoints

### Public Endpoints

- `POST /register`: Register a new user
- `POST /login`: Authenticate a user
- `POST /reset-request`: Request a password reset
- `POST /reset-password`: Reset password with a valid token
- `POST /verify-email`: Verify user's email address

### Protected Endpoints (Require Authentication)

- `POST /logout`: Log out the current user
- `GET /protected`: A sample protected route
- `GET /user`: Get current user's profile
- `GET /users`: Get all users (admin only)
- `GET /users/:id`: Get a specific user's profile
- `PUT /users/:id`: Update a user's profile
- `DELETE /users/:id`: Delete a user (admin only)

## Implementation Guide

To integrate this API with your UI:

1. **User Registration**:
   - Send a POST request to `/register` with username, email, and password.
   - Implement Cloudflare Turnstile for CAPTCHA protection.
   - Handle the email verification process.

2. **User Login**:
   - Send a POST request to `/login` with email/username and password.
   - Store the returned JWT token securely (e.g., in HttpOnly cookies).

3. **Protected Routes**:
   - Include the JWT token in the Authorization header for all protected routes.

4. **Password Reset**:
   - Implement the forgot password flow using `/reset-request` and `/reset-password` endpoints.

5. **User Management**:
   - Use the appropriate endpoints for viewing and editing user profiles.

6. **Error Handling**:
   - Implement proper error handling for various HTTP status codes returned by the API.

## Security Considerations

- Use HTTPS in production.
- Implement proper CORS settings.
- Never store or transmit passwords in plain text.
- Use secure session management techniques.

## Rate Limiting

- General rate limit: 60 requests per minute per IP.
- Login: 5 requests per minute per IP.
- Registration: 3 requests per minute per IP.
- Password reset request: 2 requests per minute per IP.

## Setup and Configuration

1. Clone the repository.
2. Set up environment variables (see `.env.example`).
3. Run `go mod tidy` to install dependencies.
4. Start the server with `go run main.go`.

## Testing

Run the test suite with:

```bash
go test ./...

