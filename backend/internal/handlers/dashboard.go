package handlers

import (
	"consultant-management/backend/internal/db"
	"consultant-management/backend/internal/middleware"
	"consultant-management/backend/internal/utils"
	"consultant-management/backend/pkg/viewmodels"
	"html/template"
	"net/http"
)

// RenderDashboardPage handler
func RenderDashboardPage(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := middleware.GetIsAuthenticated(r)

	var tmpl *template.Template
	var err error
	if isAuthenticated {
		tmpl, err = template.ParseFiles(
			"/home/mattias/consultant-management/frontend/templates/base.html",
			"/home/mattias/consultant-management/frontend/templates/dashboard_authenticated.html",
		)
	} else {
		tmpl, err = template.ParseFiles(
			"/home/mattias/consultant-management/frontend/templates/base.html",
			"/home/mattias/consultant-management/frontend/templates/dashboard_unauthenticated.html",
		)
	}
	if err != nil {
		utils.HandleError(w, err, "Error parsing templates", http.StatusInternalServerError)
		return
	}

	conn := db.GetDB()

	// Fetch the number of consultants
	var numConsultants int
	err = conn.QueryRow("SELECT COUNT(*) FROM consultants").Scan(&numConsultants)
	if err != nil {
		utils.HandleError(w, err, "Failed to fetch number of consultants", http.StatusInternalServerError)
		return
	}

	// Fetch the number of tasks
	var numTasks int
	err = conn.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&numTasks)
	if err != nil {
		utils.HandleError(w, err, "Failed to fetch number of tasks", http.StatusInternalServerError)
		return
	}

	data := viewmodels.DashboardData{
		Title:              "Dashboard",
		NumTasks:           numTasks,
		NumConsultants:     numConsultants,
		IsAuthenticatedKey: isAuthenticated,
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		utils.HandleError(w, err, "Error executing template", http.StatusInternalServerError)
	}
}
