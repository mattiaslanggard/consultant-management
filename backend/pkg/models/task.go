package models

import "time"

type Task struct {
	ID              int       `json:"id"`
	ConsultantID    int       `json:"consultant_id"`
	CustomerName    string    `json:"customer_name"`
	TaskDescription string    `json:"task_description"`
	AssignedHours   int       `json:"assigned_hours"`
	Deadline        time.Time `json:"deadline"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

type CustomerTasks struct {
	CustomerName string
	Tasks        []Task
}
