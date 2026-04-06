package auth

import (
	"github.com/gofrs/uuid"
)

type Role struct {
	ID          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

type RequestRole struct {
	Name        string `json:"nama" db:"nama"`
	Description string `json:"description" db:"description"`
}

func (role *Role) NewRoleFormat(reqFormat Role) (newRole Role, err error) {
	newID, _ := uuid.NewV4()
	newRole = Role{
		ID:          newID.String(),
		Name:        reqFormat.Name,
		Description: reqFormat.Description,
	}
	return
}
