package graph

import (
	"context"
	"errors"
	"fitquest-backend/auth"
	"fitquest-backend/database"
	"fitquest-backend/graph/model"
	jwtlib "fitquest-backend/jwt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Register is the resolver for the register field.
func (r *mutationResolver) Register(ctx context.Context, username string, password string) (*model.AuthPayload, error) {
	collection := r.DB.Database("fitquest").Collection("users")

	var existing database.MongoUser
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&existing)
	if err == nil {
		return nil, errors.New("username already taken")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return nil, err
	}

	user := database.MongoUser{
		ID:           primitive.NewObjectID(),
		Username:     username,
		PasswordHash: string(hash),
	}

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		return nil, err
	}

	token, err := jwtlib.GenerateToken(user.ID.Hex(), username)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return nil, err
	}

	return &model.AuthPayload{
		Token: token,
		User:  mapMongoUserToGQL(user),
	}, nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, username string, password string) (*model.AuthPayload, error) {
	collection := r.DB.Database("fitquest").Collection("users")

	var user database.MongoUser
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	token, err := jwtlib.GenerateToken(user.ID.Hex(), username)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return nil, err
	}

	return &model.AuthPayload{
		Token: token,
		User:  mapMongoUserToGQL(user),
	}, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }

// RequireAuth checks if a user is authenticated and returns the user context or an error.
// Call this at the start of any resolver that needs authentication.
func RequireAuth(ctx context.Context) (*auth.UserCtx, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return nil, errors.New("authentication required")
	}
	return user, nil
}

// GetUserIdByUsername checks if a user exists in database by given username.
func GetUserIdByUsername(client *mongo.Client, ctx context.Context, username string) (string, error) {
	collection := client.Database("fitquest").Collection("users")

	var user database.MongoUser
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return "", err
	}

	return user.ID.Hex(), nil
}
