package handlers

import (
	"consultant-management/backend/internal/db"
	"consultant-management/backend/pkg/models"
	"html/template"
	"log"
	"net/http"
)

// RenderReportPage handler
func RenderReportPage(w http.ResponseWriter, r *http.Request) {
	conn := db.GetDB()
	rows, err := conn.Query("SELECT id, name, hours_available FROM consultants")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	consultants := []models.ConsultantReport{}
	for rows.Next() {
		var c models.ConsultantReport
		var consultantID int
		if err := rows.Scan(&consultantID, &c.Name, &c.HoursAvailable); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Fetch tasks for each consultant
		taskRows, err := conn.Query("SELECT customer_name, task_description, assigned_hours, status, deadline FROM tasks WHERE consultant_id = $1", consultantID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer taskRows.Close()

		tasks := []models.TaskReport{}
		for taskRows.Next() {
			var t models.TaskReport
			if err := taskRows.Scan(&t.CustomerName, &t.TaskDescription, &t.AssignedHours, &t.Status, &t.Deadline); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			tasks = append(tasks, t)
		}
		c.Tasks = tasks

		consultants = append(consultants, c)
	}

	tmpl, err := template.ParseFiles(
		"/home/mattias/consultant-management/frontend/templates/base.html",
		"/home/mattias/consultant-management/frontend/templates/report.html",
	)
	if err != nil {
		log.Printf("Error parsing templates: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", map[string]interface{}{
		"Title":       "Consultant Report",
		"Consultants": consultants,
	})
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
