package resolvers

import (
	"database/sql"

	"github.com/kotojo/life-manager/models"
)

//go:generate go run github.com/99designs/gqlgen
type Resolver struct {
	DB               *sql.DB
	UsersService     *models.UserService
	DocumentsService *models.DocumentService
}
