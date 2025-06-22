package biz

import (
	"context"
	"errors"
	"time"

	challengestorage "hub-service/module/challenge/storage"
	scoremodel "hub-service/module/score/model"
	scorestorage "hub-service/module/score/storage"
	translatebiz "hub-service/module/translate/biz"
	translatemodel "hub-service/module/translate/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrChallengeNotFound = errors.New("challenge not found")
)

type ScoreBiz struct {
	scoreStorage     *scorestorage.Storage
	challengeStorage *challengestorage.Storage
	translateBiz     *translatebiz.TranslateBiz
}

func NewScoreBiz(scoreStorage *scorestorage.Storage, challengeStorage *challengestorage.Storage, translateBiz *translatebiz.TranslateBiz) *ScoreBiz {
	return &ScoreBiz{
		scoreStorage:     scoreStorage,
		challengeStorage: challengeStorage,
		translateBiz:     translateBiz,
	}
}

func (biz *ScoreBiz) SubmitScore(ctx context.Context, userID primitive.ObjectID, req *scoremodel.SubmitScoreRequest) (*scoremodel.SubmitScoreResponse, error) {
	// Parse challenge ID
	challengeID, err := primitive.ObjectIDFromHex(req.ChallengeID)
	if err != nil {
		return nil, err
	}

	// Get challenge details
	challenge, err := biz.challengeStorage.Get(ctx, challengeID)
	if err != nil {
		return nil, err
	}
	if challenge == nil {
		return nil, ErrChallengeNotFound
	}

	// Get DeepL score for user translation
	scoreRequest := &translatemodel.ScoreRequest{
		ChallengeID:     req.ChallengeID,
		SentenceIndex:   0, // For now, we'll use the first sentence
		UserTranslation: req.UserTranslation,
	}

	scoreResponse, err := biz.translateBiz.ScoreTranslation(ctx, *scoreRequest)
	if err != nil {
		return nil, err
	}

	// Check if user already has a score for this challenge
	existingScore, err := biz.scoreStorage.GetScoreByUserAndChallenge(ctx, userID, challengeID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var attemptCount int
	var bestScore float64
	isNewBest := false

	if existingScore == nil {
		// This is the first attempt
		attemptCount = 1
		bestScore = scoreResponse.Score
		isNewBest = true

		scoreCreate := &scoremodel.ScoreCreate{
			UserID:           userID,
			ChallengeID:      challengeID,
			UserTranslation:  req.UserTranslation,
			DeepLScore:       scoreResponse.Score,
			DeepLTranslation: scoreResponse.DeepLTranslation,
			OriginalContent:  challenge.Content,
			AttemptCount:     attemptCount,
			BestScore:        bestScore,
			CreatedAt:        now,
			UpdatedAt:        now,
		}
		err = biz.scoreStorage.CreateScore(ctx, scoreCreate)
	} else {
		// This is a subsequent attempt
		attemptCount = existingScore.AttemptCount + 1
		bestScore = existingScore.BestScore

		if scoreResponse.Score > existingScore.BestScore {
			bestScore = scoreResponse.Score
			isNewBest = true
		}

		scoreUpdate := &scoremodel.ScoreUpdate{
			UserTranslation:  &req.UserTranslation,
			DeepLScore:       &scoreResponse.Score,
			DeepLTranslation: &scoreResponse.DeepLTranslation,
			AttemptCount:     &attemptCount,
			BestScore:        &bestScore,
			UpdatedAt:        &now,
		}
		err = biz.scoreStorage.UpdateScore(ctx, existingScore.ID, scoreUpdate)
	}

	if err != nil {
		return nil, err
	}

	return &scoremodel.SubmitScoreResponse{
		Score:            scoreResponse.Score,
		UserTranslation:  req.UserTranslation,
		DeepLTranslation: scoreResponse.DeepLTranslation,
		OriginalContent:  challenge.Content,
		AttemptCount:     attemptCount,
		BestScore:        bestScore,
		IsNewBest:        isNewBest,
	}, nil
}

func (biz *ScoreBiz) GetUserScores(ctx context.Context, userID primitive.ObjectID) (*scoremodel.GetUserScoresResponse, error) {
	// Get user score summary
	summary, err := biz.scoreStorage.GetUserScoreSummary(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get all user scores
	scores, err := biz.scoreStorage.GetUserScores(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Convert to challenge scores
	challengeScores := make([]scoremodel.ChallengeScore, len(scores))
	for i, score := range scores {
		// Get challenge details
		challenge, err := biz.challengeStorage.Get(ctx, score.ChallengeID)
		if err != nil {
			continue // Skip if challenge not found
		}

		// Handle legacy data where best_score might be 0 or missing
		bestScoreForDisplay := score.BestScore
		if bestScoreForDisplay == 0 {
			bestScoreForDisplay = score.DeepLScore
		}

		challengeScores[i] = scoremodel.ChallengeScore{
			ChallengeID:      score.ChallengeID,
			ChallengeTitle:   challenge.Title,
			BestScore:        bestScoreForDisplay,
			AttemptCount:     score.AttemptCount,
			LastAttemptAt:    score.UpdatedAt,
			UserTranslation:  score.UserTranslation,
			DeepLTranslation: score.DeepLTranslation,
			OriginalContent:  score.OriginalContent,
		}
	}

	return &scoremodel.GetUserScoresResponse{
		Summary: *summary,
		Scores:  challengeScores,
	}, nil
}

func (biz *ScoreBiz) GetTotalScore(ctx context.Context, userID primitive.ObjectID) (*scoremodel.GetTotalScoreResponse, error) {
	return biz.scoreStorage.GetTotalScore(ctx, userID)
}
