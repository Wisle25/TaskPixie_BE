package security

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/matthewhartstonge/argon2"
	"github.com/wisle25/be-template/applications/security"
)

// Argon2PasswordHash Using argon2 for password hashing
type Argon2PasswordHash struct /* implements PasswordHash */ {
	argon argon2.Config
}

func NewArgon2() security.PasswordHash {
	return &Argon2PasswordHash{
		argon: argon2.DefaultConfig(),
	}
}

func (a *Argon2PasswordHash) Hash(password string) string {
	hashed, err := a.argon.HashEncoded([]byte(password))

	if err != nil {
		panic(fmt.Errorf("argon2_err: hash: %v", err))
	}

	return string(hashed)
}

func (a *Argon2PasswordHash) Compare(password string, hashedPassword string) {
	match, err := argon2.VerifyEncoded([]byte(password), []byte(hashedPassword))

	if err != nil {
		panic(fmt.Errorf("argon2_err: compare: %v", err))
	}

	if !match {
		panic(fiber.NewError(fiber.StatusUnauthorized, "Password is incorrect!"))
	}
}
