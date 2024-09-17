package http

import (
	"ms-sv-jira/helper/logger"
	"ms-sv-jira/middleware"
	"ms-sv-jira/module"
	"os"

	"github.com/labstack/echo/v4"
)

type ResponseError struct {
	Message string
}

type Handler struct {
	usecase module.Usecase
	log     logger.Logger
}

func NewHandler(e *echo.Echo, middl *middleware.GoMiddleware, usecase module.Usecase, log logger.Logger) {
	handler := &Handler{
		log:     log,
		usecase: usecase,
	}
	e.POST("/issues/:id", handler.GetIssuesByProjectId)
	e.POST("/download/:id", handler.DownloadCsv)
}

func (h *Handler) GetIssuesByProjectId(c echo.Context) error {
	id := c.Param("id")
	data, err := h.usecase.GetIssuesByProjectId(id)
	if err != nil {
		return c.JSON(500, "error internal")
	}
	return c.JSON(200, data)
}

func (h *Handler) DownloadCsv(c echo.Context) error {
	table := c.Param("id")
	h.usecase.Csv(table)
	filename := table + ".csv"
	defer os.Remove(filename)
	c.Response().Header().Add("Content-Disposition", `attachment; filename="`+filename+`"`)
	c.Response().Header().Add("Access-Control-Expose-Headers", "Content-Disposition")
	return c.File(filename)
}
