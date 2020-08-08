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
func NewHTTP(svc streamer.Service, r *echo.Group, jwtMiddleware echo.MiddlewareFunc) {
	h := HTTP{svc: svc}

	ur := r.Group("/streamers")
	ur.Use(jwtMiddleware)

	ur.POST("/{name}", h.saveStreamer)
	ur.GET("", h.getStreamers)
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

func (h *HTTP) getStreamers(c echo.Context) error {
	streamers, err := h.svc.GetStreamersForUser(c)
	if err != nil {
		// TODO provide better error codes
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, streamers)
}
