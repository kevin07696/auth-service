package domain

import (
	"context"
)

// Defines an abstraction for database operations related to the Login entity.
type Repositor interface {
    // Applies automatic table and index migrations for the Login entity.
    Migrate()

    // Inserts a new login record into the database.
    // Returns a status code.
    CreateLogin(ctx context.Context, login *Login) (status StatusCode)

    // Reads a login record by username.
    // Returns the login record and a status code.
    GetLoginByUsername(ctx context.Context, username Username) (login Login, status StatusCode)

    // Reads a login record by email address.
    // Returns the login record and a status code.
    GetLoginByEmail(ctx context.Context, emailAddress Email) (login Login, status StatusCode)
}

// The unvalidated service payload used to create a new login.
type CreateLoginRequest struct {
    Username string `json:"username" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

// The unvalidated service payload used to verify a login.
type LoginRequest struct {
	UserInput string
	Password  string
}

// The login service response
type LoginResponse struct {
	LoginID string
}
