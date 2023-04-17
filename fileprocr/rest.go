package fileprocr

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type RestHandler struct {
}

func (h *RestHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/files", h.uploadFile)
}

func (h *RestHandler) uploadFile(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
