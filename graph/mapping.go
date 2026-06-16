package graph

import (
	"fitquest-backend/database"
	"fitquest-backend/graph/model"
)

func mapMongoUserToGQL(mu database.MongoUser) *model.User {
	return &model.User{
		ID:       mu.ID.Hex(),
		Username: mu.Username,
	}
}

func mapMongoWorkoutToGQL(mw database.MongoWorkout) *model.Workout {
	exercises := make([]*model.Exercise, 0, len(mw.Exercises))
	for _, me := range mw.Exercises {
		exercises = append(exercises, &model.Exercise{
			Name:              me.Name,
			EquipmentRequired: me.EquipmentRequired,
			PrimaryMuscle:     me.PrimaryMuscle,
		})
	}

	return &model.Workout{
		ID:                  mw.ID.Hex(),
		Title:               mw.Title,
		Category:            mw.Category,
		Difficulty:          mw.Difficulty,
		Description:         mw.Description,
		TargetMuscleGroups:  mw.TargetMuscleGroups,
		Exercises:           exercises,
		IsCommunityTemplate: mw.IsCommunityTemplate,
	}
}
