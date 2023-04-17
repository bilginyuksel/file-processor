package fileprocr

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type fileprocrService interface {
	Store(io.Reader) (string, error)
}

type RestHandler struct {
	svc fileprocrService
}

func NewRestHandler(svc fileprocrService) *RestHandler {
	return &RestHandler{svc: svc}
}

func (h *RestHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/files", h.uploadFile)
}

func (h *RestHandler) uploadFile(c echo.Context) error {
	fileheader, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "could not get file"})
	}

	f, err := fileheader.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "could not open file"})
	}
	defer f.Close()

	filename, err := h.svc.Store(f)
	if err != nil {
		zap.L().Error("File could not stored", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not process file"})
	}

	return c.JSON(http.StatusOK, map[string]string{"filename": filename})
}
