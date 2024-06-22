package validation

import (
	"github.com/wisle25/task-pixie/applications/validation"
	"github.com/wisle25/task-pixie/domains/entity"
	"github.com/wisle25/task-pixie/infrastructures/services"
)

type GoValidateProject struct {
	validation *services.Validation
}

func NewValidateProject(validation *services.Validation) validation.ValidateProject {
	return &GoValidateProject{
		validation: validation,
	}
}

func (v *GoValidateProject) ValidatePayload(payload *entity.ProjectPayload) {
	schema := map[string]string{
		"Title":     "required,min=3,max=100",
		"Detail":    "required,min=3,max=1000",
		"Priority":  "required,oneof=Low High Urgent",
		"Status":    "required,oneof='To Do' 'In Progress' Completed Canceled",
		"MembersId": "required,dive,uuid",
	}

	services.Validate(payload, schema, v.validation)
}
