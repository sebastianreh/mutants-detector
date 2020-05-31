package controllers

import (
	. "ExamenMeLiMutante/settings"
	"github.com/labstack/echo/v4"

	"net/http"
	"time"
)

type Status struct {
	Version string    `json:"name" xml:"name" form:"name" query:"name"`
	Name    string    `json:"email" xml:"email" form:"email" query:"email"`
	Date    time.Time `json:"date" xml:"date" form:"date" query:"date"`
}

func StatusController(c echo.Context) error {
	return c.JSON(http.StatusOK, &Status{
		Name:    ProjectSettings.ProjectName,
		Version: ProjectSettings.ProjectVersion,
		Date:    time.Now(),
	})
}