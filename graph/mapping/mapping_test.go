package mapping

import (
	"fitquest-backend/database"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMapMongoUserToGQL_NotFound(t *testing.T) {
	got := MapMongoUserToGQL(database.MongoUser{})

	if got == nil {
		t.Fatal("expected non-nil result")
	}
	if got.ID != primitive.NilObjectID.Hex() {
		t.Errorf("expected zero ObjectID, got %q", got.ID)
	}
	if got.Username != "" {
		t.Errorf("expected empty Username, got %q", got.Username)
	}
}

func TestMapMongoUserToGQL(t *testing.T) {
	mu := database.MongoUser{
		ID:           primitive.NewObjectID(),
		Username:     "testuser",
		PasswordHash: "hash",
	}

	got := MapMongoUserToGQL(mu)

	if got.ID != mu.ID.Hex() {
		t.Errorf("ID = %q, want %q", got.ID, mu.ID.Hex())
	}
	if got.Username != mu.Username {
		t.Errorf("Username = %q, want %q", got.Username, mu.Username)
	}
}

func TestMapMongoExerciseToGQL_Empty(t *testing.T) {
	got := MapMongoExerciseToGQL(database.MongoExercise{})

	if got == nil {
		t.Fatal("expected non-nil result")
	}
	if got.ID != primitive.NilObjectID.Hex() {
		t.Errorf("expected zero ObjectID, got %q", got.ID)
	}
	if got.Name != "" {
		t.Errorf("expected empty Name, got %q", got.Name)
	}
	if got.Category == nil || *got.Category != "" {
		t.Errorf("expected empty Category, got %v", got.Category)
	}
	if got.Mechanic == nil || *got.Mechanic != "" {
		t.Errorf("expected empty Mechanic, got %v", got.Mechanic)
	}
	if len(got.PrimaryMuscles) != 0 {
		t.Errorf("expected empty PrimaryMuscles, got %v", got.PrimaryMuscles)
	}
	if len(got.SecondaryMuscles) != 0 {
		t.Errorf("expected empty SecondaryMuscles, got %v", got.SecondaryMuscles)
	}
}

func TestMapMongoExerciseToGQL(t *testing.T) {
	me := database.MongoExercise{
		ID:               primitive.NewObjectID(),
		Name:             "Bench Press",
		Category:         "Strength",
		Mechanic:         "Compound",
		PrimaryMuscles:   []string{"Chest", "Triceps"},
		SecondaryMuscles: []string{"Front Delts"},
	}

	got := MapMongoExerciseToGQL(me)

	if got.ID != me.ID.Hex() {
		t.Errorf("ID = %q, want %q", got.ID, me.ID.Hex())
	}
	if got.Name != me.Name {
		t.Errorf("Name = %q, want %q", got.Name, me.Name)
	}
	if *got.Category != me.Category {
		t.Errorf("Category = %q, want %q", *got.Category, me.Category)
	}
	if *got.Mechanic != me.Mechanic {
		t.Errorf("Mechanic = %q, want %q", *got.Mechanic, me.Mechanic)
	}
	if len(got.PrimaryMuscles) != len(me.PrimaryMuscles) {
		t.Fatalf("PrimaryMuscles length = %d, want %d", len(got.PrimaryMuscles), len(me.PrimaryMuscles))
	}
	for i, p := range got.PrimaryMuscles {
		if *p != me.PrimaryMuscles[i] {
			t.Errorf("PrimaryMuscles[%d] = %q, want %q", i, *p, me.PrimaryMuscles[i])
		}
	}
	if len(got.SecondaryMuscles) != len(me.SecondaryMuscles) {
		t.Fatalf("SecondaryMuscles length = %d, want %d", len(got.SecondaryMuscles), len(me.SecondaryMuscles))
	}
	for i, s := range got.SecondaryMuscles {
		if *s != me.SecondaryMuscles[i] {
			t.Errorf("SecondaryMuscles[%d] = %q, want %q", i, *s, me.SecondaryMuscles[i])
		}
	}
}
