package entity

// ProjectPayload represents the payload for creating or updating a project.
type ProjectPayload struct {
	Title     string   `json:"title"`
	Detail    string   `json:"detail"`
	Priority  string   `json:"priority"`
	Status    string   `json:"status"`
	MembersId []string `json:"members"` // User IDs
}

// PreviewProject represents a brief overview of a project.
type PreviewProject struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

// Project represents the detailed view of a project in the system.
type Project struct {
	Id              string   `json:"id"`
	Title           string   `json:"title"`
	Detail          string   `json:"detail"`
	Priority        string   `json:"priority"`
	Status          string   `json:"status"`
	MembersUsername []string `json:"members"` // Usernames or User IDs as needed
	CreatedAt       string   `json:"createdAt"`
	UpdatedAt       string   `json:"updatedAt"`
}
