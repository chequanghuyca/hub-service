package model

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const SectionName = "sections"

// Validation error constants
var (
	ErrInvalidSection = errors.New("invalid section")
)

// Section represents a section of a challenge.
// @Description Section of a challenge containing title and content.
type Section struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title"`
	Content   string             `json:"content" bson:"content"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Image     string             `json:"image" bson:"image"`
}

func (Section) TableName() string {
	return SectionName
}

// SectionCreate is the model for creating a new section.
// @Description Required fields for creating a new section.
type SectionCreate struct {
	ID        primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title"`
	Content   string             `json:"content" bson:"content"`
	CreatedAt *time.Time         `json:"-" bson:"created_at"`
	UpdatedAt *time.Time         `json:"-" bson:"updated_at"`
	Image     string             `json:"image" bson:"image"`
}

func (SectionCreate) TableName() string {
	return SectionName
}

// SectionUpdate is the model for updating a section.
// @Description Optional fields for updating a section.
type SectionUpdate struct {
	Title     *string    `json:"title,omitempty" bson:"title,omitempty"`
	Content   *string    `json:"content,omitempty" bson:"content,omitempty"`
	UpdatedAt *time.Time `json:"-" bson:"updated_at,omitempty"`
	Image     *string    `json:"image,omitempty" bson:"image,omitempty"`
}

type SectionCreateResponse struct {
	Status string        `json:"status"`
	Data   SectionCreate `json:"data"`
}

type SectionUpdateResponse struct {
	Status string        `json:"status"`
	Data   SectionUpdate `json:"data"`
}

type SectionDeleteResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type PaginationMetadata struct {
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int64 `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

type SectionListResponse struct {
	Status   string             `json:"status"`
	Data     []Section          `json:"data"`
	Metadata PaginationMetadata `json:"metadata"`
}

type SectionGetDetailResponse struct {
	Status string  `json:"status"`
	Data   Section `json:"data"`
}

// SectionWithChallenges represents a section with its related challenges
type SectionWithChallenges struct {
	Section    Section           `json:"section"`
	Challenges []Challenge       `json:"challenges"`
	UserScore  *UserScoreSummary `json:"user_score,omitempty"`
}

// SectionWithScore represents a section with user score summary
type SectionWithScore struct {
	Section   Section           `json:"section"`
	UserScore *UserScoreSummary `json:"user_score,omitempty"`
}

// SectionSimple represents a simplified section with only id and title
// @Description Simplified section containing only id and title for list operations
type SectionSimple struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title string             `json:"title" bson:"title"`
}

// Challenge represents a challenge (imported from challenge module)
type Challenge struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title      string             `json:"title" bson:"title"`
	Content    string             `json:"content" bson:"content"`
	SourceLang string             `json:"source_lang" bson:"source_lang"`
	TargetLang string             `json:"target_lang" bson:"target_lang"`
	Difficulty string             `json:"difficulty" bson:"difficulty"`
	Category   string             `json:"category" bson:"category"`
	SectionID  primitive.ObjectID `json:"section_id" bson:"section_id"`
	CreatedAt  *time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt  *time.Time         `json:"updated_at" bson:"updated_at"`
}

// UserScoreSummary represents a summary of user's scores (imported from score module)
type UserScoreSummary struct {
	UserID          primitive.ObjectID `json:"user_id"`
	TotalScore      float64            `json:"total_score"`
	TotalChallenges int                `json:"total_challenges"`
	AverageScore    float64            `json:"average_score"`
	BestScore       float64            `json:"best_score"`
}
