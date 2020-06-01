package repositories

import (
	"github.com/sebastianreh/mutants-detector/models"
	"github.com/sebastianreh/mutants-detector/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Genera sujetos a partir de Dna_Id de manera asincrónica

func generateSubjectsAsync(subjectFunc func([]string, bool) *[]models.Subject, subjectsIds []string, status string) <-chan *[]models.Subject {
	future := make(chan *[]models.Subject)
	var isMutant bool
	if status == MutantStatus {
		isMutant = true
	} else {
		isMutant = false
	}
	go func() {
		future <- subjectFunc(subjectsIds, isMutant)
		close(future)
	}()

	return future
}

// Genera sujetos a partir de Dna_Id

func generateSubjects(mutantsIds []string, isMutant bool) *[]models.Subject {
	subjects := new([]models.Subject)
	for _, mutantId := range mutantsIds {
		mutant := models.Subject{
			Id:       mutantId,
			Dna:      utils.ConvertIdToDna(mutantId),
			IsMutant: isMutant,
		}
		*subjects = append(*subjects, mutant)
	}
	return subjects
}

// Genera los modelos para insertar en Mongo

func BuildSubjectModelDB(subject models.Subject) mongo.WriteModel {
	model := mongo.NewUpdateOneModel()
	model.Filter = bson.M{
		"dna_id": subject.Id,
	}
	model.Update = bson.M{
		"$setOnInsert": subject,
	}
	model.Upsert = &upsert
	return model
}

// Devuelve Keys para la transacción en redis segun el sujeto

func selectConditionKey(subject models.Subject) string {
	if subject.IsMutant {
		return MutantsNotSavedKey
	} else {
		return HumansNotSavedKey
	}
}

// Devuelve Keys para la transacción en redis segun el status

func selectKeysByStatus(status string) (string, string) {
	if status == MutantStatus {
		return MutantsSavedKey, MutantsNotSavedKey
	} else {
		return HumansSavedKey, HumansNotSavedKey
	}
}
