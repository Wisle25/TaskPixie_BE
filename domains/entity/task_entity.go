package entity

// Task represents a task in the system.
type Task struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ProjectId   string `json:"projectId,omitempty"` // Optional
	UserId      string `json:"userId"`
	Completed   bool   `json:"completed"`
}
