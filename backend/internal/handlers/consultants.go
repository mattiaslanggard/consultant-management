package handlers

import (
	"consultant-management/backend/internal/db"
	"consultant-management/backend/internal/middleware"
	"consultant-management/backend/internal/utils"
	"consultant-management/backend/pkg/models"
	"consultant-management/backend/pkg/viewmodels"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// RenderConsultantsPage handler
func RenderConsultantsPage(w http.ResponseWriter, r *http.Request) {
	conn := db.GetDB()
	rows, err := conn.Query("SELECT id, name, hours_available, skillset FROM consultants")
	if err != nil {
		utils.HandleError(w, err, "Error fetching consultants", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	consultants := []models.Consultant{}
	for rows.Next() {
		var c models.Consultant
		if err := rows.Scan(&c.ID, &c.Name, &c.HoursAvailable, &c.Skillset); err != nil {
			utils.HandleError(w, err, "Error scanning consultant", http.StatusInternalServerError)
			return
		}
		consultants = append(consultants, c)
	}

	tmpl, err := template.ParseFiles(
		"/home/mattias/consultant-management/frontend/templates/base.html",
		"/home/mattias/consultant-management/frontend/templates/consultants.html",
	)
	if err != nil {
		utils.HandleError(w, err, "Error parsing templates", http.StatusInternalServerError)
		return
	}

	isAuthenticatedKey := middleware.GetIsAuthenticated(r)

	err = tmpl.ExecuteTemplate(w, "base", viewmodels.ConsultantsData{
		Title:              "Consultants",
		Consultants:        consultants,
		IsAuthenticatedKey: isAuthenticatedKey,
	})
	if err != nil {
		utils.HandleError(w, err, "Error executing template", http.StatusInternalServerError)
	}
}

// RenderConsultantList handler
func RenderConsultantList(w http.ResponseWriter, r *http.Request) {
	conn := db.GetDB()
	rows, err := conn.Query("SELECT id, name, hours_available, skillset FROM consultants")
	if err != nil {
		utils.HandleError(w, err, "Error fetching consultants", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	consultants := []models.Consultant{}
	for rows.Next() {
		var c models.Consultant
		if err := rows.Scan(&c.ID, &c.Name, &c.HoursAvailable, &c.Skillset); err != nil {
			utils.HandleError(w, err, "Error scanning consultant", http.StatusInternalServerError)
			return
		}
		consultants = append(consultants, c)
	}

	tmpl, err := template.ParseFiles("/home/mattias/consultant-management/frontend/templates/consultant_list.html")
	if err != nil {
		utils.HandleError(w, err, "Error parsing template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, consultants)
	if err != nil {
		utils.HandleError(w, err, "Error executing template", http.StatusInternalServerError)
	}
}

// AddConsultant handler
func AddConsultant(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		utils.HandleError(w, err, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	hoursAvailable, err := strconv.Atoi(r.FormValue("hours_available"))
	if err != nil {
		utils.HandleError(w, err, "Invalid hours available", http.StatusBadRequest)
		return
	}
	skillset := r.FormValue("skillset")

	conn := db.GetDB()
	_, err = conn.Exec("INSERT INTO consultants (name, hours_available, skillset) VALUES ($1, $2, $3)", name, hoursAvailable, skillset)
	if err != nil {
		utils.HandleError(w, err, "Failed to add consultant", http.StatusInternalServerError)
		return
	}

	RenderConsultantList(w, r)
}

// EditConsultantForm handler
func EditConsultantForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.HandleError(w, err, "Invalid consultant ID", http.StatusBadRequest)
		return
	}

	conn := db.GetDB()
	row := conn.QueryRow("SELECT id, name, hours_available, skillset FROM consultants WHERE id = $1", id)

	var c models.Consultant
	if err := row.Scan(&c.ID, &c.Name, &c.HoursAvailable, &c.Skillset); err != nil {
		utils.HandleError(w, err, "Error scanning consultant", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("/home/mattias/consultant-management/frontend/templates/edit_form.html")
	if err != nil {
		utils.HandleError(w, err, "Error parsing template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, c)
	if err != nil {
		utils.HandleError(w, err, "Error executing template", http.StatusInternalServerError)
	}
}

// EditConsultant handler
func EditConsultant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.HandleError(w, err, "Invalid consultant ID", http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		utils.HandleError(w, err, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	hoursAvailable, err := strconv.Atoi(r.FormValue("hours_available"))
	if err != nil {
		utils.HandleError(w, err, "Invalid hours available", http.StatusBadRequest)
		return
	}
	skillset := r.FormValue("skillset")

	conn := db.GetDB()
	_, err = conn.Exec("UPDATE consultants SET name = $1, hours_available = $2, skillset = $3 WHERE id = $4", name, hoursAvailable, skillset, id)
	if err != nil {
		utils.HandleError(w, err, "Failed to update consultant", http.StatusInternalServerError)
		return
	}

	RenderConsultantList(w, r)
}

// DeleteConsultant handler
func DeleteConsultant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.HandleError(w, err, "Invalid consultant ID", http.StatusBadRequest)
		return
	}

	conn := db.GetDB()
	_, err = conn.Exec("DELETE FROM consultants WHERE id = $1", id)
	if err != nil {
		utils.HandleError(w, err, "Failed to delete consultant", http.StatusInternalServerError)
		return
	}

	RenderConsultantList(w, r)
}
