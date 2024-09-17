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
	e.GET("/issues/:id", handler.GetIssuesByProjectId)
	e.GET("/download/:id", handler.DownloadCsv)
	e.GET("/projects", handler.GetProjects)
	e.GET("/issues", handler.GetIssues)
	e.GET("/users", handler.GetUsers)
	e.GET("/boards", handler.GetBoards)
	e.GET("/sprints", handler.GetSprints)
	e.GET("/attachments", handler.GetAttachments)
	e.GET("/comments", handler.GetComments)
	e.GET("/issue-links", handler.GetIssueLinks)
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
	filename := table + ".csv"
	err := h.usecase.Csv(table)
	defer os.Remove(filename)
	if err != nil {
		return c.JSON(500, err.Error())
	}
	c.Response().Header().Add("Content-Disposition", `attachment; filename="`+filename+`"`)
	c.Response().Header().Add("Content-Type", "text/csv")
	c.Response().Header().Add("Access-Control-Expose-Headers", "Content-Disposition")
	return c.File(filename)
}

func (h *Handler) GetProjects(c echo.Context) error {
	data, err := h.usecase.GetProjects()
	if err != nil {
		return c.JSON(500, "error internal")
	}
	return c.JSON(200, data)
}

func (h *Handler) GetIssues(c echo.Context) error {
	data, err := h.usecase.GetIssues()
	if err != nil {
		return c.JSON(500, "error internal")
	}
	return c.JSON(200, data)
}

func (h *Handler) GetUsers(c echo.Context) error {
	data, err := h.usecase.GetUsers()
	if err != nil {
		return c.JSON(500, "error internal")
	}
	return c.JSON(200, data)
}

func (h *Handler) GetBoards(c echo.Context) error {
	data, err := h.usecase.GetBoards()
	if err != nil {
		return c.JSON(500, "error internal")
	}
	return c.JSON(200, data)
}

func (h *Handler) GetSprints(c echo.Context) error {
	data, err := h.usecase.GetSprints()
	if err != nil {
		return c.JSON(500, "error internal")
	}
	return c.JSON(200, data)
}

func (h *Handler) GetAttachments(c echo.Context) error {
	data, err := h.usecase.GetAttachments()
	if err != nil {
		return c.JSON(500, "error internal")
	}
	return c.JSON(200, data)
}

func (h *Handler) GetComments(c echo.Context) error {
	data, err := h.usecase.GetComments()
	if err != nil {
		return c.JSON(500, "error internal")
	}
	return c.JSON(200, data)
}

func (h *Handler) GetIssueLinks(c echo.Context) error {
	data, err := h.usecase.GetIssueLinks()
	if err != nil {
		return c.JSON(500, "error internal")
	}
	return c.JSON(200, data)
}
