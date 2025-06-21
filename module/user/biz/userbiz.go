package biz

import (
	"context"
	"errors"
	"hub-service/component/appctx"
	"hub-service/component/hasher"
	"hub-service/component/tokenprovider"
	"hub-service/module/user/model"
	"hub-service/module/user/storage"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserBiz struct {
	appCtx    appctx.AppContext
	store     *storage.UserStorage
	hasher    hasher.Hasher
	tokenProv tokenprovider.Provider
}

func NewUserBiz(appCtx appctx.AppContext) *UserBiz {
	return &UserBiz{
		appCtx:    appCtx,
		store:     storage.NewUserStorage(appCtx),
		hasher:    hasher.NewMd5Hash(),
		tokenProv: tokenprovider.NewJWTProvider(appCtx.SecretKey()),
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
	// Get user by email
	user, err := biz.store.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Verify password
	hashedPassword := biz.hasher.Hash(password)
	if hashedPassword != user.Password {
		return nil, errors.New("invalid password")
	}

	// Generate access token
	userId, _ := strconv.Atoi(user.ID.Hex()[:8]) // Convert ObjectID to int for compatibility
	payload := tokenprovider.AccessTokenPayload{
		UserId:    userId,
		Role:      user.Role, // Use role from DB
		Email:     user.Email,
		FirstName: user.Name,
		LastName:  "",
	}

	token, err := biz.tokenProv.Generate(payload, 24*60*60) // 24 hours
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		AccessToken: token.AccessToken,
		User: model.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			Avatar:    user.Avatar,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
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

func (biz *UserBiz) ListUsers(ctx context.Context, page, limit int64) ([]model.UserResponse, error) {
	// Calculate offset from page (page starts from 1)
	offset := (page - 1) * limit

	users, err := biz.store.List(ctx, limit, offset)
	if err != nil {
		return nil, err
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

	return responses, nil
}
