package models

// Modelos de los sujetos, prestats (para su calculo) y stats

type (
	Subject struct {
		Id       string   `json:"dna_id" bson:"dna_id"`
		Dna      []string `json:"dna" bson:"dna"`
		IsMutant bool     `json:"is_mutant" bson:"is_mutant"`
	}

	MutantsPreStats struct {
		CountMutants int
		CountHumans  int
	}

	MutantsStats struct {
		CountMutantDna float64 `json:"count_mutant_dna"`
		CountHumanDna  float64 `json:"count_human_dna"`
		Ratio          float64 `json:"ratio"`
	}
)
