package validation

import "github.com/wisle25/task-pixie/domains/entity"

// ValidateProject interface defines methods for validating project-related payloads.
type ValidateProject interface {
	ValidatePayload(payload *entity.ProjectPayload)
}
