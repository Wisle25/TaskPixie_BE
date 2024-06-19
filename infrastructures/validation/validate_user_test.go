package validation_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/wisle25/be-template/domains/entity"
	"github.com/wisle25/be-template/infrastructures/services"
	"testing"

	"github.com/wisle25/be-template/infrastructures/validation"
)

func TestValidateUser(t *testing.T) {
	validator := services.NewValidation()
	validateUser := validation.NewValidateUser(validator)

	t.Run("Register and Update Validation", func(t *testing.T) {
		t.Run("Username Validation", func(t *testing.T) {
			t.Run("Should return error when reaches the max limit", func(t *testing.T) {
				// Arrange
				registerPayload := &entity.RegisterUserPayload{
					Username:        "awekrhnawklejrklewhjrtleksjrtklerwntklrsjgklrdjgtyklrdjgtyklrdjgkltydejrt",
					Password:        "Helkdaoskd",
					ConfirmPassword: "Helkdaoskd",
					Email:           "aiwerkajwe@gmail.com",
				}

				updatePayload := &entity.UpdateUserPayload{
					Username:        "awekrhnawklejrklewhjrtleksjrtklerwntklrsjgklrdjgtyklrdjgtyklrdjgkltydejrt",
					Password:        "Helkdaoskd",
					ConfirmPassword: "Helkdaoskd",
					Email:           "aiwerkajwe@gmail.com",
				}

				// Action and Assert
				assert.Panics(t, func() {
					validateUser.ValidateRegisterPayload(registerPayload)
				})
				assert.Panics(t, func() {
					validateUser.ValidateUpdatePayload(updatePayload)
				})
			})

			t.Run("Should return error when username is empty", func(t *testing.T) {
				// Arrange
				registerPayload := &entity.RegisterUserPayload{
					Username:        "",
					Password:        "Helkdaoskd",
					ConfirmPassword: "Helkdaoskd",
					Email:           "aiwerkajwe@gmail.com",
				}

				updatePayload := &entity.UpdateUserPayload{
					Username:        "",
					Password:        "Helkdaoskd",
					ConfirmPassword: "Helkdaoskd",
					Email:           "aiwerkajwe@gmail.com",
				}

				// Action and Assert
				assert.Panics(t, func() {
					validateUser.ValidateRegisterPayload(registerPayload)
				})
				assert.Panics(t, func() {
					validateUser.ValidateUpdatePayload(updatePayload)
				})
			})

			t.Run("Should return error when username contains illegal chars", func(t *testing.T) {
				// Arrange
				registerPayload := &entity.RegisterUserPayload{
					Username:        "AS awdawd",
					Password:        "Helkdaoskd",
					ConfirmPassword: "Helkdaoskd",
					Email:           "aiwerkajwe@gmail.com",
				}

				updatePayload := &entity.UpdateUserPayload{
					Username:        "AS awdawd",
					Password:        "Helkdaoskd",
					ConfirmPassword: "Helkdaoskd",
					Email:           "aiwerkajwe@gmail.com",
				}

				// Action and Assert
				assert.Panics(t, func() {
					validateUser.ValidateRegisterPayload(registerPayload)
				})
				assert.Panics(t, func() {
					validateUser.ValidateUpdatePayload(updatePayload)
				})
			})

			t.Run("Shouldn't raise error when everything is met", func(t *testing.T) {
				// Arrange
				registerPayload := &entity.RegisterUserPayload{
					Username:        "myusername",
					Password:        "mypassword",
					ConfirmPassword: "mypassword",
					Email:           "aiwerkajwe@gmail.com",
				}

				updatePayload := &entity.UpdateUserPayload{
					Username:        "myusername",
					Password:        "mypassword",
					ConfirmPassword: "mypassword",
					Email:           "aiwerkajwe@gmail.com",
				}

				// Action and Assert
				assert.NotPanics(t, func() {
					validateUser.ValidateRegisterPayload(registerPayload)
				})
				assert.NotPanics(t, func() {
					validateUser.ValidateUpdatePayload(updatePayload)
				})
			})
		})

		t.Run("Password Validation", func(t *testing.T) {
			t.Run("Should return error when password is too short", func(t *testing.T) {
				// Arrange
				registerPayload := &entity.RegisterUserPayload{
					Username:        "myusername",
					Password:        "short",
					ConfirmPassword: "Helkdaoskd",
					Email:           "aiwerkajwe@gmail.com",
				}

				updatePayload := &entity.UpdateUserPayload{
					Username:        "myusername",
					Password:        "short",
					ConfirmPassword: "Helkdaoskd",
					Email:           "aiwerkajwe@gmail.com",
				}

				// Action and Assert
				assert.Panics(t, func() {
					validateUser.ValidateRegisterPayload(registerPayload)
				})
				assert.Panics(t, func() {
					validateUser.ValidateUpdatePayload(updatePayload)
				})
			})

			t.Run("Should return error when password is empty", func(t *testing.T) {
				// Arrange
				registerPayload := &entity.RegisterUserPayload{
					Username:        "myusername",
					Password:        "",
					ConfirmPassword: "Helkdaoskd",
					Email:           "aiwerkajwe@gmail.com",
				}

				updatePayload := &entity.UpdateUserPayload{
					Username:        "myusername",
					Password:        "",
					ConfirmPassword: "Helkdaoskd",
					Email:           "aiwerkajwe@gmail.com",
				}

				// Action and Assert
				assert.Panics(t, func() {
					validateUser.ValidateRegisterPayload(registerPayload)
				})
				assert.Panics(t, func() {
					validateUser.ValidateUpdatePayload(updatePayload)
				})
			})

			t.Run("Should return error when confirm password doesn't match", func(t *testing.T) {
				// Arrange
				registerPayload := &entity.RegisterUserPayload{
					Username:        "myusername",
					Password:        "sdawdasdwd",
					ConfirmPassword: "Helkdaoskd",
					Email:           "aiwerkajwe@gmail.com",
				}

				updatePayload := &entity.UpdateUserPayload{
					Username:        "myusername",
					Password:        "sdawdasdwd",
					ConfirmPassword: "Helkdaoskd",
					Email:           "aiwerkajwe@gmail.com",
				}

				// Action and Assert
				assert.Panics(t, func() {
					validateUser.ValidateRegisterPayload(registerPayload)
				})
				assert.Panics(t, func() {
					validateUser.ValidateUpdatePayload(updatePayload)
				})
			})

			t.Run("Shouldn't raise error when password is valid", func(t *testing.T) {
				// Arrange
				registerPayload := &entity.RegisterUserPayload{
					Username:        "myusername",
					Password:        "validpassword",
					ConfirmPassword: "validpassword",
					Email:           "aiwerkajwe@gmail.com",
				}

				updatePayload := &entity.UpdateUserPayload{
					Username:        "myusername",
					Password:        "validpassword",
					ConfirmPassword: "validpassword",
					Email:           "aiwerkajwe@gmail.com",
				}

				// Action and Assert
				assert.NotPanics(t, func() {
					validateUser.ValidateRegisterPayload(registerPayload)
				})
				assert.NotPanics(t, func() {
					validateUser.ValidateUpdatePayload(updatePayload)
				})
			})
		})

		t.Run("Email Validation", func(t *testing.T) {
			t.Run("Should return error when email is invalid", func(t *testing.T) {
				// Arrange
				registerPayload := &entity.RegisterUserPayload{
					Username:        "myusername",
					Password:        "validpassword",
					ConfirmPassword: "validpassword",
					Email:           "invalid-email",
				}

				updatePayload := &entity.UpdateUserPayload{
					Username:        "myusername",
					Password:        "validpassword",
					ConfirmPassword: "validpassword",
					Email:           "invalid-email",
				}

				// Action and Assert
				assert.Panics(t, func() {
					validateUser.ValidateRegisterPayload(registerPayload)
				})
				assert.Panics(t, func() {
					validateUser.ValidateUpdatePayload(updatePayload)
				})
			})

			t.Run("Should return error when email is empty", func(t *testing.T) {
				// Arrange
				registerPayload := &entity.RegisterUserPayload{
					Username:        "myusername",
					Password:        "validpassword",
					ConfirmPassword: "validpassword",
					Email:           "",
				}

				updatePayload := &entity.UpdateUserPayload{
					Username:        "myusername",
					Password:        "validpassword",
					ConfirmPassword: "validpassword",
					Email:           "",
				}

				// Action and Assert
				assert.Panics(t, func() {
					validateUser.ValidateRegisterPayload(registerPayload)
				})
				assert.Panics(t, func() {
					validateUser.ValidateUpdatePayload(updatePayload)
				})
			})

			t.Run("Shouldn't raise error when email is valid", func(t *testing.T) {
				// Arrange
				registerPayload := &entity.RegisterUserPayload{
					Username:        "myusername",
					Password:        "validpassword",
					ConfirmPassword: "validpassword",
					Email:           "valid.email@example.com",
				}

				updatePayload := &entity.UpdateUserPayload{
					Username:        "myusername",
					Password:        "validpassword",
					ConfirmPassword: "validpassword",
					Email:           "valid.email@example.com",
				}

				// Action and Assert
				assert.NotPanics(t, func() {
					validateUser.ValidateRegisterPayload(registerPayload)
				})
				assert.NotPanics(t, func() {
					validateUser.ValidateUpdatePayload(updatePayload)
				})
			})
		})
	})

	t.Run("Login Validation", func(t *testing.T) {

		t.Run("Identity Validation", func(t *testing.T) {
			t.Run("Should return error when identity is empty", func(t *testing.T) {
				// Arrange
				payload := &entity.LoginUserPayload{
					Identity: "",
					Password: "anypassword",
				}

				// Action and Assert
				assert.Panics(t, func() {
					validateUser.ValidateLoginPayload(payload)
				})
			})

			t.Run("Should return error when identity is less than 3", func(t *testing.T) {
				// Arrange
				payload := &entity.LoginUserPayload{
					Identity: "sw",
					Password: "anypassword",
				}

				// Action and Assert
				assert.Panics(t, func() {
					validateUser.ValidateLoginPayload(payload)
				})
			})

			t.Run("Shouldn't raise error when identity is valid", func(t *testing.T) {
				// Arrange
				payload := &entity.LoginUserPayload{
					Identity: "swdw",
					Password: "anypassword",
				}

				// Action and Assert
				assert.NotPanics(t, func() {
					validateUser.ValidateLoginPayload(payload)
				})
			})
		})

		t.Run("Password Validation", func(t *testing.T) {
			t.Run("Should return error when password is too short", func(t *testing.T) {
				// Arrange
				payload := &entity.LoginUserPayload{
					Identity: "anyidentity",
					Password: "short",
				}

				// Action and Assert
				assert.Panics(t, func() {
					validateUser.ValidateLoginPayload(payload)
				})
			})

			t.Run("Should return error when password is empty", func(t *testing.T) {
				// Arrange
				payload := &entity.LoginUserPayload{
					Identity: "anyidentity",
					Password: "",
				}

				// Action and Assert
				assert.Panics(t, func() {
					validateUser.ValidateLoginPayload(payload)
				})
			})

			t.Run("Shouldn't raise error when password is valid", func(t *testing.T) {
				// Arrange
				payload := &entity.LoginUserPayload{
					Identity: "anyidentity",
					Password: "validpassword",
				}

				// Action and Assert
				assert.NotPanics(t, func() {
					validateUser.ValidateLoginPayload(payload)
				})
			})
		})
	})
}
