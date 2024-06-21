//go:build wireinject
// +build wireinject

package container

import (
	"database/sql"
	"github.com/google/wire"
	"github.com/wisle25/task-pixie/applications/cache"
	"github.com/wisle25/task-pixie/applications/file_statics"
	"github.com/wisle25/task-pixie/applications/generator"
	"github.com/wisle25/task-pixie/applications/use_case"
	"github.com/wisle25/task-pixie/commons"
	"github.com/wisle25/task-pixie/infrastructures/repository"
	"github.com/wisle25/task-pixie/infrastructures/security"
	"github.com/wisle25/task-pixie/infrastructures/services"
	"github.com/wisle25/task-pixie/infrastructures/validation"
)

// Dependency Injection for User Use Case
func NewUserContainer(
	config *commons.Config,
	db *sql.DB,
	cache cache.Cache,
	idGenerator generator.IdGenerator,
	fileProcessing file_statics.FileProcessing,
	fileUpload file_statics.FileUpload,
	validator *services.Validation,
) *use_case.UserUseCase {
	wire.Build(
		repository.NewUserRepositoryPG,
		security.NewArgon2,
		validation.NewValidateUser,
		security.NewJwtToken,
		use_case.NewUserUseCase,
	)

	return nil
}

// Dependency Injection for Project Use Case
func NewProjectContainer(idGenerator generator.IdGenerator, db *sql.DB) *use_case.ProjectUseCase {
	wire.Build(
		repository.NewProjectRepositoryPG,
		use_case.NewProjectUseCase,
	)
	return nil
}

// Dependency Injection for Task Use Case
func NewTaskContainer(
	idgenerator generator.IdGenerator,
	db *sql.DB,
) *use_case.TaskUseCase {
	wire.Build(
		repository.NewTaskRepositoryPG,
		use_case.NewTaskUseCase,
	)

	return nil
}
