package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	TranslationCollectionName = "translations"
	SentenceCollectionName    = "translation_sentences"
)

// Translation represents a complete text passage for translation
type Translation struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title      string             `json:"title" bson:"title"`
	Content    string             `json:"content" bson:"content"` // Full text passage
	SourceLang string             `json:"source_lang" bson:"source_lang"`
	TargetLang string             `json:"target_lang" bson:"target_lang"`
	Category   string             `json:"category" bson:"category"`
	Difficulty string             `json:"difficulty" bson:"difficulty"`
	TotalScore float64            `json:"total_score" bson:"total_score"` // Total possible score
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
	Image      string             `json:"image" bson:"image"`
}

func (Translation) TableName() string {
	return TranslationCollectionName
}

// TranslationCreate is the model for creating a new translation
type TranslationCreate struct {
	ID         primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Title      string             `json:"title" bson:"title" binding:"required"`
	Content    string             `json:"content" bson:"content" binding:"required"`
	SourceLang string             `json:"source_lang" bson:"source_lang" binding:"required"`
	TargetLang string             `json:"target_lang" bson:"target_lang" binding:"required"`
	Category   string             `json:"category" bson:"category"`
	Difficulty string             `json:"difficulty" bson:"difficulty" binding:"required,oneof=easy medium hard"`
	CreatedAt  *time.Time         `json:"-" bson:"created_at"`
	UpdatedAt  *time.Time         `json:"-" bson:"updated_at"`
	Image      string             `json:"image" bson:"image"`
}

func (TranslationCreate) TableName() string {
	return Translation{}.TableName()
}

// TranslationUpdate is the model for updating an existing translation
type TranslationUpdate struct {
	Title      *string    `json:"title,omitempty" bson:"title,omitempty"`
	Content    *string    `json:"content,omitempty" bson:"content,omitempty"`
	SourceLang *string    `json:"source_lang,omitempty" bson:"source_lang,omitempty"`
	TargetLang *string    `json:"target_lang,omitempty" bson:"target_lang,omitempty"`
	Category   *string    `json:"category,omitempty" bson:"category,omitempty"`
	Difficulty *string    `json:"difficulty,omitempty" bson:"difficulty,omitempty"`
	UpdatedAt  *time.Time `json:"-" bson:"updated_at,omitempty"`
	Image      *string    `json:"image,omitempty" bson:"image,omitempty"`
}

func (TranslationUpdate) TableName() string {
	return Translation{}.TableName()
}

// TranslationSentence represents a single sentence within a translation
type TranslationSentence struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TranslationID primitive.ObjectID `json:"translation_id" bson:"translation_id"`
	SentenceIndex int                `json:"sentence_index" bson:"sentence_index"`
	Content       string             `json:"content" bson:"content"`
	MaxScore      float64            `json:"max_score" bson:"max_score"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at" bson:"updated_at"`
}

func (TranslationSentence) TableName() string {
	return SentenceCollectionName
}

// TranslationSentenceCreate is the model for creating a new sentence
type TranslationSentenceCreate struct {
	ID            primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	TranslationID primitive.ObjectID `json:"translation_id" bson:"translation_id"`
	SentenceIndex int                `json:"sentence_index" bson:"sentence_index"`
	Content       string             `json:"content" bson:"content" binding:"required"`
	MaxScore      float64            `json:"max_score" bson:"max_score"`
	CreatedAt     *time.Time         `json:"-" bson:"created_at"`
	UpdatedAt     *time.Time         `json:"-" bson:"updated_at"`
}

func (TranslationSentenceCreate) TableName() string {
	return TranslationSentence{}.TableName()
}

// UserTranslationScore represents a user's score for a specific sentence
type UserTranslationScore struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID          primitive.ObjectID `json:"user_id" bson:"user_id"`
	TranslationID   primitive.ObjectID `json:"translation_id" bson:"translation_id"`
	SentenceID      primitive.ObjectID `json:"sentence_id" bson:"sentence_id"`
	SentenceIndex   int                `json:"sentence_index" bson:"sentence_index"`
	UserTranslation string             `json:"user_translation" bson:"user_translation"`
	Score           float64            `json:"score" bson:"score"`
	Feedback        string             `json:"feedback" bson:"feedback"`
	Errors          string             `json:"errors" bson:"errors"`
	Suggestions     string             `json:"suggestions" bson:"suggestions"`
	AttemptCount    int                `json:"attempt_count" bson:"attempt_count"`
	BestScore       float64            `json:"best_score" bson:"best_score"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
}

const UserTranslationScoreCollectionName = "user_translation_scores"

func (UserTranslationScore) TableName() string {
	return UserTranslationScoreCollectionName
}

// UserTranslationScoreCreate is the model for creating a new user score
type UserTranslationScoreCreate struct {
	UserID          primitive.ObjectID `json:"user_id" bson:"user_id"`
	TranslationID   primitive.ObjectID `json:"translation_id" bson:"translation_id"`
	SentenceID      primitive.ObjectID `json:"sentence_id" bson:"sentence_id"`
	SentenceIndex   int                `json:"sentence_index" bson:"sentence_index"`
	UserTranslation string             `json:"user_translation" bson:"user_translation"`
	Score           float64            `json:"score" bson:"score"`
	Feedback        string             `json:"feedback" bson:"feedback"`
	Errors          string             `json:"errors" bson:"errors"`
	Suggestions     string             `json:"suggestions" bson:"suggestions"`
	AttemptCount    int                `json:"attempt_count" bson:"attempt_count"`
	BestScore       float64            `json:"best_score" bson:"best_score"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
}

func (UserTranslationScoreCreate) TableName() string {
	return UserTranslationScore{}.TableName()
}

// TranslationWithSentences represents a translation with all its sentences
type TranslationWithSentences struct {
	Translation Translation           `json:"translation"`
	Sentences   []TranslationSentence `json:"sentences"`
}

// TranslationWithUserProgress represents a translation with user's progress
type TranslationWithUserProgress struct {
	Translation     Translation            `json:"translation"`
	Sentences       []TranslationSentence  `json:"sentences"`
	UserScores      []UserTranslationScore `json:"user_scores"`
	TotalUserScore  float64                `json:"total_user_score"`
	CompletedCount  int                    `json:"completed_count"`
	ProgressPercent float64                `json:"progress_percent"`
}

// SubmitSentenceTranslationRequest for submitting a sentence translation
type SubmitSentenceTranslationRequest struct {
	UserTranslation string `json:"user_translation" binding:"required"`
}

// SubmitSentenceTranslationResponse for sentence translation result
type SubmitSentenceTranslationResponse struct {
	Score           float64 `json:"score"`
	UserTranslation string  `json:"user_translation"`
	Feedback        string  `json:"feedback"`
	Errors          string  `json:"errors"`
	Suggestions     string  `json:"suggestions"`
	OriginalContent string  `json:"original_content"`
	AttemptCount    int     `json:"attempt_count"`
	BestScore       float64 `json:"best_score"`
	IsNewBest       bool    `json:"is_new_best"`
	TotalUserScore  float64 `json:"total_user_score"`
	ProgressPercent float64 `json:"progress_percent"`
}

// TranslationSummary represents a summary of user's translation progress
type TranslationSummary struct {
	TranslationID    primitive.ObjectID `json:"translation_id"`
	Title            string             `json:"title"`
	TotalSentences   int                `json:"total_sentences"`
	CompletedCount   int                `json:"completed_count"`
	TotalUserScore   float64            `json:"total_user_score"`
	MaxPossibleScore float64            `json:"max_possible_score"`
	ProgressPercent  float64            `json:"progress_percent"`
	LastAttemptAt    *time.Time         `json:"last_attempt_at"`
}

// GetUserTranslationScoresRequest for getting user's translation scores
type GetUserTranslationScoresRequest struct {
	UserID string `uri:"user_id" binding:"required"`
}

// GetUserTranslationScoresResponse for user's translation scores
type GetUserTranslationScoresResponse struct {
	Summaries []TranslationSummary `json:"summaries"`
}

// GetUserTranslationScoresAPIResponse API response wrapper
type GetUserTranslationScoresAPIResponse struct {
	Data GetUserTranslationScoresResponse `json:"data"`
}
