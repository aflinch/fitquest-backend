package tests

import (
	"context"
	"errors"
	"fitquest-backend/database"
	"fitquest-backend/graph/resolvers"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepo struct {
	findByUsername func(ctx context.Context, username string) (*database.MongoUser, error)
	insert         func(ctx context.Context, user database.MongoUser) error
}

func (m *mockUserRepo) FindByUsername(ctx context.Context, username string) (*database.MongoUser, error) {
	return m.findByUsername(ctx, username)
}

func (m *mockUserRepo) Insert(ctx context.Context, user database.MongoUser) error {
	return m.insert(ctx, user)
}

func TestRegister_UsernameTaken(t *testing.T) {
	existing := &database.MongoUser{
		ID:           primitive.NewObjectID(),
		Username:     "existinguser",
		PasswordHash: "hash",
	}

	r := &resolvers.MutationResolver{Resolver: &resolvers.Resolver{
		Users: &mockUserRepo{
			findByUsername: func(ctx context.Context, username string) (*database.MongoUser, error) {
				return existing, nil
			},
			insert: func(ctx context.Context, user database.MongoUser) error {
				return nil
			},
		},
	}}

	_, err := r.Register(context.Background(), "existinguser", "password123")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "username already taken" {
		t.Errorf("expected 'username already taken', got %q", err.Error())
	}
}

func TestRegister_Success(t *testing.T) {
	var insertedUser database.MongoUser
	r := &resolvers.MutationResolver{Resolver: &resolvers.Resolver{
		Users: &mockUserRepo{
			findByUsername: func(ctx context.Context, username string) (*database.MongoUser, error) {
				return nil, nil
			},
			insert: func(ctx context.Context, user database.MongoUser) error {
				insertedUser = user
				return nil
			},
		},
	}}

	payload, err := r.Register(context.Background(), "newuser", "password123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if payload.User.Username != "newuser" {
		t.Errorf("Username = %q, want %q", payload.User.Username, "newuser")
	}
	if payload.Token == "" {
		t.Error("expected non-empty token")
	}
	if insertedUser.Username != "newuser" {
		t.Errorf("inserted Username = %q, want %q", insertedUser.Username, "newuser")
	}
	if insertedUser.PasswordHash == "" {
		t.Error("expected non-empty password hash")
	}
	if insertedUser.PasswordHash == "password123" {
		t.Error("password should be hashed, not stored in plaintext")
	}
}

func TestRegister_FindByUsernameError(t *testing.T) {
	r := &resolvers.MutationResolver{Resolver: &resolvers.Resolver{
		Users: &mockUserRepo{
			findByUsername: func(ctx context.Context, username string) (*database.MongoUser, error) {
				return nil, errors.New("connection failed")
			},
			insert: func(ctx context.Context, user database.MongoUser) error {
				return nil
			},
		},
	}}

	_, err := r.Register(context.Background(), "newuser", "password123")
	if err == nil || err.Error() != "connection failed" {
		t.Fatalf("expected 'connection failed', got %v", err)
	}
}

func TestRegister_InsertError(t *testing.T) {
	r := &resolvers.MutationResolver{Resolver: &resolvers.Resolver{
		Users: &mockUserRepo{
			findByUsername: func(ctx context.Context, username string) (*database.MongoUser, error) {
				return nil, nil
			},
			insert: func(ctx context.Context, user database.MongoUser) error {
				return errors.New("insert failed")
			},
		},
	}}

	_, err := r.Register(context.Background(), "newuser", "password123")
	if err == nil || err.Error() != "insert failed" {
		t.Fatalf("expected 'insert failed', got %v", err)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	r := &resolvers.MutationResolver{Resolver: &resolvers.Resolver{
		Users: &mockUserRepo{
			findByUsername: func(ctx context.Context, username string) (*database.MongoUser, error) {
				return nil, nil
			},
		},
	}}

	_, err := r.Login(context.Background(), "unknown", "password123")
	if err == nil || err.Error() != "invalid username or password" {
		t.Fatalf("expected 'invalid username or password', got %v", err)
	}
}

func TestLogin_FindByUsernameError(t *testing.T) {
	r := &resolvers.MutationResolver{Resolver: &resolvers.Resolver{
		Users: &mockUserRepo{
			findByUsername: func(ctx context.Context, username string) (*database.MongoUser, error) {
				return nil, errors.New("connection failed")
			},
		},
	}}

	_, err := r.Login(context.Background(), "someone", "password123")
	if err == nil || err.Error() != "invalid username or password" {
		t.Fatalf("expected 'invalid username or password', got %v", err)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	user := &database.MongoUser{
		ID:           primitive.NewObjectID(),
		Username:     "testuser",
		PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // bcrypt hash of "correctpassword"
	}

	r := &resolvers.MutationResolver{&resolvers.Resolver{
		Users: &mockUserRepo{
			findByUsername: func(ctx context.Context, username string) (*database.MongoUser, error) {
				return user, nil
			},
		},
	}}

	_, err := r.Login(context.Background(), "testuser", "wrongpassword")
	if err == nil || err.Error() != "invalid username or password" {
		t.Fatalf("expected 'invalid username or password', got %v", err)
	}
}

func TestLogin_Success(t *testing.T) {
	hash, err := bcryptHash("correctpassword")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	user := &database.MongoUser{
		ID:           primitive.NewObjectID(),
		Username:     "testuser",
		PasswordHash: hash,
	}

	r := &resolvers.MutationResolver{Resolver: &resolvers.Resolver{
		Users: &mockUserRepo{
			findByUsername: func(ctx context.Context, username string) (*database.MongoUser, error) {
				return user, nil
			},
		},
	}}

	payload, err := r.Login(context.Background(), "testuser", "correctpassword")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if payload.User.Username != "testuser" {
		t.Errorf("Username = %q, want %q", payload.User.Username, "testuser")
	}
	if payload.Token == "" {
		t.Error("expected non-empty token")
	}
}

func bcryptHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
