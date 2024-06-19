package users

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/wisle25/be-template/applications/use_case"
	"github.com/wisle25/be-template/domains/entity"
	"strings"
	"time"
)

type UserHandler struct {
	useCase *use_case.UserUseCase
}

func NewUserHandler(useCase *use_case.UserUseCase) *UserHandler {
	return &UserHandler{
		useCase: useCase,
	}
}

func (h *UserHandler) AddUser(c *fiber.Ctx) error {
	// Payload
	var payload entity.RegisterUserPayload
	_ = c.BodyParser(&payload)

	// Use Case
	returnedId := h.useCase.ExecuteAdd(&payload)

	// Response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"data":    returnedId,
		"message": "Successfully registering new user! Welcome!",
	})
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	// Payload
	var payload entity.LoginUserPayload
	_ = c.BodyParser(&payload)

	// Use Case
	accessTokenDetail, refreshTokenDetail := h.useCase.ExecuteLogin(&payload)

	// Insert the tokens to cookies
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessTokenDetail.Token,
		Path:     "/",
		MaxAge:   accessTokenDetail.MaxAge,
		Secure:   true,
		HTTPOnly: true,
		Domain:   "localhost",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshTokenDetail.Token,
		Path:     "/",
		MaxAge:   refreshTokenDetail.MaxAge,
		Secure:   true,
		HTTPOnly: true,
		Domain:   "localhost",
	})

	// Response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully logged in!",
	})
}

func (h *UserHandler) RefreshToken(c *fiber.Ctx) error {
	// Payload
	refreshToken := c.Cookies("refresh_token")

	// Use Case
	accessTokenDetail := h.useCase.ExecuteRefreshToken(refreshToken)

	// Insert the tokens to cookies
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessTokenDetail.Token,
		Path:     "/",
		MaxAge:   accessTokenDetail.MaxAge,
		Secure:   true,
		HTTPOnly: true,
		Domain:   "localhost",
	})

	// Response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
	})
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
	// Payload
	refreshToken := c.Cookies("refresh_token")
	accessTokenId := c.Locals("accessTokenId").(string)

	// Use Case
	h.useCase.ExecuteLogout(refreshToken, accessTokenId)

	// Remove from cookie
	expiredTime := time.Now().Add(-time.Hour * 24)

	c.Cookie(&fiber.Cookie{
		Name:    "access_token",
		Value:   "",
		Expires: expiredTime,
	})
	c.Cookie(&fiber.Cookie{
		Name:    "refresh_token",
		Value:   "",
		Expires: expiredTime,
	})

	// Response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully logged out!",
	})
}

func (h *UserHandler) GetUserById(c *fiber.Ctx) error {
	// Payload
	id := c.Params("id")

	// Use Case
	user := h.useCase.ExecuteGetUserById(id)

	// Response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}

func (h *UserHandler) GetLoggedUser(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   c.Locals("userInfo").(entity.User),
	})
}

func (h *UserHandler) UpdateUserById(c *fiber.Ctx) error {
	var err error

	// Make sure to update self (not by others)
	id := c.Params("id")
	loggedUserId := c.Locals("userInfo").(entity.User).Id

	if loggedUserId != id {
		return fiber.NewError(
			fiber.StatusForbidden,
			"You are not able to edit other user's profile!",
		)
	}

	// Payload
	var payload entity.UpdateUserPayload
	_ = c.BodyParser(&payload)

	payload.Avatar, err = c.FormFile("avatar")
	if err != nil {
		if !strings.Contains(err.Error(), "there is no uploaded") {
			return fmt.Errorf("upload avatar: %v", err)
		}

		payload.Avatar = nil
	}

	// Use Case
	h.useCase.ExecuteUpdateUserById(id, &payload)

	// Response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully update user!",
	})
}
