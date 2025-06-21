package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CollectionName = "challenges"

// Challenge represents a translation challenge stored in the database.
// @Description Contains the details of a translation challenge.
type Challenge struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty" example:"62b4c3789196e8a159933552"`
	Title      string             `json:"title" bson:"title" example:"Greetings"`
	Content    string             `json:"content" bson:"content" example:"Hello, world!"`
	SourceLang string             `json:"source_lang" bson:"source_lang" example:"EN"`
	TargetLang string             `json:"target_lang" bson:"target_lang" example:"VI"`
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
	SourceLang string             `json:"source_lang" bson:"source_lang" binding:"required" example:"EN"`
	TargetLang string             `json:"target_lang" bson:"target_lang" binding:"required" example:"VI"`
	CreatedAt  *time.Time         `json:"-" bson:"created_at"`
	UpdatedAt  *time.Time         `json:"-" bson:"updated_at"`
}

func (ChallengeCreate) TableName() string {
	return Challenge{}.TableName()
}

// ChallengeUpdate is the model for updating an existing challenge.
// @Description Fields available for updating a translation challenge. All fields are optional.
type ChallengeUpdate struct {
	Title      *string    `json:"title,omitempty" bson:"title,omitempty" example:"Formal Greetings"`
	Content    *string    `json:"content,omitempty" bson:"content,omitempty" example:"Good morning, everyone."`
	SourceLang *string    `json:"source_lang,omitempty" bson:"source_lang,omitempty" example:"EN-US"`
	TargetLang *string    `json:"target_lang,omitempty" bson:"target_lang,omitempty" example:"DE"`
	UpdatedAt  *time.Time `json:"-" bson:"updated_at,omitempty"`
}

func (ChallengeUpdate) TableName() string {
	return Challenge{}.TableName()
}
