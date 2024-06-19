package db_helper

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/wisle25/be-template/domains/entity"
	"github.com/wisle25/be-template/infrastructures/services"
)

type UserHelperDB struct {
	DB *sql.DB
}

func (h *UserHelperDB) AddUserDB(payload *entity.RegisterUserPayload) {
	id, _ := uuid.NewV7()
	query := "INSERT INTO users(id, username, password, email) VALUES ($1, $2, $3, $4)"
	_, err := h.DB.Exec(
		query,
		id,
		payload.Username,
		payload.Password,
		payload.Email,
	)

	if err != nil {
		panic(err)
	}
}

func (h *UserHelperDB) GetUsers() []entity.User {
	query := "SELECT id, username, email, avatar_link FROM users"
	rows, _ := h.DB.Query(query)

	return services.GetTableDB[entity.User](rows)
}

func (h *UserHelperDB) CleanUserDB() {
	_, _ = h.DB.Query("TRUNCATE TABLE users")
}
