package model

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CollectionName = "challenges"

// Validation constants
const (
	DifficultyEasy   = "easy"
	DifficultyMedium = "medium"
	DifficultyHard   = "hard"

	CategoryWork          = "work"
	CategoryLife          = "life"
	CategoryTravel        = "travel"
	CategoryDailyLife     = "daily_life"
	CategoryEntertainment = "entertainment"
	CategoryEducation     = "education"
	CategoryEconomy       = "economy"
	CategoryHealth        = "health"
	CategorySport         = "sport"
	CategoryTechnology    = "technology"
	CategoryCulture       = "culture"
)

// Validation error constants
var (
	ErrInvalidDifficulty = errors.New("invalid difficulty value")
	ErrInvalidCategory   = errors.New("invalid category value")
)

// Challenge represents a translation challenge stored in the database.
// @Description Contains the details of a translation challenge.
type Challenge struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty" example:"62b4c3789196e8a159933552"`
	Title      string             `json:"title" bson:"title" example:"Greetings"`
	Content    string             `json:"content" bson:"content" example:"Hello, world!"`
	SourceLang string             `json:"source_lang" bson:"source_lang" example:"VI"`
	TargetLang string             `json:"target_lang" bson:"target_lang" example:"EN"`
	Difficulty string             `json:"difficulty" bson:"difficulty" example:"easy"`
	Category   string             `json:"category" bson:"category" example:"work"`
	SectionID  primitive.ObjectID `json:"section_id" bson:"section_id" example:"62b4c3789196e8a159933552"`
	CreatedAt  *time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt  *time.Time         `json:"updated_at" bson:"updated_at"`
}

func (Challenge) TableName() string {
	return CollectionName
}

// ChallengeCreate is the model for creating a new challenge.
// @Description Required fields for creating a new translation challenge.
type ChallengeCreate struct {
	ID         primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Title      string             `json:"title" bson:"title" binding:"required" example:"Greetings"`
	Content    string             `json:"content" bson:"content" binding:"required" example:"Hello, world!"`
	SourceLang string             `json:"source_lang" bson:"source_lang" binding:"required" example:"VI"`
	TargetLang string             `json:"target_lang" bson:"target_lang" binding:"required" example:"EN"`
	Difficulty string             `json:"difficulty" bson:"difficulty" binding:"required,oneof=easy medium hard" example:"easy"`
	Category   string             `json:"category" bson:"category" binding:"omitempty" example:"work"`
	SectionID  primitive.ObjectID `json:"section_id" bson:"section_id" example:"62b4c3789196e8a159933552"`
	CreatedAt  *time.Time         `json:"-" bson:"created_at"`
	UpdatedAt  *time.Time         `json:"-" bson:"updated_at"`
}

func (ChallengeCreate) TableName() string {
	return Challenge{}.TableName()
}

// ChallengeUpdate is the model for updating an existing challenge.
// @Description Fields available for updating a translation challenge. All fields are optional.
type ChallengeUpdate struct {
	Title      *string             `json:"title,omitempty" bson:"title,omitempty" binding:"omitempty,min=1" example:"Formal Greetings"`
	Content    *string             `json:"content,omitempty" bson:"content,omitempty" binding:"omitempty,min=1" example:"Good morning, everyone."`
	SourceLang *string             `json:"source_lang,omitempty" bson:"source_lang,omitempty" binding:"omitempty,min=2,max=2" example:"VI"`
	TargetLang *string             `json:"target_lang,omitempty" bson:"target_lang,omitempty" binding:"omitempty,min=2,max=2" example:"EN"`
	Difficulty *string             `json:"difficulty,omitempty" bson:"difficulty,omitempty" binding:"omitempty,oneof=easy medium hard" example:"easy"`
	Category   *string             `json:"category,omitempty" bson:"category,omitempty" example:"work"`
	SectionID  *primitive.ObjectID `json:"section_id,omitempty" bson:"section_id,omitempty" example:"62b4c3789196e8a159933552"`
	UpdatedAt  *time.Time          `json:"-" bson:"updated_at,omitempty"`
}

func (ChallengeUpdate) TableName() string {
	return Challenge{}.TableName()
}

// HasUpdates returns true if at least one field is provided for update
func (cu ChallengeUpdate) HasUpdates() bool {
	return cu.Title != nil || cu.Content != nil || cu.SourceLang != nil ||
		cu.TargetLang != nil || cu.Difficulty != nil || cu.Category != nil
}

// Helper functions to generate validation strings
func GetDifficultyValidation() string {
	return "oneof=" + DifficultyEasy + " " + DifficultyMedium + " " + DifficultyHard
}

func GetCategoryValidation() string {
	categories := GetValidCategories()
	result := "oneof="
	for i, category := range categories {
		if i > 0 {
			result += " "
		}
		result += category
	}
	return result
}

// GetValidDifficulties returns all valid difficulty values
func GetValidDifficulties() []string {
	return []string{DifficultyEasy, DifficultyMedium, DifficultyHard}
}

// GetValidCategories returns all valid category values
func GetValidCategories() []string {
	return []string{
		CategoryWork, CategoryLife, CategoryTravel, CategoryDailyLife,
		CategoryEntertainment, CategoryEducation, CategoryEconomy,
		CategoryHealth, CategorySport, CategoryTechnology, CategoryCulture,
	}
}
