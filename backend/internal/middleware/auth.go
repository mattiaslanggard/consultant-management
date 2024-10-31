package middleware

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type contextKey string

const IsAuthenticatedKey contextKey = "isAuthenticated"

// AuthMiddleware to protect routes
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				ctx := context.WithValue(r.Context(), IsAuthenticatedKey, false)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := c.Value
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if !tkn.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), IsAuthenticatedKey, true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetIsAuthenticated retrieves the isAuthenticated value from the context
func GetIsAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(IsAuthenticatedKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
