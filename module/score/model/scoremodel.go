package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CollectionName = "scores"

// Score represents a user's score for a specific challenge
// Đã chuyển sang dùng Gemini hoàn toàn
// Lưu ý: GeminiErrors và GeminiSuggestions có thể là JSON/text tuỳ theo cách lưu trữ
// Nếu cần lưu dạng []string hoặc []struct, hãy điều chỉnh lại kiểu dữ liệu

type Score struct {
	ID                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID            primitive.ObjectID `json:"user_id" bson:"user_id"`
	ChallengeID       primitive.ObjectID `json:"challenge_id" bson:"challenge_id"`
	UserTranslation   string             `json:"user_translation" bson:"user_translation"`
	GeminiScore       float64            `json:"gemini_score" bson:"gemini_score"`
	GeminiFeedback    string             `json:"gemini_feedback" bson:"gemini_feedback"`
	GeminiErrors      string             `json:"gemini_errors" bson:"gemini_errors"`
	GeminiSuggestions string             `json:"gemini_suggestions" bson:"gemini_suggestions"`
	OriginalContent   string             `json:"original_content" bson:"original_content"`
	AttemptCount      int                `json:"attempt_count" bson:"attempt_count"`
	BestScore         float64            `json:"best_score" bson:"best_score"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at" bson:"updated_at"`
}

func (Score) TableName() string {
	return CollectionName
}

// ScoreCreate is the model for creating a new score
// Đã chuyển sang dùng Gemini hoàn toàn
// GeminiErrors, GeminiSuggestions có thể lưu dạng JSON string hoặc []byte

type ScoreCreate struct {
	UserID            primitive.ObjectID `json:"user_id" bson:"user_id"`
	ChallengeID       primitive.ObjectID `json:"challenge_id" bson:"challenge_id"`
	UserTranslation   string             `json:"user_translation" bson:"user_translation"`
	GeminiScore       float64            `json:"gemini_score" bson:"gemini_score"`
	GeminiFeedback    string             `json:"gemini_feedback" bson:"gemini_feedback"`
	GeminiErrors      string             `json:"gemini_errors" bson:"gemini_errors"`
	GeminiSuggestions string             `json:"gemini_suggestions" bson:"gemini_suggestions"`
	OriginalContent   string             `json:"original_content" bson:"original_content"`
	AttemptCount      int                `json:"attempt_count" bson:"attempt_count"`
	BestScore         float64            `json:"best_score" bson:"best_score"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at" bson:"updated_at"`
}

func (ScoreCreate) TableName() string {
	return Score{}.TableName()
}

// ScoreUpdate is the model for updating an existing score

type ScoreUpdate struct {
	UserTranslation   *string    `json:"user_translation,omitempty" bson:"user_translation,omitempty"`
	GeminiScore       *float64   `json:"gemini_score,omitempty" bson:"gemini_score,omitempty"`
	GeminiFeedback    *string    `json:"gemini_feedback,omitempty" bson:"gemini_feedback,omitempty"`
	GeminiErrors      *string    `json:"gemini_errors,omitempty" bson:"gemini_errors,omitempty"`
	GeminiSuggestions *string    `json:"gemini_suggestions,omitempty" bson:"gemini_suggestions,omitempty"`
	AttemptCount      *int       `json:"attempt_count,omitempty" bson:"attempt_count,omitempty"`
	BestScore         *float64   `json:"best_score,omitempty" bson:"best_score,omitempty"`
	UpdatedAt         *time.Time `json:"updated_at" bson:"updated_at,omitempty"`
}

func (ScoreUpdate) TableName() string {
	return Score{}.TableName()
}

// ChallengeScore đại diện cho điểm số của user cho một challenge
// Đã chuyển sang dùng Gemini hoàn toàn

type ChallengeScore struct {
	ChallengeID       primitive.ObjectID `json:"challenge_id"`
	ChallengeTitle    string             `json:"challenge_title"`
	BestScore         float64            `json:"best_score"`
	AttemptCount      int                `json:"attempt_count"`
	LastAttemptAt     time.Time          `json:"last_attempt_at"`
	UserTranslation   string             `json:"user_translation"`
	GeminiFeedback    string             `json:"gemini_feedback"`
	GeminiErrors      string             `json:"gemini_errors"`
	GeminiSuggestions string             `json:"gemini_suggestions"`
	OriginalContent   string             `json:"original_content"`
}

// SubmitScoreRequest giữ nguyên

type SubmitScoreRequest struct {
	ChallengeID     string `json:"challenge_id" binding:"required"`
	UserTranslation string `json:"user_translation" binding:"required"`
}

// SubmitScoreResponse trả về kết quả chấm điểm của Gemini

type SubmitScoreResponse struct {
	Score             float64 `json:"score"`
	UserTranslation   string  `json:"user_translation"`
	GeminiFeedback    string  `json:"gemini_feedback"`
	GeminiErrors      string  `json:"gemini_errors"`
	GeminiSuggestions string  `json:"gemini_suggestions"`
	OriginalContent   string  `json:"original_content"`
	AttemptCount      int     `json:"attempt_count"`
	BestScore         float64 `json:"best_score"`
	IsNewBest         bool    `json:"is_new_best"`
}

// ScoreRequest giữ nguyên

type ScoreRequest struct {
	ChallengeID     string `json:"challenge_id" binding:"required"`
	SentenceIndex   int    `json:"sentence_index"`
	UserTranslation string `json:"user_translation" binding:"required"`
}

// ScoreResponse trả về kết quả chấm điểm của Gemini

type ScoreResponse struct {
	Score             float64 `json:"score"`
	UserTranslation   string  `json:"user_translation"`
	GeminiFeedback    string  `json:"gemini_feedback"`
	GeminiErrors      string  `json:"gemini_errors"`
	GeminiSuggestions string  `json:"gemini_suggestions"`
	OriginalSentence  string  `json:"original_sentence"`
}

// Các struct khác giữ nguyên

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
