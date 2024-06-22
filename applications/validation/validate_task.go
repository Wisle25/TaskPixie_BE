package validation

import "github.com/wisle25/task-pixie/domains/entity"

// ValidateTask interface defines methods for validating task-related payloads.
type ValidateTask interface {
	ValidatePayload(payload *entity.TaskPayload)
}
