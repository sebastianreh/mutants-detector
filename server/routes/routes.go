package routes

import (
	"ExamenMeLiMutante/controllers"
	"github.com/labstack/echo/v4"
)

// Aquí se configuran las rutas para el servidor
// Tenemos un controlador de status, para conocer el health-state del ms
// Mutant es la validación del mutante
// Stats son las estadisticas de los mutantes validados

func SetupRoutes(server *echo.Echo, mutantController controllers.IMutantController) {
	base := server.Group("")
	base.GET("/health", controllers.StatusController)
	base.POST("/mutant", mutantController.VerifyMutantStatus)
	base.GET("/stats", mutantController.GetMutantStats)
}
