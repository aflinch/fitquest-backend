package mapping

import (
	"fitquest-backend/database"
	"fitquest-backend/graph/helper"
	"fitquest-backend/graph/model"
)

func MapMongoUserToGQL(mu database.MongoUser) *model.User {
	return &model.User{
		ID:       mu.ID.Hex(),
		Username: mu.Username,
	}
}

func MapMongoExerciseToGQL(me database.MongoExercise) *model.Exercise {
	return &model.Exercise{
		ID:               me.ID.Hex(),
		Name:             me.Name,
		Category:         &me.Category,
		Mechanic:         &me.Mechanic,
		PrimaryMuscles:   helper.ToPtrSlice(me.PrimaryMuscles),
		SecondaryMuscles: helper.ToPtrSlice(me.SecondaryMuscles),
	}
}
