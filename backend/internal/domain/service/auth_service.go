package service

import (
	"context"
	"fmt"
	"time"

	"github.com/familybass/drive/internal/domain/entity"
	"github.com/familybass/drive/internal/domain/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService gère l'authentification et les tokens JWT
type AuthService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

// NewAuthService crée un nouveau service d'authentification
func NewAuthService(userRepo repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// TokenClaims représente les claims du token JWT
type TokenClaims struct {
	UserID  string
	IsAdmin bool
	Email   string
}

// HashPassword hash un mot de passe
func (s *AuthService) HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashed), nil
}

// ComparePassword compare un mot de passe avec son hash
func (s *AuthService) ComparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// Register crée un nouvel utilisateur
func (s *AuthService) Register(ctx context.Context, email, password string) (*entity.User, error) {
	// Vérifier que l'utilisateur n'existe pas
	exists, err := s.userRepo.Exists(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("user already exists")
	}

	// Hash le mot de passe
	hashedPw, err := s.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Créer l'utilisateur
	user := &entity.User{
		ID:        generateUUID(),
		Email:     email,
		Password:  hashedPw,
		IsActive:  false, // Nécessite validation par admin
		IsAdmin:   false,
		CreatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Login authentifie un utilisateur et retourne un token JWT
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	// Récupérer l'utilisateur
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	// Vérifier l'activation
	if !user.IsActive {
		return "", fmt.Errorf("account not validated")
	}

	// Vérifier le mot de passe
	if err := s.ComparePassword(user.Password, password); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	// Créer le token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"admin": user.IsAdmin,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenStr, nil
}

// VerifyToken valide un token JWT et retourne ses claims
func (s *AuthService) VerifyToken(tokenStr string) (*TokenClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	userID, _ := claims["sub"].(string)
	isAdmin, _ := claims["admin"].(bool)
	email, _ := claims["email"].(string)

	return &TokenClaims{
		UserID:  userID,
		IsAdmin: isAdmin,
		Email:   email,
	}, nil
}

// ValidateUser valide un utilisateur (admin only)
func (s *AuthService) ValidateUser(ctx context.Context, userID string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	user.IsActive = true
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// EnsureAdmin crée ou active un utilisateur admin
func (s *AuthService) EnsureAdmin(ctx context.Context, email, password string) error {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		// Créer l'admin
		hashedPw, err := s.HashPassword(password)
		if err != nil {
			return err
		}

		admin := &entity.User{
			ID:        generateUUID(),
			Email:     email,
			Password:  hashedPw,
			IsActive:  true,
			IsAdmin:   true,
			CreatedAt: time.Now(),
		}

		return s.userRepo.Create(ctx, admin)
	}

	// Activation si existe
	user.IsActive = true
	user.IsAdmin = true
	return s.userRepo.Update(ctx, user)
}

// generateUUID génère un UUID v4 simple
func generateUUID() string {
	return "uuid_" + fmt.Sprintf("%d", time.Now().UnixNano())
	// À remplacer par une vraie impl si nécessaire
}
