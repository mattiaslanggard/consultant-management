package models

type Consultant struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	HoursAvailable int    `json:"hours_available"`
	Skillset       string `json:"skillset"`
	Tasks          []Task `json:"tasks"`
}
