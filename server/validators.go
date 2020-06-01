package server

import (
	"github.com/sebastianreh/mutants-detector/models"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"regexp"
)

const (
	allowedChard = "ATGC"
)

const (
	minDnaDimension = 4
)

// Expresión regular para verificar que los caracteres dentro de la cadena de Adn sean las permitidas

var (
	r, _ = regexp.Compile(fmt.Sprintf(("^[%s]*$"), allowedChard))
)

type CustomValidator struct {
	validator *validator.Validate
}

// Agrega los validadores a Echo

func SetupValidator(e *echo.Echo) {
	e.Validator = &CustomValidator{validator: validator.New()}
}

// Valida los requests: la estructura del adn y el struct del request

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	request, ok := i.(*models.MutantRequest)
	if ok {
		if err := cv.ValidateDnaChain(*request); err != nil {
			return err
		}
	}
	return nil
}

func (cv *CustomValidator) ValidateDnaChain(request models.MutantRequest) error {
	dimension := len(request.DnaChain)
	for _, dnaBit := range request.DnaChain {
		if err := cv.validateDimension(dnaBit, dimension); err != nil {
			return err
		}
		if err := cv.validateDnaBits(dnaBit); err != nil {
			return err
		}
	}
	return nil
}

// Valida la dimensión NxN

func (cv *CustomValidator) validateDimension(dnaBit string, dimension int) error {
	dnaLength := len(dnaBit)
	if dnaLength != dimension{
		return fmt.Errorf("dimensions of DNA are not NxN")
	}
	if dnaLength < minDnaDimension {
		return fmt.Errorf("dimensions of DNA are not above the minimum required")
	}
	return nil
}

// Valida que las cadenas de adn solo estén formadas por los caracteres permitidos

func (cv *CustomValidator) validateDnaBits(dnaBit string) error {
	if !r.MatchString(dnaBit) {
		return fmt.Errorf("characters inside of DNA are not any combinations of: %s", allowedChard)
	}
	return nil
}
