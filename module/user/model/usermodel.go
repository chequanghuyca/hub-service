package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email      string             `bson:"email" json:"email"`
	Name       string             `bson:"name" json:"name"`
	Avatar     string             `bson:"avatar,omitempty" json:"avatar,omitempty"`
	Role       string             `bson:"role" json:"role"`
	Provider   string             `bson:"provider,omitempty" json:"provider,omitempty"`
	ProviderID string             `bson:"provider_id,omitempty" json:"provider_id,omitempty"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

type UserCreate struct {
	Email      string `json:"email" binding:"required,email"`
	Name       string `json:"name" binding:"required"`
	Avatar     string `json:"avatar"`
	Role       string `json:"role" binding:"required,oneof=admin client super_admin"`
	Provider   string `json:"provider,omitempty"`
	ProviderID string `json:"provider_id,omitempty"`
}

type UserUpdate struct {
	Name string `json:"name,omitempty"`
}

type UserResponse struct {
	ID        primitive.ObjectID `json:"id"`
	Email     string             `json:"email"`
	Name      string             `json:"name"`
	Avatar    string             `json:"avatar,omitempty"`
	Role      string             `json:"role"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// API Response Models for Swagger
type CreateUserResponse struct {
	Status string       `json:"status"`
	Data   UserResponse `json:"data"`
}

type LoginAPIResponse struct {
	Status string        `json:"status"`
	Data   LoginResponse `json:"data"`
}

type GetUserResponse struct {
	Status string       `json:"status"`
	Data   UserResponse `json:"data"`
}

type UpdateUserResponse struct {
	Status string       `json:"status"`
	Data   UserResponse `json:"data"`
}

type DeleteUserResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// PaginationMetadata contains pagination information
type PaginationMetadata struct {
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int64 `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

type ListUsersResponse struct {
	Status   string             `json:"status"`
	Data     []UserResponse     `json:"data"`
	Metadata PaginationMetadata `json:"metadata"`
}

// ListUsersRequest represents the request parameters for listing users
type ListUsersRequest struct {
	Page      int64  `form:"page" binding:"omitempty,min=1" example:"1"`
	Limit     int64  `form:"limit" binding:"omitempty,min=1,max=100" example:"10"`
	SortBy    string `form:"sort_by" binding:"omitempty,oneof=name email created_at updated_at" example:"created_at"`
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc" example:"desc"`
	Search    string `form:"search" binding:"omitempty" example:"huy"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SocialLoginRequest struct {
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Avatar     string `json:"avatar"`
	IdToken    string `json:"id_token" binding:"required"`
	Provider   string `json:"provider" binding:"required"`
	ProviderID string `json:"provider_id" binding:"required"`
}

type UpdateRoleRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required,oneof=admin client super_admin"`
}
