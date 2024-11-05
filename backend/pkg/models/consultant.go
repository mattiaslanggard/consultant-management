package models

type Consultant struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	HoursAvailable int     `json:"hours_available"`
	Skillset       string  `json:"skillset"`
	OfficeID       int     `json:"office_id"`
	Office         *Office `json:"office,omitempty"`
	Tasks          []Task  `json:"tasks"`
}
