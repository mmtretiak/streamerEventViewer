package transport

import (
	"net/http"
	"streamerEventViewer/internal/services/user"

	"github.com/labstack/echo"
)

// HTTP represents user http service
type HTTP struct {
	svc user.Service
}

// NewHTTP creates new user http service
func NewHTTP(svc user.Service, r *echo.Group) {
	h := HTTP{svc}
	ur := r.Group("/users")

	ur.GET("/login", h.login)
	ur.GET("/login/redirect", h.redirect)
}

func (h *HTTP) login(c echo.Context) error {
	err := h.svc.Login(c)
	if err != nil {
		// TODO provide better error codes
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}

func (h *HTTP) redirect(c echo.Context) error {
	token, err := h.svc.Redirect(c)
	if err != nil {
		// TODO provide better error codes
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, &token)
}