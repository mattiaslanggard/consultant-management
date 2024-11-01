package handlers

import (
	"consultant-management/backend/internal/db"
	"consultant-management/backend/internal/logger"
	"consultant-management/backend/internal/utils"
	"consultant-management/backend/pkg/models"
	"html/template"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// LoginPage handler
func LoginPage(w http.ResponseWriter, r *http.Request) {
	renderLoginPage(w, "")
}

// renderLoginPage renders the login page with an optional error message
func renderLoginPage(w http.ResponseWriter, errorMessage string) {
	tmpl, err := template.ParseFiles(
		"/home/mattias/consultant-management/frontend/templates/base.html",
		"/home/mattias/consultant-management/frontend/templates/login.html",
	)
	if err != nil {
		utils.HandleError(w, err, "Error parsing templates", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":        "Login",
		"ErrorMessage": errorMessage,
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		utils.HandleError(w, err, "Error executing template", http.StatusInternalServerError)
	}
}

// RegisterPage handler
func RegisterPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"/home/mattias/consultant-management/frontend/templates/base.html",
		"/home/mattias/consultant-management/frontend/templates/register.html",
	)
	if err != nil {
		logger.ErrorLogger.Printf("Error parsing templates: %v", err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", map[string]interface{}{
		"Title": "Register",
	})
	if err != nil {
		utils.HandleError(w, err, "Error executing template", http.StatusInternalServerError)
	}
}

var jwtKey = []byte("your_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Register handler
func Register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		utils.HandleError(w, err, "Invalid request payload", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		utils.HandleError(w, nil, "Username and password are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		utils.HandleError(w, err, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	conn := db.GetDB()
	_, err = conn.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, string(hashedPassword))
	if err != nil {
		utils.HandleError(w, err, "Failed to register user", http.StatusInternalServerError)
		return
	}

	logger.InfoLogger.Printf("User registered: %s", username)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Login handler
func Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		utils.HandleError(w, err, "Error parsing form", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		renderLoginPage(w, "Username and password are required")
		return
	}

	conn := db.GetDB()
	row := conn.QueryRow("SELECT id, password FROM users WHERE username = $1", username)
	var storedUser models.User
	err = row.Scan(&storedUser.ID, &storedUser.Password)
	if err != nil {
		renderLoginPage(w, "Invalid username or password")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(password))
	if err != nil {
		renderLoginPage(w, "Invalid username or password")
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		utils.HandleError(w, err, "Error generating token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	logger.InfoLogger.Printf("User logged in: %s", username)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// Logout handler
func Logout(w http.ResponseWriter, r *http.Request) {
	// Clear the authentication cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	logger.InfoLogger.Println("User logged out")
	// Redirect to the login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Middleware to protect routes
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

		logger.InfoLogger.Printf("User authenticated: %s", claims.Username)
		next.ServeHTTP(w, r)
	})
}
