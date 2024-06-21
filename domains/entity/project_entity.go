package entity

// ProjectPayload represents the payload for creating or updating a project.
type ProjectPayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Project represents a project in the system.
type Project struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	OwnerId       string `json:"ownerId"`
	OwnerUsername string `json:"ownerUsername"`
	CreatedAt     string `json:"createdAt"`
}
