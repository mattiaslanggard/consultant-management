package handlers

import (
	"consultant-management/backend/internal/db"
	"consultant-management/backend/internal/middleware"
	"consultant-management/backend/pkg/viewmodels"
	"html/template"
	"log"
	"net/http"
)

// RenderDashboardPage handler
func RenderDashboardPage(w http.ResponseWriter, r *http.Request) {
	conn := db.GetDB()

	// Fetch the number of consultants
	var numConsultants int
	err := conn.QueryRow("SELECT COUNT(*) FROM consultants").Scan(&numConsultants)
	if err != nil {
		http.Error(w, "Failed to fetch number of consultants", http.StatusInternalServerError)
		return
	}

	// Fetch the number of tasks
	var numTasks int
	err = conn.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&numTasks)
	if err != nil {
		http.Error(w, "Failed to fetch number of tasks", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(
		"/home/mattias/consultant-management/frontend/templates/base.html",
		"/home/mattias/consultant-management/frontend/templates/dashboard.html",
	)
	if err != nil {
		log.Printf("Error parsing templates: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	isAuthenticatedKey := middleware.GetIsAuthenticated(r)

	err = tmpl.ExecuteTemplate(w, "base", viewmodels.DashboardData{
		Title:              "Dashboard",
		NumTasks:           numTasks,
		NumConsultants:     numConsultants,
		IsAuthenticatedKey: isAuthenticatedKey,
	})
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
