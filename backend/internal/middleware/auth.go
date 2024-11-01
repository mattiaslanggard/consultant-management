package middleware

import (
	"context"
	"net/http"

	"consultant-management/backend/internal/utils"

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
				utils.HandleError(w, err, "Unauthorized", http.StatusUnauthorized)
				return
			}
			utils.HandleError(w, err, "Bad request", http.StatusBadRequest)
			return
		}

		tokenStr := c.Value
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				utils.HandleError(w, err, "Unauthorized", http.StatusUnauthorized)
				return
			}
			utils.HandleError(w, err, "Bad request", http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			utils.HandleError(w, err, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), IsAuthenticatedKey, true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// SetAuthContext sets the isAuthenticated context value without blocking access
func SetAuthContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				// Set isAuthenticated to false and continue
				ctx := context.WithValue(r.Context(), IsAuthenticatedKey, false)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			utils.HandleError(w, err, "Bad request", http.StatusBadRequest)
			return
		}

		tokenStr := c.Value
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				// Set isAuthenticated to false and continue
				ctx := context.WithValue(r.Context(), IsAuthenticatedKey, false)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			utils.HandleError(w, err, "Bad request", http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			// Set isAuthenticated to false and continue
			ctx := context.WithValue(r.Context(), IsAuthenticatedKey, false)
			next.ServeHTTP(w, r.WithContext(ctx))
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
