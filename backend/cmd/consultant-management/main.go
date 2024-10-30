package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"consultant-management/backend/internal/db"
	"consultant-management/backend/internal/handlers"
	"consultant-management/backend/internal/middleware"

	"github.com/gorilla/mux"
)

func main() {
	// Read database credentials from environment variables
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbSSLMode := os.Getenv("DB_SSLMODE")
	// Initialize database connection
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=%s", dbUser, dbName, dbPassword, dbSSLMode)
	err := db.InitDB(connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.CloseDB()

	// Set up the router
	router := mux.NewRouter()
	router.HandleFunc("/register", handlers.Register).Methods("POST")
	router.HandleFunc("/login", handlers.Login).Methods("POST")
	router.HandleFunc("/login", handlers.LoginPage).Methods("GET")
	router.HandleFunc("/register", handlers.RegisterPage).Methods("GET")
	router.HandleFunc("/dashboard", handlers.RenderDashboardPage).Methods("GET")

	// Apply authentication middleware to protected routes
	protected := router.PathPrefix("/").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/consultants", handlers.RenderConsultantsPage).Methods("GET")
	protected.HandleFunc("/report", handlers.RenderReportPage).Methods("GET")
	protected.HandleFunc("/consultants", handlers.AddConsultant).Methods("POST")
	protected.HandleFunc("/consultants/{id}", handlers.EditConsultant).Methods("POST")
	protected.HandleFunc("/consultants/{id}/edit", handlers.EditConsultantForm).Methods("GET")
	protected.HandleFunc("/consultants/{id}/delete", handlers.DeleteConsultant).Methods("POST")
	protected.HandleFunc("/assign_task", handlers.AssignTask).Methods("POST")

	// Serve static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/static"))))

	// Serve the index.html template for the root URL
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}).Methods("GET")

	// Start the server
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
