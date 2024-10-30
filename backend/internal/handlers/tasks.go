package handlers

import (
	"consultant-management/backend/internal/db"
	"consultant-management/backend/pkg/models"
	"net/http"

	"github.com/gorilla/schema" // Use schema to parse form data
)

var decoder = schema.NewDecoder()

// Task assignment handler
func AssignTask(w http.ResponseWriter, r *http.Request) {
	db := db.GetDB()
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Decode form values into the Task struct
	var task models.Task
	if err := decoder.Decode(&task, r.PostForm); err != nil {
		http.Error(w, "Failed to parse task data", http.StatusBadRequest)
		return
	}

	// Prepare SQL statement
	_, err := db.Exec(`
        INSERT INTO tasks (consultant_id, customer_name, task_description, assigned_hours, deadline, status)
        VALUES ($1, $2, $3, $4, $5, $6)`,
		task.ConsultantID, task.CustomerName, task.TaskDescription, task.AssignedHours, task.Deadline, task.Status,
	)

	if err != nil {
		http.Error(w, "Failed to assign task", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Task assigned successfully"))
}
