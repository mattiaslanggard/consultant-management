package viewmodels

import "consultant-management/backend/pkg/models"

type DashboardData struct {
	Title              string
	NumConsultants     int
	NumTasks           int
	IsAuthenticatedKey bool
}

type ConsultantsData struct {
	Title              string
	Consultants        []models.Consultant
	IsAuthenticatedKey bool
}

type TaskReport struct {
	CustomerName    string
	TaskDescription string
	AssignedHours   int
	Status          string
	Deadline        string
}

type CustomerTasks struct {
	CustomerName string
	Tasks        []models.Task
}

type TasksData struct {
	Title              string
	CustomerTasks      []CustomerTasks
	IsAuthenticatedKey bool
}

type ConsultantReport struct {
	Name           string
	HoursAvailable int
	Tasks          []TaskReport
}

type ReportData struct {
	Title              string
	Consultants        []ConsultantReport
	IsAuthenticatedKey bool
}
