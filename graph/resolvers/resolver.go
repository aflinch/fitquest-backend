package resolvers

import "fitquest-backend/database"

type Resolver struct {
	Exercises database.ExerciseRepository
	Users     database.UserRepository
}
