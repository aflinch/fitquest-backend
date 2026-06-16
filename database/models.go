package database

import "go.mongodb.org/mongo-driver/bson/primitive"

// MongoUser MongoUsers matches your exact BSON layout in Atlas
type MongoUser struct {
	ID           primitive.ObjectID `bson:"_id"`
	Username     string             `bson:"username"`
	PasswordHash string             `bson:"password_hash"`
}

// MongoExercise matches your exact BSON layout in Atlas
type MongoExercise struct {
	Name              string `bson:"name"`
	EquipmentRequired string `bson:"equipment_required"`
	PrimaryMuscle     string `bson:"primary_muscle"`
}

// MongoWorkout handles BSON database unmarshalling smoothly
type MongoWorkout struct {
	ID                  primitive.ObjectID `bson:"_id"`
	Title               string             `bson:"title"`
	Category            string             `bson:"category"`
	Difficulty          string             `bson:"difficulty"`
	Description         string             `bson:"description"`
	TargetMuscleGroups  []string           `bson:"target_muscle_groups"`
	Exercises           []MongoExercise    `bson:"exercises"`
	IsCommunityTemplate bool               `bson:"is_community_template"`
}
