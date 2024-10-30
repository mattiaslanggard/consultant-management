package handlers

import (
	"consultant-management/backend/internal/db"
	"html/template"
	"log"
	"net/http"
)

type DashboardData struct {
	NumConsultants int
	NumTasks       int
}

// RenderDashboardPage handler
func RenderDashboardPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"/home/mattias/consultant-management/frontend/templates/base.html",
		"/home/mattias/consultant-management/frontend/templates/dashboard.html",
	)
	if err != nil {
		log.Printf("Error parsing templates: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", map[string]interface{}{
		"Title": "Dashboard",
	})
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
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

	data := DashboardData{
		NumConsultants: numConsultants,
		NumTasks:       numTasks,
	}

	tmpl := template.Must(template.ParseFiles("/home/mattias/consultant-management/frontend/templates/dashboard.html"))
	tmpl.Execute(w, data)
}
