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

	challengemodel "hub-service/module/challenge/model"
	challengestorage "hub-service/module/challenge/storage"
	scoremodel "hub-service/module/score/model"
	scorestorage "hub-service/module/score/storage"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrChallengeNotFound = errors.New("challenge not found")
)

type GeminiAnalyzer interface {
	AnalyzeGrammar(ctx context.Context, originalText, userTranslation, targetLanguage string) (*GrammarAnalysis, error)
}

type GeminiBiz struct {
	apiKey  string
	client  *http.Client
	baseURL string
}

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

func NewGeminiBiz(apiKey string, baseURL string) *GeminiBiz {
	return &GeminiBiz{
		apiKey:  apiKey,
		client:  &http.Client{},
		baseURL: baseURL,
	}
}

func (b *GeminiBiz) AnalyzeGrammar(ctx context.Context, originalText, userTranslation, targetLanguage string) (*GrammarAnalysis, error) {
	promptText := fmt.Sprintf(GeminiGrammarPrompt, originalText, userTranslation, targetLanguage)

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
		return nil, errors.New("failed to marshal request")
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", b.baseURL, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, errors.New("failed to create request")
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-goog-api-key", b.apiKey)

	resp, err := b.client.Do(httpReq)
	if err != nil {
		return nil, errors.New("failed to execute request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Gemini error status: %d, body: %s\n", resp.StatusCode, string(body))
		return nil, fmt.Errorf("status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var geminiResp GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return nil, errors.New("failed to decode response")
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, errors.New("no response from Gemini")
	}

	responseText := geminiResp.Candidates[0].Content.Parts[0].Text
	fmt.Printf("Gemini responseText: %s\n", responseText)

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
		fmt.Printf("Gemini JSON parse error: %v\n", err)
		return nil, errors.New("failed to parse Gemini analysis")
	}

	return &analysis, nil
}

type ChallengeStore interface {
	Get(ctx context.Context, id primitive.ObjectID) (*challengemodel.Challenge, error)
}

type ScoreBiz struct {
	scoreStorage     *scorestorage.Storage
	challengeStorage *challengestorage.Storage
	geminiBiz        GeminiAnalyzer
}

func NewScoreBiz(scoreStorage *scorestorage.Storage, challengeStorage *challengestorage.Storage, geminiBiz GeminiAnalyzer) *ScoreBiz {
	return &ScoreBiz{
		scoreStorage:     scoreStorage,
		challengeStorage: challengeStorage,
		geminiBiz:        geminiBiz,
	}
}

func (biz *ScoreBiz) SubmitScore(ctx context.Context, userID primitive.ObjectID, req *scoremodel.SubmitScoreRequest, targetLanguage string) (*scoremodel.SubmitScoreResponse, error) {
	challengeID, err := primitive.ObjectIDFromHex(req.ChallengeID)
	if err != nil {
		return nil, err
	}

	challenge, err := biz.challengeStorage.Get(ctx, challengeID)
	if err != nil {
		return nil, err
	}
	if challenge == nil {
		return nil, ErrChallengeNotFound
	}

	analysis, err := biz.geminiBiz.AnalyzeGrammar(ctx, challenge.Content, req.UserTranslation, targetLanguage)
	if err != nil {
		return nil, err
	}

	existingScore, err := biz.scoreStorage.GetScoreByUserAndChallenge(ctx, userID, challengeID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var attemptCount int
	var bestScore float64
	isNewBest := false

	geminiErrors := ""
	if len(analysis.Errors) > 0 {
		b, _ := json.Marshal(analysis.Errors)
		geminiErrors = string(b)
	}
	geminiSuggestions := ""
	if len(analysis.Suggestions) > 0 {
		b, _ := json.Marshal(analysis.Suggestions)
		geminiSuggestions = string(b)
	}

	if existingScore == nil {
		attemptCount = 1
		bestScore = analysis.Score
		isNewBest = true

		scoreCreate := &scoremodel.ScoreCreate{
			UserID:            userID,
			ChallengeID:       challengeID,
			UserTranslation:   req.UserTranslation,
			GeminiScore:       analysis.Score,
			GeminiFeedback:    analysis.Feedback,
			GeminiErrors:      geminiErrors,
			GeminiSuggestions: geminiSuggestions,
			OriginalContent:   challenge.Content,
			AttemptCount:      attemptCount,
			BestScore:         bestScore,
			CreatedAt:         now,
			UpdatedAt:         now,
		}
		err = biz.scoreStorage.CreateScore(ctx, scoreCreate)
	} else {
		attemptCount = existingScore.AttemptCount + 1
		bestScore = existingScore.BestScore
		if analysis.Score > existingScore.BestScore {
			bestScore = analysis.Score
			isNewBest = true
		}
		scoreUpdate := &scoremodel.ScoreUpdate{
			UserTranslation:   &req.UserTranslation,
			GeminiScore:       &analysis.Score,
			GeminiFeedback:    &analysis.Feedback,
			GeminiErrors:      &geminiErrors,
			GeminiSuggestions: &geminiSuggestions,
			AttemptCount:      &attemptCount,
			BestScore:         &bestScore,
			UpdatedAt:         &now,
		}
		err = biz.scoreStorage.UpdateScore(ctx, existingScore.ID, scoreUpdate)
	}

	if err != nil {
		return nil, err
	}

	return &scoremodel.SubmitScoreResponse{
		Score:             analysis.Score,
		UserTranslation:   req.UserTranslation,
		GeminiFeedback:    analysis.Feedback,
		GeminiErrors:      geminiErrors,
		GeminiSuggestions: geminiSuggestions,
		OriginalContent:   challenge.Content,
		AttemptCount:      attemptCount,
		BestScore:         bestScore,
		IsNewBest:         isNewBest,
	}, nil
}

func (biz *ScoreBiz) GetUserScores(ctx context.Context, userID primitive.ObjectID) (*scoremodel.GetUserScoresResponse, error) {
	summary, err := biz.scoreStorage.GetUserScoreSummary(ctx, userID)
	if err != nil {
		return nil, err
	}

	scores, err := biz.scoreStorage.GetUserScores(ctx, userID)
	if err != nil {
		return nil, err
	}

	challengeScores := make([]scoremodel.ChallengeScore, len(scores))
	for i, score := range scores {
		challenge, err := biz.challengeStorage.Get(ctx, score.ChallengeID)
		if err != nil {
			continue
		}
		challengeScores[i] = scoremodel.ChallengeScore{
			ChallengeID:       score.ChallengeID,
			ChallengeTitle:    challenge.Title,
			BestScore:         score.BestScore,
			AttemptCount:      score.AttemptCount,
			LastAttemptAt:     score.UpdatedAt,
			UserTranslation:   score.UserTranslation,
			GeminiFeedback:    score.GeminiFeedback,
			GeminiErrors:      score.GeminiErrors,
			GeminiSuggestions: score.GeminiSuggestions,
			OriginalContent:   score.OriginalContent,
		}
	}

	return &scoremodel.GetUserScoresResponse{
		Summary: *summary,
		Scores:  challengeScores,
	}, nil
}
