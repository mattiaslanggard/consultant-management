package handlers

import (
	"consultant-management/backend/internal/db"
	"consultant-management/backend/internal/middleware"
	"consultant-management/backend/internal/utils"
	"consultant-management/backend/pkg/viewmodels"
	"html/template"
	"net/http"
)

// RenderReportPage handler
func RenderReportPage(w http.ResponseWriter, r *http.Request) {
	conn := db.GetDB()
	rows, err := conn.Query("SELECT id, name, hours_available FROM consultants")
	if err != nil {
		utils.HandleError(w, err, "Error fetching consultants", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	consultants := []viewmodels.ConsultantReport{}
	for rows.Next() {
		var c viewmodels.ConsultantReport
		var consultantID int
		if err := rows.Scan(&consultantID, &c.Name, &c.HoursAvailable); err != nil {
			utils.HandleError(w, err, "Error scanning consultant", http.StatusInternalServerError)
			return
		}

		// Fetch tasks for each consultant
		taskRows, err := conn.Query("SELECT customer_name, task_description, assigned_hours, status, deadline FROM tasks WHERE consultant_id = $1", consultantID)
		if err != nil {
			utils.HandleError(w, err, "Error fetching tasks", http.StatusInternalServerError)
			return
		}
		defer taskRows.Close()

		tasks := []viewmodels.TaskReport{}
		for taskRows.Next() {
			var t viewmodels.TaskReport
			if err := taskRows.Scan(&t.CustomerName, &t.TaskDescription, &t.AssignedHours, &t.Status, &t.Deadline); err != nil {
				utils.HandleError(w, err, "Error scanning task", http.StatusInternalServerError)
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
		utils.HandleError(w, err, "Error parsing templates", http.StatusInternalServerError)
		return
	}

	isAuthenticatedKey := middleware.GetIsAuthenticated(r)

	err = tmpl.ExecuteTemplate(w, "base", viewmodels.ReportData{
		Title:              "Consultant Report",
		Consultants:        consultants,
		IsAuthenticatedKey: isAuthenticatedKey,
	})
	if err != nil {
		utils.HandleError(w, err, "Error executing template", http.StatusInternalServerError)
	}
}
