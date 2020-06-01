package controllers

import (
	. "ExamenMeLiMutante/settings"
	"github.com/labstack/echo/v4"

	"net/http"
	"time"
)

type Status struct {
	Name    string    `json:"email"`
	Version string    `json:"name"`
	Date    time.Time `json:"date"`
}

func StatusController(c echo.Context) error {
	return c.JSON(http.StatusOK, &Status{
		Name:    ProjectSettings.ProjectName,
		Version: ProjectSettings.ProjectVersion,
		Date:    time.Now(),
	})
}