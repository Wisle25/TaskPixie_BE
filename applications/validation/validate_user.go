package validation

import "github.com/wisle25/task-pixie/domains/entity"

// ValidateUser interface defines methods for validating user-related payloads.
type ValidateUser interface {
	ValidateRegisterPayload(payload *entity.RegisterUserPayload)
	ValidateLoginPayload(payload *entity.LoginUserPayload)
	ValidateUpdatePayload(payload *entity.UpdateUserPayload)
}
