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

	ur.POST("/:id", h.saveClip)
	ur.GET("/:id", h.getClipsForStreamer)

	views := ur.Group("/views")

	views.GET("", h.getTotalViews)
	views.GET("/perStreamer/:id", h.getTotalViewsByStreamer)
	views.GET("/perStreamer", h.getTotalViewsPerStreamer)
}

func (h *HTTP) saveClip(c echo.Context) error {
	streamerID := c.Param("id")

	err := h.svc.SaveClip(c, streamerID)
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

func (h *HTTP) getTotalViewsByStreamer(c echo.Context) error {
	streamerID := c.Param("id")

	total, err := h.svc.GetTotalViewsByStreamer(c, streamerID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, total)
}

func (h *HTTP) getTotalViews(c echo.Context) error {
	total, err := h.svc.GetTotalViews(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, total)
}

func (h *HTTP) getTotalViewsPerStreamer(c echo.Context) error {
	perStreamer, err := h.svc.GetTotalViewsPerStreamer(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, perStreamer)
}
