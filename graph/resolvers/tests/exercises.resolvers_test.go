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
	findById     func(ctx context.Context, id *string) ([]database.MongoExercise, error)
	findFiltered func(ctx context.Context, muscle *string, ids []*string) ([]database.MongoExercise, error)
}

func (m *mockExerciseRepo) FindById(ctx context.Context, id *string) ([]database.MongoExercise, error) {
	return m.findById(ctx, id)
}

func (m *mockExerciseRepo) FindFiltered(ctx context.Context, muscle *string, ids []*string) ([]database.MongoExercise, error) {
	return m.findFiltered(ctx, muscle, ids)
}

func (m *mockExerciseRepo) FindAll(ctx context.Context) ([]database.MongoExercise, error) {
	return m.findAll(ctx)
}

func authContext() context.Context {
	return auth.WithUserCtx(context.Background(), &auth.UserCtx{
		UserID:   "test-id",
		Username: "testuser",
	})
}

func TestExercises_AuthRequired(t *testing.T) {
	r := &resolvers.QueryResolver{Resolver: &resolvers.Resolver{
		Exercises: &mockExerciseRepo{
			findAll: func(ctx context.Context) ([]database.MongoExercise, error) {
				return nil, nil
			},
		},
	}}

	_, err := r.Exercises(context.Background(), nil)
	if err == nil || err.Error() != "authentication required" {
		t.Fatalf("expected auth error, got %v", err)
	}
}

func TestExercises_Success(t *testing.T) {
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

	got, err := r.Exercises(authContext(), nil)
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

func TestExercises_Empty(t *testing.T) {
	r := &resolvers.QueryResolver{&resolvers.Resolver{
		Exercises: &mockExerciseRepo{
			findAll: func(ctx context.Context) ([]database.MongoExercise, error) {
				return []database.MongoExercise{}, nil
			},
		},
	}}

	got, err := r.Exercises(authContext(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 0 {
		t.Fatalf("expected 0 exercises, got %d", len(got))
	}
}

func TestExercises_RepoError(t *testing.T) {
	r := &resolvers.QueryResolver{&resolvers.Resolver{
		Exercises: &mockExerciseRepo{
			findAll: func(ctx context.Context) ([]database.MongoExercise, error) {
				return nil, errors.New("db error")
			},
		},
	}}

	_, err := r.Exercises(authContext(), nil)
	if err == nil || err.Error() != "db error" {
		t.Fatalf("expected db error, got %v", err)
	}
}

func TestExercises_ByID_Success(t *testing.T) {
	exercise := database.MongoExercise{
		ID:               primitive.NewObjectID(),
		Name:             "Squat",
		Category:         "Strength",
		Mechanic:         "Compound",
		PrimaryMuscles:   []string{"Quads", "Glutes"},
		SecondaryMuscles: []string{"Hamstrings"},
	}
	id := exercise.ID.Hex()

	r := &resolvers.QueryResolver{&resolvers.Resolver{
		Exercises: &mockExerciseRepo{
			findById: func(ctx context.Context, gotID *string) ([]database.MongoExercise, error) {
				if gotID == nil || *gotID != id {
					t.Fatalf("expected id=%s, got %v", id, gotID)
				}
				return []database.MongoExercise{exercise}, nil
			},
		},
	}}

	got, err := r.Exercises(authContext(), &id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 exercise, got %d", len(got))
	}
	if got[0].Name != "Squat" {
		t.Errorf("Name = %q, want %q", got[0].Name, "Squat")
	}
}

func TestExercises_ByID_NotFound(t *testing.T) {
	id := primitive.NewObjectID().Hex()
	r := &resolvers.QueryResolver{&resolvers.Resolver{
		Exercises: &mockExerciseRepo{
			findById: func(ctx context.Context, gotID *string) ([]database.MongoExercise, error) {
				return []database.MongoExercise{}, nil
			},
		},
	}}

	got, err := r.Exercises(authContext(), &id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected 0 exercises, got %d", len(got))
	}
}

func TestExercises_ByID_RepoError(t *testing.T) {
	id := primitive.NewObjectID().Hex()
	r := &resolvers.QueryResolver{&resolvers.Resolver{
		Exercises: &mockExerciseRepo{
			findById: func(ctx context.Context, gotID *string) ([]database.MongoExercise, error) {
				return nil, errors.New("db error")
			},
		},
	}}

	_, err := r.Exercises(authContext(), &id)
	if err == nil || err.Error() != "db error" {
		t.Fatalf("expected db error, got %v", err)
	}
}

func TestFilteredExercises_AuthRequired(t *testing.T) {
	r := &resolvers.QueryResolver{&resolvers.Resolver{
		Exercises: &mockExerciseRepo{
			findFiltered: func(ctx context.Context, muscle *string, ids []*string) ([]database.MongoExercise, error) {
				return nil, nil
			},
		},
	}}

	_, err := r.FilteredExercises(context.Background(), model.ExerciseFilter{})
	if err == nil || err.Error() != "authentication required" {
		t.Fatalf("expected auth error, got %v", err)
	}
}

func TestFilteredExercises_NilWhere(t *testing.T) {
	var calledWith *string
	r := &resolvers.QueryResolver{&resolvers.Resolver{
		Exercises: &mockExerciseRepo{
			findFiltered: func(ctx context.Context, muscle *string, ids []*string) ([]database.MongoExercise, error) {
				calledWith = muscle
				return []database.MongoExercise{
					{ID: primitive.NewObjectID(), Name: "Push Up"},
				}, nil
			},
		},
	}}

	got, err := r.FilteredExercises(authContext(), model.ExerciseFilter{})
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

func TestFilteredExercises_WithMuscle(t *testing.T) {
	muscle := "Chest"
	r := &resolvers.QueryResolver{&resolvers.Resolver{
		Exercises: &mockExerciseRepo{
			findFiltered: func(ctx context.Context, m *string, ids []*string) ([]database.MongoExercise, error) {
				if m == nil || *m != "Chest" {
					t.Fatalf("expected muscle=Chest, got %v", m)
				}
				return []database.MongoExercise{
					{ID: primitive.NewObjectID(), Name: "Bench Press"},
				}, nil
			},
		},
	}}

	got, err := r.FilteredExercises(authContext(), model.ExerciseFilter{Muscle: &muscle})
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

func TestFilteredExercises_WithIDs(t *testing.T) {
	id1 := primitive.NewObjectID().Hex()
	id2 := primitive.NewObjectID().Hex()
	r := &resolvers.QueryResolver{&resolvers.Resolver{
		Exercises: &mockExerciseRepo{
			findFiltered: func(ctx context.Context, muscle *string, ids []*string) ([]database.MongoExercise, error) {
				if muscle != nil {
					t.Fatalf("expected nil muscle, got %v", *muscle)
				}
				if len(ids) != 2 || *ids[0] != id1 || *ids[1] != id2 {
					t.Fatalf("unexpected ids: %v", ids)
				}
				return []database.MongoExercise{
					{ID: primitive.NewObjectID(), Name: "Push Up"},
				}, nil
			},
		},
	}}

	got, err := r.FilteredExercises(authContext(), model.ExerciseFilter{
		Ids: []*string{&id1, &id2},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 exercise, got %d", len(got))
	}
	if got[0].Name != "Push Up" {
		t.Errorf("Name = %q, want %q", got[0].Name, "Push Up")
	}
}

func TestFilteredExercises_EmptyResult(t *testing.T) {
	r := &resolvers.QueryResolver{&resolvers.Resolver{
		Exercises: &mockExerciseRepo{
			findFiltered: func(ctx context.Context, muscle *string, ids []*string) ([]database.MongoExercise, error) {
				return []database.MongoExercise{}, nil
			},
		},
	}}

	got, err := r.FilteredExercises(authContext(), model.ExerciseFilter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 0 {
		t.Fatalf("expected 0 exercises, got %d", len(got))
	}
}

func TestFilteredExercises_RepoError(t *testing.T) {
	r := &resolvers.QueryResolver{&resolvers.Resolver{
		Exercises: &mockExerciseRepo{
			findFiltered: func(ctx context.Context, muscle *string, ids []*string) ([]database.MongoExercise, error) {
				return nil, errors.New("db error")
			},
		},
	}}

	_, err := r.FilteredExercises(authContext(), model.ExerciseFilter{})
	if err == nil || err.Error() != "db error" {
		t.Fatalf("expected db error, got %v", err)
	}
}
