package handlers

import (
	"consultant-management/backend/internal/db"
	"consultant-management/backend/internal/middleware"
	"consultant-management/backend/internal/utils"
	"consultant-management/backend/pkg/models"
	"consultant-management/backend/pkg/viewmodels"
	"encoding/json"
	"html/template"
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
		utils.HandleError(w, err, "Error fetching tasks", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tasksByCustomer := make(map[string][]models.Task)
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.ConsultantID, &task.CustomerName, &task.TaskDescription, &task.AssignedHours, &task.Deadline, &task.Status); err != nil {
			utils.HandleError(w, err, "Error scanning task", http.StatusInternalServerError)
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
		utils.HandleError(w, err, "Error parsing templates", http.StatusInternalServerError)
		return
	}

	isAuthenticatedKey := middleware.GetIsAuthenticated(r)

	err = tmpl.ExecuteTemplate(w, "base", viewmodels.TasksData{
		Title:              "Tasks",
		CustomerTasks:      customerTasks,
		IsAuthenticatedKey: isAuthenticatedKey,
	})
	if err != nil {
		utils.HandleError(w, err, "Error executing template", http.StatusInternalServerError)
	}
}

// EditTaskForm handler
func EditTaskForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.HandleError(w, err, "Invalid task ID", http.StatusBadRequest)
		return
	}

	conn := db.GetDB()
	row := conn.QueryRow(`
        SELECT id, consultant_id, customer_name, task_description, assigned_hours, deadline, status, created_at
        FROM tasks WHERE id = $1`, id)

	var task models.Task
	err = row.Scan(&task.ID, &task.ConsultantID, &task.CustomerName, &task.TaskDescription, &task.AssignedHours, &task.Deadline, &task.Status, &task.CreatedAt)
	if err != nil {
		utils.HandleError(w, err, "Failed to fetch task", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("/home/mattias/consultant-management/frontend/templates/edit_task.html")
	if err != nil {
		utils.HandleError(w, err, "Error parsing template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, task)
	if err != nil {
		utils.HandleError(w, err, "Error executing template", http.StatusInternalServerError)
	}
}

// GetTask handler
func GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.HandleError(w, err, "Invalid task ID", http.StatusBadRequest)
		return
	}

	conn := db.GetDB()
	row := conn.QueryRow(`
        SELECT id, consultant_id, customer_name, task_description, assigned_hours, deadline, status, created_at
        FROM tasks WHERE id = $1`, id)

	var task models.Task
	err = row.Scan(&task.ID, &task.ConsultantID, &task.CustomerName, &task.TaskDescription, &task.AssignedHours, &task.Deadline, &task.Status, &task.CreatedAt)
	if err != nil {
		utils.HandleError(w, err, "Failed to fetch task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// AssignTask handler
func AssignTask(w http.ResponseWriter, r *http.Request) {
	conn := db.GetDB()
	// Parse form data
	if err := r.ParseForm(); err != nil {
		utils.HandleError(w, err, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Decode form values into the Task struct
	var task models.Task
	if err := decoder.Decode(&task, r.PostForm); err != nil {
		utils.HandleError(w, err, "Failed to decode form data", http.StatusBadRequest)
		return
	}

	// Prepare SQL statement
	_, err := conn.Exec(`
        INSERT INTO tasks (consultant_id, customer_name, task_description, assigned_hours, deadline, status)
        VALUES ($1, $2, $3, $4, $5, $6)`,
		task.ConsultantID, task.CustomerName, task.TaskDescription, task.AssignedHours, task.Deadline, task.Status,
	)

	if err != nil {
		utils.HandleError(w, err, "Failed to insert task", http.StatusInternalServerError)
		return
	}

	ListTasks(w, r)
}

// UpdateTask handler
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.HandleError(w, err, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		utils.HandleError(w, err, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Decode form values into the Task struct
	var task models.Task
	if err := decoder.Decode(&task, r.PostForm); err != nil {
		utils.HandleError(w, err, "Failed to decode form data", http.StatusBadRequest)
		return
	}

	conn := db.GetDB()
	_, err = conn.Exec(`
        UPDATE tasks SET consultant_id = $1, customer_name = $2, task_description = $3, assigned_hours = $4, deadline = $5, status = $6
        WHERE id = $7`,
		task.ConsultantID, task.CustomerName, task.TaskDescription, task.AssignedHours, task.Deadline, task.Status, id,
	)
	if err != nil {
		utils.HandleError(w, err, "Failed to update task", http.StatusInternalServerError)
		return
	}

	ListTasks(w, r)
}

// DeleteTask handler
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.HandleError(w, err, "Invalid task ID", http.StatusBadRequest)
		return
	}

	conn := db.GetDB()
	_, err = conn.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		utils.HandleError(w, err, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	ListTasks(w, r)
}
