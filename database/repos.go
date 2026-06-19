package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//go:generate go run github.com/vektra/mockery/v2 --name=ExerciseRepository --output=../graph/resolvers/mocks --outpkg=mocks
type ExerciseRepository interface {
	FindAll(ctx context.Context) ([]MongoExercise, error)
	FindFiltered(ctx context.Context, muscle *string) ([]MongoExercise, error)
}

//go:generate go run github.com/vektra/mockery/v2 --name=UserRepository --output=../graph/resolvers/mocks --outpkg=mocks
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

func (r *mongoExerciseRepo) FindFiltered(ctx context.Context, muscle *string) ([]MongoExercise, error) {
	filter := bson.M{}
	if muscle != nil {
		filter = bson.M{
			"$or": []bson.M{
				{"primary_muscles": bson.M{"$in": []string{*muscle}}},
				{"secondary_muscles": bson.M{"$in": []string{*muscle}}},
			},
		}
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
