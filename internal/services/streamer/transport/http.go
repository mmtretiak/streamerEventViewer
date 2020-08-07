package transport

import (
	"github.com/labstack/echo"
	"net/http"
	"streamerEventViewer/internal/services/streamer"
)

// HTTP represents user http service
type HTTP struct {
	svc streamer.Service
}

// NewHTTP creates new user http service
func NewHTTP(svc streamer.Service, r *echo.Group) {
	h := HTTP{svc: svc}
	ur := r.Group("/streamer")

	ur.POST("/{name}", h.saveStreamer)
}

func (h *HTTP) saveStreamer(c echo.Context) error {
	streamerName := c.Param("name")

	err := h.svc.SaveStreamer(c, streamerName)
	if err != nil {
		// TODO provide better error codes
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}
