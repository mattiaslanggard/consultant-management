package handlers

import (
	"consultant-management/backend/internal/db"
	"consultant-management/backend/internal/middleware"
	"consultant-management/backend/pkg/models"
	"consultant-management/backend/pkg/viewmodels"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

// ListTasks handler
func ListTasks(w http.ResponseWriter, r *http.Request) {
	conn := db.GetDB()
	rows, err := conn.Query(`
        SELECT id, consultant_id, customer_name, task_description, assigned_hours, deadline, status
        FROM tasks
        ORDER BY customer_name
    `)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tasksByCustomer := make(map[string][]models.Task)
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.ConsultantID, &task.CustomerName, &task.TaskDescription, &task.AssignedHours, &task.Deadline, &task.Status); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tasksByCustomer[task.CustomerName] = append(tasksByCustomer[task.CustomerName], task)
	}

	customerTasks := []viewmodels.CustomerTasks{}
	for customer, tasks := range tasksByCustomer {
		customerTasks = append(customerTasks, viewmodels.CustomerTasks{
			CustomerName: customer,
			Tasks:        tasks,
		})
	}

	tmpl, err := template.ParseFiles(
		"/home/mattias/consultant-management/frontend/templates/base.html",
		"/home/mattias/consultant-management/frontend/templates/tasks.html",
	)
	if err != nil {
		log.Printf("Error parsing templates: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	isAuthenticatedKey := middleware.GetIsAuthenticated(r)

	err = tmpl.ExecuteTemplate(w, "base", viewmodels.TasksData{
		Title:              "Tasks",
		CustomerTasks:      customerTasks,
		IsAuthenticatedKey: isAuthenticatedKey,
	})
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// AssignTask handler
func AssignTask(w http.ResponseWriter, r *http.Request) {
	conn := db.GetDB()
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
	_, err := conn.Exec(`
        INSERT INTO tasks (consultant_id, customer_name, task_description, assigned_hours, deadline, status)
        VALUES ($1, $2, $3, $4, $5, $6)`,
		task.ConsultantID, task.CustomerName, task.TaskDescription, task.AssignedHours, task.Deadline, task.Status,
	)

	if err != nil {
		http.Error(w, "Failed to assign task", http.StatusInternalServerError)
		return
	}

	ListTasks(w, r)
}

// DeleteTask handler
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	conn := db.GetDB()
	_, err = conn.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	ListTasks(w, r)
}
