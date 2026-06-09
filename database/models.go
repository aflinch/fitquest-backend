package database

import "go.mongodb.org/mongo-driver/v2/bson"

// ExerciseDB maps the nested array items in Atlas
type ExerciseDB struct {
	Name              string `bson:"name"`
	EquipmentRequired string `bson:"equipment_required"`
	PrimaryMuscle     string `bson:"primary_muscle"`
}

// WorkoutTemplateDB maps precisely to your screenshot's fields
type WorkoutTemplateDB struct {
	ID                  bson.ObjectID `bson:"_id"`
	Title               string        `bson:"title"`
	Category            string        `bson:"category"`
	Difficulty          string        `bson:"difficulty"`
	Description         string        `bson:"description"`
	TargetMuscleGroups  []string      `bson:"target_muscle_groups"`
	Exercises           []ExerciseDB  `bson:"exercises"`
	IsCommunityTemplate bool          `bson:"is_community_template"`
}
