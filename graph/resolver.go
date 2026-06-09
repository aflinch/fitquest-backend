package graph

import "fitness-backend/graph/model"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	workoutTemplates []*model.WorkoutTemplate
}

func NewResolver() *Resolver {
	return &Resolver{
		workoutTemplates: []*model.WorkoutTemplate{
			{
				ID:                 "template-strength-001",
				Title:              "Full Body Strength",
				Category:           "Strength",
				Difficulty:         "Beginner",
				Description:        "A simple starter workout for major muscle groups.",
				TargetMuscleGroups: []string{"Chest", "Back", "Legs", "Core"},
				Exercises: []*model.Exercise{
					{
						Name:              "Bodyweight Squat",
						EquipmentRequired: "None",
						PrimaryMuscle:     "Legs",
					},
					{
						Name:              "Push-Up",
						EquipmentRequired: "None",
						PrimaryMuscle:     "Chest",
					},
					{
						Name:              "Plank",
						EquipmentRequired: "None",
						PrimaryMuscle:     "Core",
					},
				},
				IsCommunityTemplate: false,
			},
			{
				ID:                 "template-cardio-001",
				Title:              "Quick Cardio Builder",
				Category:           "Cardio",
				Difficulty:         "Intermediate",
				Description:        "A short conditioning session for quick retrieval checks.",
				TargetMuscleGroups: []string{"Legs", "Cardio"},
				Exercises: []*model.Exercise{
					{
						Name:              "Jumping Jacks",
						EquipmentRequired: "None",
						PrimaryMuscle:     "Cardio",
					},
					{
						Name:              "Mountain Climbers",
						EquipmentRequired: "None",
						PrimaryMuscle:     "Core",
					},
				},
				IsCommunityTemplate: true,
			},
		},
	}
}
