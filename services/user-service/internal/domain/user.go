package domain

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    ID        uuid.UUID `json:"id" db:"id"`
    Email     string    `json:"email" db:"email"`
    Password  string    `json:"-" db:"password"`
    FirstName string    `json:"first_name" db:"first_name"`
    LastName  string    `json:"last_name" db:"last_name"`
    Status    string    `json:"status" db:"status"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UserRepository interface {
    Create(user *User) error
    GetByID(id uuid.UUID) (*User, error)
    GetByEmail(email string) (*User, error)
    Update(user *User) error
    Delete(id uuid.UUID) error
}

type UserService interface {
    Register(email, password, firstName, lastName string) (*User, error)
    Login(email, password string) (string, error)
    GetProfile(userID uuid.UUID) (*User, error)
    UpdateProfile(userID uuid.UUID, firstName, lastName string) error
}
