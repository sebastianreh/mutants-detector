package server

import (
	"github.com/sebastianreh/mutants-detector/container"
	"github.com/sebastianreh/mutants-detector/server/routes"
	. "github.com/sebastianreh/mutants-detector/settings"
	"fmt"
	"github.com/labstack/echo/v4"
)

var (
	url = fmt.Sprintf("%s:%s", ProjectSettings.Host, ProjectSettings.Port)
)

// Inicia el servidor, agrega las validaciones y las rutas

func SetupServer() {
	server := echo.New()
	SetupValidator(server)
	routes.SetupRoutes(server, container.MutantController)
	server.Logger.Fatal(server.Start(url))
}