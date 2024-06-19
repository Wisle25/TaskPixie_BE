//go:build wireinject
// +build wireinject

package container

import (
	"database/sql"
	"github.com/google/wire"
	"github.com/wisle25/be-template/applications/cache"
	"github.com/wisle25/be-template/applications/file_statics"
	"github.com/wisle25/be-template/applications/generator"
	"github.com/wisle25/be-template/applications/use_case"
	"github.com/wisle25/be-template/commons"
	"github.com/wisle25/be-template/infrastructures/repository"
	"github.com/wisle25/be-template/infrastructures/security"
	"github.com/wisle25/be-template/infrastructures/services"
	"github.com/wisle25/be-template/infrastructures/validation"
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
