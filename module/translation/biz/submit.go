package biz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	common "hub-service/common"
	"hub-service/module/translation/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubmitTranslationStore interface {
	GetTranslation(ctx context.Context, id primitive.ObjectID) (*model.Translation, error)
	GetSentencesByTranslationID(ctx context.Context, translationID primitive.ObjectID) ([]model.TranslationSentence, error)
	GetUserScore(ctx context.Context, userID, translationID primitive.ObjectID, sentenceIndex int) (*model.UserTranslationScore, error)
	CreateUserScore(ctx context.Context, data *model.UserTranslationScoreCreate) error
	UpdateUserScore(ctx context.Context, id primitive.ObjectID, data *model.UserTranslationScoreCreate) error
	GetUserScoresByTranslation(ctx context.Context, userID, translationID primitive.ObjectID) ([]model.UserTranslationScore, error)
}

type submitTranslationBiz struct {
	store   SubmitTranslationStore
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewSubmitTranslationBiz(store SubmitTranslationStore, apiKey, baseURL string) *submitTranslationBiz {
	return &submitTranslationBiz{
		store:   store,
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (biz *submitTranslationBiz) SubmitSentenceTranslation(ctx context.Context, translationID primitive.ObjectID, sentenceIndex int, userTranslation string, userID primitive.ObjectID) (*model.SubmitSentenceTranslationResponse, error) {
	// Get translation and sentence
	translation, err := biz.store.GetTranslation(ctx, translationID)
	if err != nil {
		return nil, err
	}

	sentences, err := biz.store.GetSentencesByTranslationID(ctx, translationID)
	if err != nil {
		return nil, err
	}

	if sentenceIndex >= len(sentences) {
		return nil, common.ErrInvalidRequest(errors.New("sentence index out of range"))
	}

	sentence := sentences[sentenceIndex]

	// Get existing user score
	existingScore, _ := biz.store.GetUserScore(ctx, userID, translationID, sentenceIndex)

	// Calculate score using AI
	score, feedback, errors, suggestions, err := biz.calculateScore(sentence.Content, userTranslation, translation.TargetLang)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	attemptCount := 1
	bestScore := score
	isNewBest := true

	if existingScore != nil {
		// Update existing score
		attemptCount = existingScore.AttemptCount + 1
		bestScore = existingScore.BestScore
		isNewBest = false

		if score > existingScore.BestScore {
			bestScore = score
			isNewBest = true
		}

		scoreData := &model.UserTranslationScoreCreate{
			UserID:          userID,
			TranslationID:   translationID,
			SentenceID:      sentence.ID,
			SentenceIndex:   sentenceIndex,
			UserTranslation: userTranslation,
			Score:           score,
			Feedback:        feedback,
			Errors:          errors,
			Suggestions:     suggestions,
			AttemptCount:    attemptCount,
			BestScore:       bestScore,
			CreatedAt:       existingScore.CreatedAt,
			UpdatedAt:       now,
		}

		err = biz.store.UpdateUserScore(ctx, existingScore.ID, scoreData)
	} else {
		// Create new score
		scoreData := &model.UserTranslationScoreCreate{
			UserID:          userID,
			TranslationID:   translationID,
			SentenceID:      sentence.ID,
			SentenceIndex:   sentenceIndex,
			UserTranslation: userTranslation,
			Score:           score,
			Feedback:        feedback,
			Errors:          errors,
			Suggestions:     suggestions,
			AttemptCount:    attemptCount,
			BestScore:       bestScore,
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		err = biz.store.CreateUserScore(ctx, scoreData)
	}

	if err != nil {
		return nil, err
	}

	// Calculate total user score and progress
	userScores, err := biz.store.GetUserScoresByTranslation(ctx, userID, translationID)
	if err != nil {
		return nil, err
	}

	totalUserScore := 0.0
	for _, userScore := range userScores {
		totalUserScore += userScore.BestScore
	}

	// Calculate total possible score from all sentences
	totalPossibleScore := 0.0
	for _, s := range sentences {
		totalPossibleScore += s.MaxScore
	}

	progressPercent := 0.0
	if totalPossibleScore > 0 {
		progressPercent = (totalUserScore / totalPossibleScore) * 100
	}

	return &model.SubmitSentenceTranslationResponse{
		Score:           score,
		UserTranslation: userTranslation,
		Feedback:        feedback,
		Errors:          errors,
		Suggestions:     suggestions,
		OriginalContent: sentence.Content,
		AttemptCount:    attemptCount,
		BestScore:       bestScore,
		IsNewBest:       isNewBest,
		TotalUserScore:  totalUserScore,
		ProgressPercent: progressPercent,
	}, nil
}

// Helper function to calculate score using AI
func (biz *submitTranslationBiz) calculateScore(original, translation, targetLanguage string) (float64, string, string, string, error) {
	promptText := fmt.Sprintf(GeminiGrammarPrompt, original, translation, targetLanguage)

	req := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: promptText},
				},
			},
		},
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return 0, "", "", "", errors.New("failed to marshal request")
	}

	httpReq, err := http.NewRequestWithContext(context.Background(), "POST", biz.baseURL, strings.NewReader(string(jsonData)))
	if err != nil {
		return 0, "", "", "", errors.New("failed to create request")
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-goog-api-key", biz.apiKey)

	resp, err := biz.client.Do(httpReq)
	if err != nil {
		return 0, "", "", "", errors.New("failed to execute request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, "", "", "", fmt.Errorf("status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var geminiResp GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return 0, "", "", "", errors.New("failed to decode response")
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return 0, "", "", "", errors.New("no response from Gemini")
	}

	responseText := geminiResp.Candidates[0].Content.Parts[0].Text
	responseText = strings.TrimSpace(responseText)
	if strings.HasPrefix(responseText, "```json") {
		responseText = strings.TrimPrefix(responseText, "```json")
		responseText = strings.TrimSpace(responseText)
	}
	if strings.HasPrefix(responseText, "```") {
		responseText = strings.TrimPrefix(responseText, "```")
		responseText = strings.TrimSpace(responseText)
	}
	if strings.HasSuffix(responseText, "```") {
		responseText = strings.TrimSuffix(responseText, "```")
		responseText = strings.TrimSpace(responseText)
	}

	var analysis GrammarAnalysis
	if err := json.Unmarshal([]byte(responseText), &analysis); err != nil {
		return 0, "", "", "", errors.New("failed to parse Gemini analysis")
	}

	errors := ""
	if len(analysis.Errors) > 0 {
		b, _ := json.Marshal(analysis.Errors)
		errors = string(b)
	}
	suggestions := ""
	if len(analysis.Suggestions) > 0 {
		b, _ := json.Marshal(analysis.Suggestions)
		suggestions = string(b)
	}

	return analysis.Score, analysis.Feedback, errors, suggestions, nil
}

// Gemini API types
type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content Content `json:"content"`
}

type GrammarAnalysis struct {
	Score       float64  `json:"score"`
	Errors      []Error  `json:"errors"`
	Suggestions []string `json:"suggestions"`
	Feedback    string   `json:"feedback"`
}

type Error struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Position    int    `json:"position"`
	Correction  string `json:"correction"`
}

const GeminiGrammarPrompt = `
    You are an English teacher assisting Vietnamese learners.

    Your task is to evaluate the student's English translation of a Vietnamese sentence and return structured feedback in JSON format. The response must be suitable for educational apps that teach English to Vietnamese users.

    Original Vietnamese sentence: "%s"
    Student's English translation: "%s"
    Target language: %s

    Return a JSON object with the following fields:

    {
        "score": 0 - 100,
        "errors": [
            {
                "type": "grammar | syntax | vocabulary",
                "description": "Simple explanation in Vietnamese to help learners understand the mistake",
                "position": character index of the mistake (or 0 if unknown),
                "correction": "Suggested correction in English"
            }
        ],
        "suggestions": [
            "Learning tips or revision advice in Vietnamese"
        ],
        "feedback": "Write a short comment in Vietnamese to share how you feel - whether it's great, okay, not so good, or anything else you'd like to say."
    }

    Requirements:
    - Use Vietnamese for 'description', 'suggestions', and 'feedback'.
    - Be slightly generous in scoring. Give 100 points if the student's translation is fully correct or only has very minor, acceptable differences (e.g., “Hi” vs “Hello”).
    - Only deduct points for actual mistakes that affect grammar, meaning, or clarity.
    - Do NOT return any markdown, explanation, or extra text. Only respond with the raw JSON object.
`
