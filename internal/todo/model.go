package todo

import (
	"time"
)

type Priority string

const (
	PriorityHigh   Priority = "high"
	PriorityMedium Priority = "medium"
	PriorityLow    Priority = "low"
)

type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "inprogress"
	StatusDone       Status = "done"
)

type Task struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Priority    Priority `json:"priority"`
	Status      Status   `json:"status"`
	DueDate     string   `json:"dueDate"`
	Tags        []string `json:"tags"`
}

type TaskStore struct {
	Tasks []Task `json:"tasks"`
}

func NewTask(title, description string, priority Priority, dueDate string, tags []string) *Task {
	return &Task{
		ID:          generateID(),
		Title:       title,
		Description: description,
		Priority:    priority,
		Status:      StatusPending,
		DueDate:     dueDate,
		Tags:        tags,
	}
}

func generateID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[i%len(letters)]
	}
	return string(b)
}
