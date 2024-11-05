package models

type Office struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}
