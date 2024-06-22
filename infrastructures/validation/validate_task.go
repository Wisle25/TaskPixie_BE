package validation

import (
	"github.com/wisle25/task-pixie/applications/validation"
	"github.com/wisle25/task-pixie/domains/entity"
	"github.com/wisle25/task-pixie/infrastructures/services"
)

type GoValidateTask struct {
	validation *services.Validation
}

func NewValidateTask(validation *services.Validation) validation.ValidateTask {
	return &GoValidateTask{
		validation: validation,
	}
}

func (v *GoValidateTask) ValidatePayload(payload *entity.TaskPayload) {
	schema := map[string]string{
		"Title":        "required,min=3,max=100",
		"Description":  "required,min=3,max=1000",
		"Detail":       "omitempty,min=3,max=1000",
		"Priority":     "required,oneof=Low High Urgent",
		"Status":       "required,oneof='To Do' 'In Progress' Completed Canceled",
		"ProjectId":    "omitempty,uuid",
		"DueDate":      "required",
		"AssignedToId": "omitempty,dive,uuid",
	}

	services.Validate(payload, schema, v.validation)
}
