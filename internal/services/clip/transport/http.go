package transport

import (
	"github.com/labstack/echo"
	"net/http"
	"streamerEventViewer/internal/services/clip"
)

// HTTP represents user http service
type HTTP struct {
	svc clip.Service
}

// NewHTTP creates new user http service
func NewHTTP(svc clip.Service, r *echo.Group, jwtMiddleware echo.MiddlewareFunc) {
	h := HTTP{svc: svc}

	ur := r.Group("/clips")
	ur.Use(jwtMiddleware)

	// TODO should use consist way(one type of ids for each request)
	// id - external(twitch) id for streamer
	ur.POST("/{id}", h.saveClip)
	// id - our own generated streamer id
	ur.GET("/{id}", h.getClipsForStreamer)
}

func (h *HTTP) saveClip(c echo.Context) error {
	externalStreamerID := c.Param("id")

	err := h.svc.SaveClip(c, externalStreamerID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (h *HTTP) getClipsForStreamer(c echo.Context) error {
	streamerID := c.Param("id")

	clips, err := h.svc.GetClipsForStreamer(c, streamerID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, clips)
}
