package postgres

import (
    
    "time"
    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
    "ecommerce/services/user-service/internal/domain"
)

type UserRepository struct {
    db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
    query := `
        INSERT INTO users (id, email, password, first_name, last_name, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `
    user.ID = uuid.New()
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()
    user.Status = "active"
    
    _, err := r.db.Exec(query, user.ID, user.Email, user.Password, 
        user.FirstName, user.LastName, user.Status, user.CreatedAt, user.UpdatedAt)
    return err
}

func (r *UserRepository) GetByID(id uuid.UUID) (*domain.User, error) {
    var user domain.User
    query := `SELECT id, email, first_name, last_name, status, created_at, updated_at 
              FROM users WHERE id = $1 AND status != 'deleted'`
    err := r.db.Get(&user, query, id)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
    var user domain.User
    query := `SELECT id, email, password, first_name, last_name, status, created_at, updated_at 
              FROM users WHERE email = $1 AND status != 'deleted'`
    err := r.db.Get(&user, query, email)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *UserRepository) Update(user *domain.User) error {
    query := `UPDATE users SET first_name = $2, last_name = $3, updated_at = $4 WHERE id = $1`
    user.UpdatedAt = time.Now()
    _, err := r.db.Exec(query, user.ID, user.FirstName, user.LastName, user.UpdatedAt)
    return err
}

func (r *UserRepository) Delete(id uuid.UUID) error {
    query := `UPDATE users SET status = 'deleted', updated_at = $2 WHERE id = $1`
    _, err := r.db.Exec(query, id, time.Now())
    return err
}
