package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`
	Name      string             `bson:"name" json:"name"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type UserCreate struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

type UserUpdate struct {
	Name string `json:"name" binding:"required"`
}

type UserResponse struct {
	ID        primitive.ObjectID `json:"id"`
	Email     string             `json:"email"`
	Name      string             `json:"name"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}

// API Response Models for Swagger
type CreateUserResponse struct {
	Data UserResponse `json:"data"`
}

type LoginAPIResponse struct {
	Data LoginResponse `json:"data"`
}

type GetUserResponse struct {
	Data UserResponse `json:"data"`
}

type UpdateUserResponse struct {
	Data UserResponse `json:"data"`
}

type DeleteUserResponse struct {
	Message string `json:"message"`
}

type ListUsersResponse struct {
	Data []UserResponse `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
