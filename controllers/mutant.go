package controllers

import (
	"ExamenMeLiMutante/models"
	"ExamenMeLiMutante/services"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type (
	MutantController struct {
		service services.IMutantService
	}

	IMutantController interface {
		VerifyMutantStatus(c echo.Context) error
		GetMutantStats(c echo.Context) error
	}
)

func NewMutantController(mutantService services.IMutantService) IMutantController {
	return MutantController{
		service: mutantService}
}

// Controller para verificar el estado del sujeto, devolviendo
// el http.Status correspondiente según la condición del sujeto y posteriormente

func (controller MutantController) VerifyMutantStatus(c echo.Context) error {
	mutantRequest := new(models.MutantRequest)
	mutantResponse := models.MutantResponse{}
	if err := c.Bind(&mutantRequest); err != nil {
		log.Errorf("controllers.VerifyMutantStatus | Error in VerifyMutantStatus, bad request: %v", err)
		echoErr := echo.NewHTTPError(http.StatusBadRequest, err.Error())
		return c.JSON(echoErr.Code, echoErr)
	}
	if err := c.Validate(mutantRequest); err != nil {
		log.Errorf("controllers.VerifyMutantStatus | Error in VerifyMutantStatus, bad request: %v", err)
		echoErr := echo.NewHTTPError(http.StatusBadRequest, err.Error())
		return c.JSON(echoErr.Code, echoErr)
	}
	mutantResponse = controller.service.VerifyMutant(*mutantRequest)
	if !mutantResponse.IsMutant {
		return c.JSON(http.StatusForbidden, mutantResponse)
	}
	return c.JSON(http.StatusOK, mutantResponse)
}

// Controlador de stats

func (controller MutantController) GetMutantStats(c echo.Context) error {
	mutantStats, err := controller.service.GetSubjectsStats()
	if err != nil {
		log.Errorf("Error in GetMutantStats: internal error: %v", err)
		echoErr :=  echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		return c.JSON(echoErr.Code, echoErr)
	}
	return c.JSON(http.StatusOK, mutantStats)
}
