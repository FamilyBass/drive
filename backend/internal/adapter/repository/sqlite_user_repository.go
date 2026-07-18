package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/familybass/drive/internal/domain/entity"
)

// SQLiteUserRepository implémente UserRepository avec SQLite
type SQLiteUserRepository struct {
	db *sql.DB
}

// NewSQLiteUserRepository crée un nouveau repository SQLite pour les utilisateurs
func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

// Create insère un nouvel utilisateur
func (r *SQLiteUserRepository) Create(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users(id,email,password,is_active,is_admin,created_at) 
	          VALUES(?,?,?,?,?,?)`
	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.Password,
		boolToInt(user.IsActive), boolToInt(user.IsAdmin),
		user.CreatedAt.Unix(),
	)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

// GetByID récupère un utilisateur par son ID
func (r *SQLiteUserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	user := &entity.User{}
	var isActive, isAdmin int

	query := `SELECT id,email,password,is_active,is_admin,created_at FROM users WHERE id=?`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Password,
		&isActive, &isAdmin, (*unixTime)(&user.CreatedAt),
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	user.IsActive = isActive == 1
	user.IsAdmin = isAdmin == 1
	return user, nil
}

// GetByEmail récupère un utilisateur par son email
func (r *SQLiteUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := &entity.User{}
	var isActive, isAdmin int

	query := `SELECT id,email,password,is_active,is_admin,created_at FROM users WHERE email=?`
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Password,
		&isActive, &isAdmin, (*unixTime)(&user.CreatedAt),
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	user.IsActive = isActive == 1
	user.IsAdmin = isAdmin == 1
	return user, nil
}

// Update met à jour un utilisateur
func (r *SQLiteUserRepository) Update(ctx context.Context, user *entity.User) error {
	query := `UPDATE users SET email=?,password=?,is_active=?,is_admin=? WHERE id=?`
	result, err := r.db.ExecContext(ctx, query,
		user.Email, user.Password,
		boolToInt(user.IsActive), boolToInt(user.IsAdmin),
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// Exists vérifie si un utilisateur existe
func (r *SQLiteUserRepository) Exists(ctx context.Context, email string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE email=?", email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check existence: %w", err)
	}
	return count > 0, nil
}

// Helpers

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

type unixTime time.Time

func (ut *unixTime) Scan(value interface{}) error {
	if v, ok := value.(int64); ok {
		*ut = unixTime(time.Unix(v, 0))
		return nil
	}
	return fmt.Errorf("cannot scan %v into unixTime", value)
}
