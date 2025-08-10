package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CollectionName = "scores"

// Score represents a user's score for a specific challenge
// Updated to match the new Gemini-based scoring structure

type Score struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID          primitive.ObjectID `json:"user_id" bson:"user_id"`
	ChallengeID     primitive.ObjectID `json:"challenge_id" bson:"challenge_id"`
	UserTranslation string             `json:"user_translation" bson:"user_translation"`
	Score           float64            `json:"score" bson:"score"`
	Feedback        string             `json:"feedback" bson:"feedback"`
	Errors          string             `json:"errors" bson:"errors"`
	Suggestions     string             `json:"suggestions" bson:"suggestions"`
	OriginalContent string             `json:"original_content" bson:"original_content"`
	AttemptCount    int                `json:"attempt_count" bson:"attempt_count"`
	BestScore       float64            `json:"best_score" bson:"best_score"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
}

func (Score) TableName() string {
	return CollectionName
}

// ScoreCreate is the model for creating a new score
// Updated to match the new Gemini-based scoring structure

type ScoreCreate struct {
	UserID          primitive.ObjectID `json:"user_id" bson:"user_id"`
	ChallengeID     primitive.ObjectID `json:"challenge_id" bson:"challenge_id"`
	UserTranslation string             `json:"user_translation" bson:"user_translation"`
	Score           float64            `json:"score" bson:"score"`
	Feedback        string             `json:"feedback" bson:"feedback"`
	Errors          string             `json:"errors" bson:"errors"`
	Suggestions     string             `json:"suggestions" bson:"suggestions"`
	OriginalContent string             `json:"original_content" bson:"original_content"`
	AttemptCount    int                `json:"attempt_count" bson:"attempt_count"`
	BestScore       float64            `json:"best_score" bson:"best_score"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
}

func (ScoreCreate) TableName() string {
	return Score{}.TableName()
}

// ScoreUpdate is the model for updating an existing score
// Updated to match the new Gemini-based scoring structure

type ScoreUpdate struct {
	UserTranslation *string    `json:"user_translation,omitempty" bson:"user_translation,omitempty"`
	Score           *float64   `json:"score,omitempty" bson:"score,omitempty"`
	Feedback        *string    `json:"feedback,omitempty" bson:"feedback,omitempty"`
	Errors          *string    `json:"errors,omitempty" bson:"errors,omitempty"`
	Suggestions     *string    `json:"suggestions,omitempty" bson:"suggestions,omitempty"`
	AttemptCount    *int       `json:"attempt_count,omitempty" bson:"attempt_count,omitempty"`
	BestScore       *float64   `json:"best_score,omitempty" bson:"best_score,omitempty"`
	UpdatedAt       *time.Time `json:"updated_at" bson:"updated_at,omitempty"`
}

func (ScoreUpdate) TableName() string {
	return Score{}.TableName()
}

// ChallengeScore đại diện cho điểm số của user cho một challenge
// Updated to match the new Gemini-based scoring structure

type ChallengeScore struct {
	ChallengeID     primitive.ObjectID `json:"challenge_id"`
	ChallengeTitle  string             `json:"challenge_title"`
	BestScore       float64            `json:"best_score"`
	AttemptCount    int                `json:"attempt_count"`
	LastAttemptAt   time.Time          `json:"last_attempt_at"`
	UserTranslation string             `json:"user_translation"`
	Feedback        string             `json:"feedback"`
	Errors          string             `json:"errors"`
	Suggestions     string             `json:"suggestions"`
	OriginalContent string             `json:"original_content"`
}

// SubmitScoreRequest giữ nguyên

type SubmitScoreRequest struct {
	ChallengeID     string `json:"challenge_id" binding:"required"`
	UserTranslation string `json:"user_translation" binding:"required"`
}

// SubmitScoreResponse trả về kết quả chấm điểm của Gemini
// Updated to match the new Gemini-based scoring structure

type SubmitScoreResponse struct {
	Score           float64 `json:"score"`
	UserTranslation string  `json:"user_translation"`
	Feedback        string  `json:"feedback"`
	Errors          string  `json:"errors"`
	Suggestions     string  `json:"suggestions"`
	OriginalContent string  `json:"original_content"`
	AttemptCount    int     `json:"attempt_count"`
	BestScore       float64 `json:"best_score"`
	IsNewBest       bool    `json:"is_new_best"`
}

// ScoreRequest giữ nguyên

type ScoreRequest struct {
	ChallengeID     string `json:"challenge_id" binding:"required"`
	SentenceIndex   int    `json:"sentence_index"`
	UserTranslation string `json:"user_translation" binding:"required"`
}

// ScoreResponse trả về kết quả chấm điểm của Gemini
// Updated to match the new Gemini-based scoring structure

type ScoreResponse struct {
	Score            float64 `json:"score"`
	UserTranslation  string  `json:"user_translation"`
	Feedback         string  `json:"feedback"`
	Errors           string  `json:"errors"`
	Suggestions      string  `json:"suggestions"`
	OriginalSentence string  `json:"original_sentence"`
}

type UserScoreSummary struct {
	UserID          primitive.ObjectID `json:"user_id"`
	TotalScore      float64            `json:"total_score"`
	TotalChallenges int                `json:"total_challenges"`
	AverageScore    float64            `json:"average_score"`
	BestScore       float64            `json:"best_score"`
}

type GetUserScoresRequest struct {
	UserID string `uri:"user_id" binding:"required"`
}

type GetUserScoresResponse struct {
	Summary UserScoreSummary `json:"summary"`
	Scores  []ChallengeScore `json:"scores"`
}

type GetUserScoresAPIResponse struct {
	Data GetUserScoresResponse `json:"data"`
}
