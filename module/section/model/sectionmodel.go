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
