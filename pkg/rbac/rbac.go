package rbac

import (
	"github.com/labstack/echo"
	"streamerEventViewer/pkg/models"
)

// RBACService represents role-based access control service interface
type RBACService interface {
	User(echo.Context) models.User
}

func New() RBACService {
	return &service{}
}

type service struct{}

func (s *service) User(c echo.Context) models.User {
	id := c.Get("id").(string)
	name := c.Get("name").(string)
	email := c.Get("email").(string)

	return models.User{
		ID:    id,
		Name:  name,
		Email: email,
	}
}
