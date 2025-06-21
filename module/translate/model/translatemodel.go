package model

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
