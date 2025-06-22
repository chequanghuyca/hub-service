package biz

import (
	"context"
	"errors"
	"hub-service/core/appctx"
	"hub-service/core/auth/tokenprovider"
	"hub-service/core/auth/tokenprovider/jwt"
	"hub-service/module/user/model"
	"hub-service/module/user/storage"
	hash "hub-service/utils/hash"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserBiz struct {
	store         *storage.UserStorage
	hasher        hash.Hasher
	tokenProvider tokenprovider.Provider
}

func NewUserBiz(appCtx appctx.AppContext) *UserBiz {
	return &UserBiz{
		store:         storage.NewUserStorage(appCtx),
		hasher:        hash.NewMd5Hash(),
		tokenProvider: jwt.NewProvider(appCtx.GetSecretKey()),
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

	// Hash password
	hashedPassword := biz.hasher.Hash(userCreate.Password)

	// Create user with hashed password
	userCreateWithHash := &model.UserCreate{
		Email:    userCreate.Email,
		Password: hashedPassword,
		Name:     userCreate.Name,
		Avatar:   userCreate.Avatar,
		Role:     userCreate.Role,
	}

	user, err := biz.store.Create(ctx, userCreateWithHash)
	if err != nil {
		return nil, err
	}

	return &model.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Avatar:    user.Avatar,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (biz *UserBiz) Login(ctx context.Context, email, password string) (*model.LoginResponse, error) {
	user, err := biz.store.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := biz.hasher.CheckPassword(user.Password, password); err != nil {
		return nil, errors.New("invalid credentials")
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
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			Avatar:    user.Avatar,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
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

	// Generate a new pair of tokens, new access token expires in 1 day
	token, err := biz.tokenProvider.Generate(newPayload, 60*60*24)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	loginResponse := &model.LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken, // Return the new refresh token as well
		User: model.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			Avatar:    user.Avatar,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
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

	return &model.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Avatar:    user.Avatar,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (biz *UserBiz) UpdateUser(ctx context.Context, id primitive.ObjectID, userUpdate *model.UserUpdate) (*model.UserResponse, error) {
	user, err := biz.store.Update(ctx, id, userUpdate)
	if err != nil {
		return nil, err
	}

	return &model.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Avatar:    user.Avatar,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
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
		responses = append(responses, model.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			Avatar:    user.Avatar,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
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
