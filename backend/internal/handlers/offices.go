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

// RenderOfficesPage handler
func RenderOfficesPage(w http.ResponseWriter, r *http.Request) {
	conn := db.GetDB()
	rows, err := conn.Query("SELECT id, name, address, zip_code, country FROM offices")
	if err != nil {
		utils.HandleError(w, err, "Error fetching offices", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	offices := []models.Office{}
	for rows.Next() {
		var o models.Office
		if err := rows.Scan(&o.ID, &o.Name, &o.Address, &o.ZipCode, &o.Country); err != nil {
			utils.HandleError(w, err, "Error scanning office", http.StatusInternalServerError)
			return
		}
		offices = append(offices, o)
	}

	tmpl, err := template.ParseFiles(
		"/home/mattias/consultant-management/frontend/templates/base.html",
		"/home/mattias/consultant-management/frontend/templates/offices.html",
		"/home/mattias/consultant-management/frontend/templates/office_list.html",
	)
	if err != nil {
		utils.HandleError(w, err, "Error parsing templates", http.StatusInternalServerError)
		return
	}

	isAuthenticatedKey := middleware.GetIsAuthenticated(r)

	err = tmpl.ExecuteTemplate(w, "base", viewmodels.OfficeData{
		Title:              "Offices",
		Offices:            offices,
		IsAuthenticatedKey: isAuthenticatedKey,
	})
	if err != nil {
		utils.HandleError(w, err, "Error executing template", http.StatusInternalServerError)
	}
}

// RenderOfficeList handler
func RenderOfficeList(w http.ResponseWriter, r *http.Request) {
	conn := db.GetDB()
	rows, err := conn.Query("SELECT id, name, address, zip_code, country FROM offices")
	if err != nil {
		utils.HandleError(w, err, "Error fetching offices", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	offices := []models.Office{}
	for rows.Next() {
		var o models.Office
		if err := rows.Scan(&o.ID, &o.Name, &o.Address, &o.ZipCode, &o.Country); err != nil {
			utils.HandleError(w, err, "Error scanning office", http.StatusInternalServerError)
			return
		}
		offices = append(offices, o)
	}

	tmpl, err := template.ParseFiles("/home/mattias/consultant-management/frontend/templates/office_list.html")
	if err != nil {
		utils.HandleError(w, err, "Error parsing template", http.StatusInternalServerError)
		return
	}

	data := viewmodels.OfficeData{
		Offices: offices,
	}

	err = tmpl.ExecuteTemplate(w, "office_list", data)
	if err != nil {
		utils.HandleError(w, err, "Error executing template", http.StatusInternalServerError)
	}
}

// AddOffice handler
func AddOffice(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		utils.HandleError(w, err, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	address := r.FormValue("address")
	zipCode := r.FormValue("zip_code")
	country := r.FormValue("country")

	conn := db.GetDB()
	_, err = conn.Exec("INSERT INTO offices (name, address, zip_code, country) VALUES ($1, $2, $3, $4)", name, address, zipCode, country)
	if err != nil {
		utils.HandleError(w, err, "Failed to add office", http.StatusInternalServerError)
		return
	}

	RenderOfficeList(w, r)
}

// EditOfficeForm handler
func EditOfficeForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.HandleError(w, err, "Invalid office ID", http.StatusBadRequest)
		return
	}

	conn := db.GetDB()
	row := conn.QueryRow("SELECT id, name, address, zip_code, country FROM offices WHERE id = $1", id)

	var o models.Office
	if err := row.Scan(&o.ID, &o.Name, &o.Address, &o.ZipCode, &o.Country); err != nil {
		utils.HandleError(w, err, "Error scanning office", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("/home/mattias/consultant-management/frontend/templates/edit_office_form.html")
	if err != nil {
		utils.HandleError(w, err, "Error parsing template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, o)
	if err != nil {
		utils.HandleError(w, err, "Error executing template", http.StatusInternalServerError)
	}
}

// EditOffice handler
func EditOffice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.HandleError(w, err, "Invalid office ID", http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		utils.HandleError(w, err, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	address := r.FormValue("address")
	zipCode := r.FormValue("zip_code")
	country := r.FormValue("country")

	conn := db.GetDB()
	_, err = conn.Exec("UPDATE offices SET name = $1, address = $2, zip_code = $3, country = $4 WHERE id = $5", name, address, zipCode, country, id)
	if err != nil {
		utils.HandleError(w, err, "Failed to update office", http.StatusInternalServerError)
		return
	}

	RenderOfficeList(w, r)
}

// DeleteOffice handler
func DeleteOffice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.HandleError(w, err, "Invalid office ID", http.StatusBadRequest)
		return
	}

	conn := db.GetDB()
	_, err = conn.Exec("DELETE FROM offices WHERE id = $1", id)
	if err != nil {
		utils.HandleError(w, err, "Failed to delete office", http.StatusInternalServerError)
		return
	}

	RenderOfficeList(w, r)
}
