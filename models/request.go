package models

// Modelo del Request http

type (
	MutantRequest struct {
		DnaChain []string `json:"dna" validate:"required"`
	}
)
