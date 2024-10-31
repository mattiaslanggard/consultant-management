package models

type ReportData struct {
	Consultants []ConsultantReport
}

type ConsultantReport struct {
	Name           string
	HoursAvailable int
	Tasks          []TaskReport
}

type TaskReport struct {
	CustomerName    string
	TaskDescription string
	AssignedHours   int
	Status          string
	Deadline        string
}
