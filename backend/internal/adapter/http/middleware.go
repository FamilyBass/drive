package http

import (
	"context"
	"net/http"

	"github.com/familybass/drive/internal/domain/service"
)

// AuthMiddleware vérifie le token JWT
func AuthMiddleware(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}

			token, err := ExtractTokenFromHeader(authHeader)
			if err != nil {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			claims, err := authService.VerifyToken(token)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			// Ajouter les claims au contexte
			ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
			ctx = context.WithValue(ctx, "is_admin", claims.IsAdmin)
			ctx = context.WithValue(ctx, "email", claims.Email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AdminMiddleware vérifie que l'utilisateur est admin
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin := r.Context().Value("is_admin").(bool)
		if !isAdmin {
			http.Error(w, "admin only", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
