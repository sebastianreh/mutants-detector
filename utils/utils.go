package utils

import (
	"github.com/sebastianreh/mutants-detector/models"
	"fmt"
	"strconv"
)

// Convierte Adn a Id -> El número es la dimensión de NxN, los strings son las candeas de Adn

func ConvertDnaToId(dna []string) string {
	dimension := len(dna)
	var content string
	{
		for _, dnaBit := range dna {
			content = content + dnaBit
		}
	}
	return fmt.Sprintf("%d%s", dimension, content)
}

// Convierte Id a Adn -> Realiza el paso inverso de la función ConvertDnaToId

func ConvertIdToDna(id string) []string {
	var dna []string
	dimension, _ := strconv.Atoi(id[:1])
	primitiveDna := id[1:]
	for i := 0; i < dimension; i++ {
		dna = append(dna, primitiveDna[i*dimension:(i+1)*dimension])
	}
	return dna
}

// Calcula las stats con las preStats

func CalculateMutantStats(preStats models.MutantsPreStats) *models.MutantsStats {
	var ratio float64
	var mutants = float64(preStats.CountMutants)
	var humans = float64(preStats.CountHumans)

	if preStats.CountMutants == 0 || preStats.CountHumans == 0 {
		ratio = 1
	}
	if preStats.CountMutants == 0 && preStats.CountHumans == 0 {
		ratio = 0
	}
	if preStats.CountMutants != 0 && preStats.CountHumans != 0 {
		ratio = mutants / humans
	}

	return &models.MutantsStats{
		CountMutantDna: mutants,
		CountHumanDna:  humans,
		Ratio:          ratio,
	}
}