package biz

import (
	"context"
	"encoding/json"
	"errors"
	"hub-service/common"
	"hub-service/core/appctx"
	"hub-service/core/auth/tokenprovider"
	"hub-service/core/auth/tokenprovider/jwt"
	"hub-service/module/email/service"
	scoremodel "hub-service/module/score/model"
	scorestorage "hub-service/module/score/storage"
	"hub-service/module/user/model"
	"hub-service/module/user/storage"
	hash "hub-service/utils/hash"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserBiz struct {
	store         *storage.UserStorage
	scoreStore    *scorestorage.Storage
	hasher        hash.Hasher
	tokenProvider tokenprovider.Provider
	emailService  *service.WelcomeEmailService
}

func NewUserBiz(appCtx appctx.AppContext) *UserBiz {
	return &UserBiz{
		store:         storage.NewUserStorage(appCtx),
		scoreStore:    scorestorage.NewStorage(appCtx.GetDatabase()),
		hasher:        hash.NewMd5Hash(),
		tokenProvider: jwt.NewProvider(appCtx.GetSecretKey()),
		emailService:  service.NewWelcomeEmailService(),
	}
}

func (biz *UserBiz) CreateUser(ctx context.Context, userCreate *model.UserCreate) (*model.UserResponse, error) {
	// Check if user already exists
	existingUser, err := biz.store.GetByEmail(ctx, userCreate.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	user, err := biz.store.Create(ctx, userCreate)
	if err != nil {
		return nil, err
	}

	return &model.UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		Name:       user.Name,
		Avatar:     user.Avatar,
		Phone:      user.Phone,
		Bio:        user.Bio,
		Role:       user.Role,
		TotalScore: 0, // New user has no score yet
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}, nil
}

func (biz *UserBiz) Login(ctx context.Context, email, password string) (*model.LoginResponse, error) {
	user, err := biz.store.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Update login status and check if this is first login
	isFirstLogin, err := biz.store.UpdateLoginStatus(ctx, user.ID)
	if err != nil {
		// Log error but don't fail the login
	}

	// Send welcome email if this is the first login
	if isFirstLogin {
		go func() {
			err := biz.emailService.SendWelcomeEmail(user.Name, user.Email)
			if err != nil {
				// Log error but don't fail the login
			}
		}()
	}

	payload := tokenprovider.TokenPayload{
		UserID: user.ID,
		Role:   user.Role,
	}

	// Access token expires in 1 day
	token, err := biz.tokenProvider.Generate(payload, 60*60*24)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}
	loginResponse := &model.LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		User: model.UserResponse{
			ID:         user.ID,
			Email:      user.Email,
			Name:       user.Name,
			Avatar:     user.Avatar,
			Phone:      user.Phone,
			Bio:        user.Bio,
			Role:       user.Role,
			TotalScore: 0, // Will be updated if needed
			CreatedAt:  user.CreatedAt,
			UpdatedAt:  user.UpdatedAt,
		},
	}
	return loginResponse, nil
}

func (biz *UserBiz) RefreshToken(ctx context.Context, refreshToken string) (*model.LoginResponse, error) {
	payload, err := biz.tokenProvider.Validate(refreshToken)
	if err != nil {
		return nil, tokenprovider.ErrInvalidToken
	}

	user, err := biz.store.GetByID(ctx, payload.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	newPayload := tokenprovider.TokenPayload{
		UserID: user.ID,
		Role:   user.Role,
	}

	// Generate only a new access token, keep the old refresh token
	newAccessToken, err := biz.tokenProvider.GenerateAccessToken(newPayload, 60*60*24)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	loginResponse := &model.LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: refreshToken, // Keep the old refresh token
		User: model.UserResponse{
			ID:         user.ID,
			Email:      user.Email,
			Name:       user.Name,
			Avatar:     user.Avatar,
			Phone:      user.Phone,
			Bio:        user.Bio,
			Role:       user.Role,
			TotalScore: 0, // Will be updated if needed
			CreatedAt:  user.CreatedAt,
			UpdatedAt:  user.UpdatedAt,
		},
	}

	return loginResponse, nil
}

func (biz *UserBiz) GetUserByID(ctx context.Context, id primitive.ObjectID) (*model.UserResponse, error) {
	user, err := biz.store.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Get user's total score
	scoreSummary, err := biz.scoreStore.GetUserScoreSummary(ctx, user.ID)
	if err != nil {
		// If error getting score, set to 0 but don't fail the request
		scoreSummary = &scoremodel.UserScoreSummary{
			UserID:     user.ID,
			TotalScore: 0,
		}
	}

	return &model.UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		Name:       user.Name,
		Avatar:     user.Avatar,
		Phone:      user.Phone,
		Bio:        user.Bio,
		Role:       user.Role,
		TotalScore: scoreSummary.TotalScore,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}, nil
}

func (biz *UserBiz) UpdateUser(ctx context.Context, id primitive.ObjectID, userUpdate *model.UserUpdate) (*model.UserResponse, error) {
	user, err := biz.store.Update(ctx, id, userUpdate)
	if err != nil {
		return nil, err
	}

	// Get user's total score
	scoreSummary, err := biz.scoreStore.GetUserScoreSummary(ctx, user.ID)
	if err != nil {
		// If error getting score, set to 0 but don't fail the request
		scoreSummary = &scoremodel.UserScoreSummary{
			UserID:     user.ID,
			TotalScore: 0,
		}
	}

	return &model.UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		Name:       user.Name,
		Avatar:     user.Avatar,
		Phone:      user.Phone,
		Bio:        user.Bio,
		Role:       user.Role,
		TotalScore: scoreSummary.TotalScore,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}, nil
}

func (biz *UserBiz) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	return biz.store.Delete(ctx, id)
}

func (biz *UserBiz) ListUsers(ctx context.Context, page, limit int64, sortBy, sortOrder, search string) ([]model.UserResponse, *model.PaginationMetadata, error) {
	offset := (page - 1) * limit

	users, err := biz.store.List(ctx, limit, offset, sortBy, sortOrder, search)
	if err != nil {
		return nil, nil, err
	}

	totalItems, err := biz.store.CountWithFilter(ctx, search)
	if err != nil {
		return nil, nil, err
	}

	var responses []model.UserResponse
	for _, user := range users {
		// Get user's total score
		scoreSummary, err := biz.scoreStore.GetUserScoreSummary(ctx, user.ID)
		if err != nil {
			// If error getting score, set to 0 but don't fail the request
			scoreSummary = &scoremodel.UserScoreSummary{
				UserID:     user.ID,
				TotalScore: 0,
			}
		}

		responses = append(responses, model.UserResponse{
			ID:         user.ID,
			Email:      user.Email,
			Name:       user.Name,
			Avatar:     user.Avatar,
			Phone:      user.Phone,
			Bio:        user.Bio,
			Role:       user.Role,
			TotalScore: scoreSummary.TotalScore,
			CreatedAt:  user.CreatedAt,
			UpdatedAt:  user.UpdatedAt,
		})
	}

	totalPages := (totalItems + limit - 1) / limit
	hasNext := page < totalPages
	hasPrev := page > 1

	metadata := &model.PaginationMetadata{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}

	return responses, metadata, nil
}

func (biz *UserBiz) SocialLogin(ctx context.Context, req *model.SocialLoginRequest) (*model.LoginResponse, error) {
	// 1. Verify id_token with Google
	resp, err := http.Get(os.Getenv("SYSTEM_GOOGLE_AUTHENTICATOR") + req.IdToken)
	if err != nil {
		return nil, errors.New("failed to verify id_token with Google")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("invalid id_token")
	}
	var googleResp struct {
		Email   string `json:"email"`
		Sub     string `json:"sub"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&googleResp); err != nil {
		return nil, errors.New("failed to decode google response")
	}
	if googleResp.Email != req.Email {
		return nil, errors.New("email does not match id_token")
	}

	// 2. Find or create user
	user, err := biz.store.GetByEmail(ctx, req.Email)
	isNewUser := false
	if err != nil {
		return nil, err
	}
	if user == nil {
		// Create new user
		userCreate := &model.UserCreate{
			Email:      req.Email,
			Name:       req.Name,
			Avatar:     req.Avatar,
			Role:       common.RoleClient,
			Provider:   req.Provider,
			ProviderID: googleResp.Sub,
		}
		user, err = biz.store.Create(ctx, userCreate)
		if err != nil {
			return nil, err
		}
		isNewUser = true
	} else {
		// Update existing user
		update := false
		updateData := make(map[string]interface{})
		if user.Name != req.Name {
			updateData["name"] = req.Name
			update = true
		}
		if user.Avatar != req.Avatar {
			updateData["avatar"] = req.Avatar
			update = true
		}
		if user.Provider != req.Provider {
			updateData["provider"] = req.Provider
			update = true
		}
		if user.ProviderID != googleResp.Sub {
			updateData["provider_id"] = googleResp.Sub
			update = true
		}
		if update {
			if err := biz.store.UpdateFields(ctx, user.ID, updateData); err == nil {
				// cập nhật lại user struct với dữ liệu mới
				if v, ok := updateData["name"]; ok {
					user.Name = v.(string)
				}
				if v, ok := updateData["avatar"]; ok {
					user.Avatar = v.(string)
				}
				if v, ok := updateData["provider"]; ok {
					user.Provider = v.(string)
				}
				if v, ok := updateData["provider_id"]; ok {
					user.ProviderID = v.(string)
				}
			}
		}

		// Update login status for existing user
		_, err = biz.store.UpdateLoginStatus(ctx, user.ID)
		if err != nil {
			// Log error but don't fail the login
		}
	}

	// Send welcome email if this is a new user
	if isNewUser {
		go func() {
			err := biz.emailService.SendWelcomeEmail(user.Name, user.Email)
			if err != nil {
				// Log error but don't fail the login
			}
		}()
	}

	payload := tokenprovider.TokenPayload{
		UserID: user.ID,
		Role:   user.Role,
	}
	token, err := biz.tokenProvider.Generate(payload, 60*60*24)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}
	loginResponse := &model.LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		User: model.UserResponse{
			ID:         user.ID,
			Email:      user.Email,
			Name:       user.Name,
			Avatar:     user.Avatar,
			Phone:      user.Phone,
			Bio:        user.Bio,
			Role:       user.Role,
			TotalScore: 0, // New user or will be updated if needed
			CreatedAt:  user.CreatedAt,
			UpdatedAt:  user.UpdatedAt,
		},
	}
	return loginResponse, nil
}

func (biz *UserBiz) StoreUpdateRole(ctx context.Context, id primitive.ObjectID, role string) error {
	return biz.store.UpdateRole(ctx, id, role)
}

func (biz *UserBiz) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return biz.store.GetByEmail(ctx, email)
}
