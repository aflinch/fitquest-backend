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
	Name             string   `bson:"name"`
	Category         string   `bson:"category"`
	Mechanic         string   `bson:"mechanic"`
	PrimaryMuscles   []string `bson:"primary_muscles"`
	SecondaryMuscles []string `bson:"secondary_muscles"`
}
