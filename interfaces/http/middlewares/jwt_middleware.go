package middlewares

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/wisle25/be-template/applications/use_case"
	"github.com/wisle25/be-template/domains/entity"
)

// JwtMiddleware gonna verifying user credential for authentication/authorization reason
// Using UserUseCase since it handles both authentication and users itself
type JwtMiddleware struct {
	userUseCase *use_case.UserUseCase
}

func NewJwtMiddleware(userUseCase *use_case.UserUseCase) *JwtMiddleware {
	return &JwtMiddleware{userUseCase}
}

func (m *JwtMiddleware) GuardJWT(c *fiber.Ctx) error {
	// Getting access token
	accessToken := c.Cookies("access_token")

	if accessToken == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "You are not logged in!")
	}

	userInfoJSON, accessTokenDetail := m.userUseCase.ExecuteGuard(accessToken)

	if userInfoJSON == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Session invalid or expired!")
	}

	// Unmarshal userInfo JSON
	var userInfo entity.User
	err := json.Unmarshal([]byte(userInfoJSON.(string)), &userInfo)
	if err != nil {
		panic(fmt.Errorf("refresh_token_err: unable to unmarshal json user info: %v", err))
	}

	// Add additional information
	c.Locals("accessTokenId", accessTokenDetail.TokenId)
	c.Locals("userInfo", userInfo)

	return c.Next()
}
