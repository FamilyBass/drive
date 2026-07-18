package entity

import "time"

// User représente un utilisateur du système
type User struct {
	ID        string
	Email     string
	Password  string // hash bcrypt
	IsActive  bool
	IsAdmin   bool
	CreatedAt time.Time
}

// NewUser crée un nouvel utilisateur
func NewUser(id, email, password string, isActive, isAdmin bool, createdAt time.Time) *User {
	return &User{
		ID:        id,
		Email:     email,
		Password:  password,
		IsActive:  isActive,
		IsAdmin:   isAdmin,
		CreatedAt: createdAt,
	}
}
