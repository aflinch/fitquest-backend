package tests

import (
	"context"
	"errors"
	"fitquest-backend/auth"
	"fitquest-backend/database"
	"fitquest-backend/graph/model"
	"fitquest-backend/graph/resolvers"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockExerciseRepo struct {
	findAll      func(ctx context.Context) ([]database.MongoExercise, error)
	findFiltered func(ctx context.Context, muscle *string) ([]database.MongoExercise, error)
}

func (m *mockExerciseRepo) FindAll(ctx context.Context) ([]database.MongoExercise, error) {
	return m.findAll(ctx)
}

func (m *mockExerciseRepo) FindFiltered(ctx context.Context, muscle *string) ([]database.MongoExercise, error) {
	return m.findFiltered(ctx, muscle)
}

func authContext() context.Context {
	return auth.WithUserCtx(context.Background(), &auth.UserCtx{
		UserID:   "test-id",
		Username: "testuser",
	})
}

func TestGetExercises_AuthRequired(t *testing.T) {
	r := &resolvers.QueryResolver{Resolver: &resolvers.Resolver{
		Exercises: &mockExerciseRepo{
			findAll: func(ctx context.Context) ([]database.MongoExercise, error) {
				return nil, nil
			},
		},
	}}

	_, err := r.GetExercises(context.Background())
	if err == nil || err.Error() != "authentication required" {
		t.Fatalf("expected auth error, got %v", err)
	}
}

func TestGetExercises_Success(t *testing.T) {
	expected := []database.MongoExercise{
		{
			ID:               primitive.NewObjectID(),
			Name:             "Bench Press",
			Category:         "Strength",
			Mechanic:         "Compound",
			PrimaryMuscles:   []string{"Chest", "Triceps"},
			SecondaryMuscles: []string{"Front Delts"},
		},
	}

	r := &resolvers.QueryResolver{&resolvers.Resolver{
		Exercises: &mockExerciseRepo{
			findAll: func(ctx context.Context) ([]database.MongoExercise, error) {
				return expected, nil
			},
		},
	}}

	got, err := r.GetExercises(authContext())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 exercise, got %d", len(got))
	}
	if got[0].Name != "Bench Press" {
		t.Errorf("Name = %q, want %q", got[0].Name, "Bench Press")
	}
	if *got[0].Category != "Strength" {
		t.Errorf("Category = %q, want %q", *got[0].Category, "Strength")
	}
}

func TestGetExercises_Empty(t *testing.T) {
	r := &resolvers.QueryResolver{&resolvers.Resolver{
		Exercises: &mockExerciseRepo{
			findAll: func(ctx context.Context) ([]database.MongoExercise, error) {
				return []database.MongoExercise{}, nil
			},
		},
	}}

	got, err := r.GetExercises(authContext())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 0 {
		t.Fatalf("expected 0 exercises, got %d", len(got))
	}
}

func TestGetExercises_RepoError(t *testing.T) {
	r := &resolvers.QueryResolver{&resolvers.Resolver{
		Exercises: &mockExerciseRepo{
			findAll: func(ctx context.Context) ([]database.MongoExercise, error) {
				return nil, errors.New("db error")
			},
		},
	}}

	_, err := r.GetExercises(authContext())
	if err == nil || err.Error() != "db error" {
		t.Fatalf("expected db error, got %v", err)
	}
}

func TestGetFilteredExercises_AuthRequired(t *testing.T) {
	r := &resolvers.QueryResolver{&resolvers.Resolver{
		Exercises: &mockExerciseRepo{
			findFiltered: func(ctx context.Context, muscle *string) ([]database.MongoExercise, error) {
				return nil, nil
			},
		},
	}}

	_, err := r.GetFilteredExercises(context.Background(), nil)
	if err == nil || err.Error() != "authentication required" {
		t.Fatalf("expected auth error, got %v", err)
	}
}

func TestGetFilteredExercises_NilWhere(t *testing.T) {
	var calledWith *string
	r := &resolvers.QueryResolver{&resolvers.Resolver{
		Exercises: &mockExerciseRepo{
			findFiltered: func(ctx context.Context, muscle *string) ([]database.MongoExercise, error) {
				calledWith = muscle
				return []database.MongoExercise{
					{ID: primitive.NewObjectID(), Name: "Push Up"},
				}, nil
			},
		},
	}}

	got, err := r.GetFilteredExercises(authContext(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if calledWith != nil {
		t.Errorf("expected nil muscle filter, got %v", *calledWith)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 exercise, got %d", len(got))
	}
}

func TestGetFilteredExercises_WithMuscle(t *testing.T) {
	muscle := "Chest"
	r := &resolvers.QueryResolver{&resolvers.Resolver{
		Exercises: &mockExerciseRepo{
			findFiltered: func(ctx context.Context, m *string) ([]database.MongoExercise, error) {
				if m == nil || *m != "Chest" {
					t.Fatalf("expected muscle=Chest, got %v", m)
				}
				return []database.MongoExercise{
					{ID: primitive.NewObjectID(), Name: "Bench Press"},
				}, nil
			},
		},
	}}

	got, err := r.GetFilteredExercises(authContext(), &model.ExerciseFilter{Muscle: &muscle})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 exercise, got %d", len(got))
	}
	if got[0].Name != "Bench Press" {
		t.Errorf("Name = %q, want %q", got[0].Name, "Bench Press")
	}
}

func TestGetFilteredExercises_EmptyResult(t *testing.T) {
	r := &resolvers.QueryResolver{&resolvers.Resolver{
		Exercises: &mockExerciseRepo{
			findFiltered: func(ctx context.Context, muscle *string) ([]database.MongoExercise, error) {
				return []database.MongoExercise{}, nil
			},
		},
	}}

	got, err := r.GetFilteredExercises(authContext(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 0 {
		t.Fatalf("expected 0 exercises, got %d", len(got))
	}
}

func TestGetFilteredExercises_RepoError(t *testing.T) {
	r := &resolvers.QueryResolver{&resolvers.Resolver{
		Exercises: &mockExerciseRepo{
			findFiltered: func(ctx context.Context, muscle *string) ([]database.MongoExercise, error) {
				return nil, errors.New("db error")
			},
		},
	}}

	_, err := r.GetFilteredExercises(authContext(), nil)
	if err == nil || err.Error() != "db error" {
		t.Fatalf("expected db error, got %v", err)
	}
}
