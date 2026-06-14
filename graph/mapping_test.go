package graph

import (
	"fitquest-backend/database"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMapMongoWorkoutToGQL(t *testing.T) {
	mw := database.MongoWorkout{
		ID:                  primitive.NewObjectID(),
		Title:               "Morning Pump",
		Category:            "Strength",
		Difficulty:          "Intermediate",
		Description:         "A solid morning routine",
		TargetMuscleGroups:  []string{"Chest", "Triceps"},
		IsCommunityTemplate: true,
		Exercises: []database.MongoExercise{
			{Name: "Push Up", EquipmentRequired: "None", PrimaryMuscle: "Chest"},
			{Name: "Dips", EquipmentRequired: "Parallel Bars", PrimaryMuscle: "Triceps"},
		},
	}

	got := mapMongoWorkoutToGQL(mw)

	if got.ID != mw.ID.Hex() {
		t.Errorf("ID = %q, want %q", got.ID, mw.ID.Hex())
	}
	if got.Title != mw.Title {
		t.Errorf("Title = %q, want %q", got.Title, mw.Title)
	}
	if got.Category != mw.Category {
		t.Errorf("Category = %q, want %q", got.Category, mw.Category)
	}
	if got.Difficulty != mw.Difficulty {
		t.Errorf("Difficulty = %q, want %q", got.Difficulty, mw.Difficulty)
	}
	if got.Description != mw.Description {
		t.Errorf("Description = %q, want %q", got.Description, mw.Description)
	}
	if len(got.TargetMuscleGroups) != len(mw.TargetMuscleGroups) {
		t.Errorf("TargetMuscleGroups length = %d, want %d", len(got.TargetMuscleGroups), len(mw.TargetMuscleGroups))
	}
	if got.IsCommunityTemplate != mw.IsCommunityTemplate {
		t.Errorf("IsCommunityTemplate = %v, want %v", got.IsCommunityTemplate, mw.IsCommunityTemplate)
	}
	if len(got.Exercises) != len(mw.Exercises) {
		t.Fatalf("Exercises length = %d, want %d", len(got.Exercises), len(mw.Exercises))
	}
	for i, e := range got.Exercises {
		if e.Name != mw.Exercises[i].Name {
			t.Errorf("Exercises[%d].Name = %q, want %q", i, e.Name, mw.Exercises[i].Name)
		}
		if e.EquipmentRequired != mw.Exercises[i].EquipmentRequired {
			t.Errorf("Exercises[%d].EquipmentRequired = %q, want %q", i, e.EquipmentRequired, mw.Exercises[i].EquipmentRequired)
		}
		if e.PrimaryMuscle != mw.Exercises[i].PrimaryMuscle {
			t.Errorf("Exercises[%d].PrimaryMuscle = %q, want %q", i, e.PrimaryMuscle, mw.Exercises[i].PrimaryMuscle)
		}
	}
}

func TestMapEmptyWorkout(t *testing.T) {
	mw := database.MongoWorkout{
		ID:    primitive.NewObjectID(),
		Title: "Empty Test",
	}

	got := mapMongoWorkoutToGQL(mw)

	if got.Title != "Empty Test" {
		t.Errorf("Title = %q, want %q", got.Title, "Empty Test")
	}
	if got.Exercises == nil {
		t.Error("Exercises should be empty slice, not nil")
	}
	if len(got.Exercises) != 0 {
		t.Errorf("Exercises length = %d, want 0", len(got.Exercises))
	}
}
