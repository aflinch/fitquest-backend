package database

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExerciseRepository interface {
	FindAll(ctx context.Context) ([]MongoExercise, error)
	FindById(ctx context.Context, id *string) ([]MongoExercise, error)
	FindFiltered(ctx context.Context, muscle *string, ids []*string) ([]MongoExercise, error)
}

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*MongoUser, error)
	Insert(ctx context.Context, user MongoUser) error
}

type mongoExerciseRepo struct {
	col *mongo.Collection
}

func NewExerciseRepo(db *mongo.Database) ExerciseRepository {
	return &mongoExerciseRepo{col: db.Collection("exercises")}
}

func (r *mongoExerciseRepo) FindAll(ctx context.Context) ([]MongoExercise, error) {
	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var exercises []MongoExercise
	if err := cursor.All(ctx, &exercises); err != nil {
		return nil, err
	}
	return exercises, nil
}

func (r *mongoExerciseRepo) FindById(ctx context.Context, id *string) ([]MongoExercise, error) {
	if id == nil {
		return nil, fmt.Errorf("id cannot be nil")
	}

	objID, err := primitive.ObjectIDFromHex(*id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %v", err)
	}

	filter := bson.M{"_id": objID}

	cursor, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var exercises []MongoExercise
	if err := cursor.All(ctx, &exercises); err != nil {
		return nil, err
	}
	return exercises, nil
}

func (r *mongoExerciseRepo) FindFiltered(ctx context.Context, muscle *string, ids []*string) ([]MongoExercise, error) {
	var objectIDs []primitive.ObjectID

	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(*id)
		if err != nil {
			return nil, fmt.Errorf("invalid id format: %v", err)
		}
		objectIDs = append(objectIDs, objID)
	}

	var orConditions []bson.M

	if muscle != nil {
		lowercaseMuscle := strings.ToLower(*muscle)
		orConditions = append(orConditions, bson.M{
			"$or": []bson.M{
				{"primary_muscles": bson.M{"$in": []string{lowercaseMuscle}}},
				{"secondary_muscles": bson.M{"$in": []string{lowercaseMuscle}}},
			},
		})
	}

	if len(objectIDs) > 0 {
		orConditions = append(orConditions, bson.M{
			"_id": bson.M{"$in": objectIDs},
		})
	}

	filter := bson.M{}
	if len(orConditions) > 0 {
		filter = bson.M{"$or": orConditions}
	}

	cursor, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var exercises []MongoExercise
	if err := cursor.All(ctx, &exercises); err != nil {
		return nil, err
	}
	return exercises, nil
}

type mongoUserRepo struct {
	col *mongo.Collection
}

func NewUserRepo(db *mongo.Database) UserRepository {
	return &mongoUserRepo{col: db.Collection("users")}
}

func (r *mongoUserRepo) FindByUsername(ctx context.Context, username string) (*MongoUser, error) {
	var user MongoUser
	err := r.col.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *mongoUserRepo) Insert(ctx context.Context, user MongoUser) error {
	_, err := r.col.InsertOne(ctx, user)
	return err
}
