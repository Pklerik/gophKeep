package middleware

import (
	"net/http"

	"github.com/Pklerik/gophKeep/internal/auth"
	"github.com/Pklerik/gophKeep/internal/service"
)

func AuthMiddleware(userService *service.UserService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tokenString, err := auth.ExtractToken(authHeader)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Validate token and get user ID
			claim, err := auth.VerifyToken(tokenString)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Optionally, you can fetch the user from the database if needed
			_, err = userService.GetUser(claim.UserID)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}
