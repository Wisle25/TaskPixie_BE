package entity

// TaskPayload represents the payload for creating or updating a task.
type TaskPayload struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Detail       string   `json:"detail"`
	Priority     string   `json:"priority"`
	Status       string   `json:"status"`
	ProjectId    string   `json:"projectId"`
	DueDate      string   `json:"dueDate"`
	AssignedToId []string `json:"assignedTo"` // User IDs
}

// PreviewTask represents a brief overview of a task.
type PreviewTask struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Status      string `json:"status"`
	Project     string `json:"project"` // Project name
}

// Task represents the detailed view of a task in the system.
type Task struct {
	ID                  string   `json:"id"`
	Title               string   `json:"title"`
	Description         string   `json:"description"`
	Detail              string   `json:"detail"`
	Priority            string   `json:"priority"`
	Status              string   `json:"status"`
	Project             string   `json:"project"`
	AssignedToUsernames []string `json:"assignedTo"` // Usernames
	DueDate             string   `json:"dueDate"`
	CreatedAt           string   `json:"createdAt"`
	UpdatedAt           string   `json:"updatedAt"`
	ProjectId           string   `json:"projectId"`
}
