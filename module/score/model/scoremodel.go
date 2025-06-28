package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CollectionName = "scores"

// Score represents a user's score for a specific challenge
type Score struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID           primitive.ObjectID `json:"user_id" bson:"user_id"`
	ChallengeID      primitive.ObjectID `json:"challenge_id" bson:"challenge_id"`
	UserTranslation  string             `json:"user_translation" bson:"user_translation"`
	DeepLScore       float64            `json:"deepl_score" bson:"deepl_score"`
	DeepLTranslation string             `json:"deepl_translation" bson:"deepl_translation"`
	OriginalContent  string             `json:"original_content" bson:"original_content"`
	AttemptCount     int                `json:"attempt_count" bson:"attempt_count"`
	BestScore        float64            `json:"best_score" bson:"best_score"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at"`
}

func (Score) TableName() string {
	return CollectionName
}

// ScoreCreate is the model for creating a new score
type ScoreCreate struct {
	UserID           primitive.ObjectID `json:"user_id" bson:"user_id"`
	ChallengeID      primitive.ObjectID `json:"challenge_id" bson:"challenge_id"`
	UserTranslation  string             `json:"user_translation" bson:"user_translation"`
	DeepLScore       float64            `json:"deepl_score" bson:"deepl_score"`
	DeepLTranslation string             `json:"deepl_translation" bson:"deepl_translation"`
	OriginalContent  string             `json:"original_content" bson:"original_content"`
	AttemptCount     int                `json:"attempt_count" bson:"attempt_count"`
	BestScore        float64            `json:"best_score" bson:"best_score"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at"`
}

func (ScoreCreate) TableName() string {
	return Score{}.TableName()
}

// ScoreUpdate is the model for updating an existing score
type ScoreUpdate struct {
	UserTranslation  *string    `json:"user_translation,omitempty" bson:"user_translation,omitempty"`
	DeepLScore       *float64   `json:"deepl_score,omitempty" bson:"deepl_score,omitempty"`
	DeepLTranslation *string    `json:"deepl_translation,omitempty" bson:"deepl_translation,omitempty"`
	AttemptCount     *int       `json:"attempt_count,omitempty" bson:"attempt_count,omitempty"`
	BestScore        *float64   `json:"best_score,omitempty" bson:"best_score,omitempty"`
	UpdatedAt        *time.Time `json:"updated_at" bson:"updated_at,omitempty"`
}

func (ScoreUpdate) TableName() string {
	return Score{}.TableName()
}

// UserScoreSummary represents a summary of user's scores
type UserScoreSummary struct {
	UserID          primitive.ObjectID `json:"user_id"`
	TotalScore      float64            `json:"total_score"`
	TotalChallenges int                `json:"total_challenges"`
	AverageScore    float64            `json:"average_score"`
	BestScore       float64            `json:"best_score"`
}

// ChallengeScore represents a user's score for a specific challenge
type ChallengeScore struct {
	ChallengeID      primitive.ObjectID `json:"challenge_id"`
	ChallengeTitle   string             `json:"challenge_title"`
	BestScore        float64            `json:"best_score"`
	AttemptCount     int                `json:"attempt_count"`
	LastAttemptAt    time.Time          `json:"last_attempt_at"`
	UserTranslation  string             `json:"user_translation"`
	DeepLTranslation string             `json:"deepl_translation"`
	OriginalContent  string             `json:"original_content"`
}

// SubmitScoreRequest represents the request for submitting a score
type SubmitScoreRequest struct {
	ChallengeID     string `json:"challenge_id" binding:"required"`
	UserTranslation string `json:"user_translation" binding:"required"`
}

// SubmitScoreResponse represents the response for submitting a score
type SubmitScoreResponse struct {
	Score            float64 `json:"score"`
	UserTranslation  string  `json:"user_translation"`
	DeepLTranslation string  `json:"deepl_translation"`
	OriginalContent  string  `json:"original_content"`
	AttemptCount     int     `json:"attempt_count"`
	BestScore        float64 `json:"best_score"`
	IsNewBest        bool    `json:"is_new_best"`
}

// GetUserScoresRequest represents the request for getting user scores
type GetUserScoresRequest struct {
	UserID string `uri:"user_id" binding:"required"`
}

// GetUserScoresResponse represents the response for getting user scores
type GetUserScoresResponse struct {
	Summary UserScoreSummary `json:"summary"`
	Scores  []ChallengeScore `json:"scores"`
}

// API Response Models for Swagger
type GetUserScoresAPIResponse struct {
	Status string                `json:"status"`
	Data   GetUserScoresResponse `json:"data"`
}

// Translation Models (moved from translate module)

// Challenge represents a text to be translated, split into sentences.
// @Description Translation challenge containing multiple sentences to translate
type Challenge struct {
	ID        string   `json:"id" example:"challenge_1"`                                    // Unique identifier for the challenge
	Title     string   `json:"title" example:"Trích đoạn 'Tôi thấy hoa vàng trên cỏ xanh'"` // Title of the challenge
	Sentences []string `json:"sentences"`                                                   // Array of sentences to translate
}

// ScoreRequest is the user's translation submission.
// @Description Request body for scoring user's translation
type ScoreRequest struct {
	ChallengeID     string `json:"challenge_id" binding:"required" example:"challenge_1"`
	SentenceIndex   int    `json:"sentence_index" binding:"min=0" example:"0"`
	UserTranslation string `json:"user_translation" binding:"required" example:"I see yellow flowers on the green grass."`
}

// ScoreResponse contains the scoring result.
// @Description Response containing the scoring result and comparison data
type ScoreResponse struct {
	Score            float64 `json:"score" example:"95.23"`                                               // Similarity score (0-100)
	UserTranslation  string  `json:"user_translation" example:"I see yellow flowers on the green grass."` // User's submitted translation
	DeepLTranslation string  `json:"deepl_translation" example:"I see yellow flowers on green grass."`    // DeepL's reference translation
	OriginalSentence string  `json:"original_sentence" example:"Tôi thấy hoa vàng trên cỏ xanh."`         // Original sentence in Vietnamese
}

// APIResponse wraps the score response with a status.
// @Description Standard API response wrapper
type APIResponse struct {
	Status string        `json:"status" example:"success"` // Response status
	Data   ScoreResponse `json:"data"`                     // Response data
}
