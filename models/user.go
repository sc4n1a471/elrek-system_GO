package models

import (
	"time"
)
import openapitypes "github.com/oapi-codegen/runtime/types"

// User defines model for User.
type User struct {
	CreatedAt *time.Time          `json:"created_at,omitempty"`
	Email     *openapitypes.Email `json:"email,omitempty"`
	Id        *openapitypes.UUID  `json:"id,omitempty"`
	IsActive  *bool               `json:"is_active,omitempty"`
	IsAdmin   *bool               `json:"is_admin,omitempty"`
	Name      *string             `json:"name,omitempty"`
	OwnerId   *openapitypes.Email `json:"owner_id,omitempty"`
	Password  *string             `json:"password,omitempty"`
	UpdatedAt *time.Time          `json:"updated_at,omitempty"`
}

// UserCreate defines model for UserCreate.
type UserCreate struct {
	Email    openapitypes.Email `json:"email"`
	Name     string             `json:"name"`
	OwnerId  openapitypes.Email `json:"owner_id"`
	Password string             `json:"password"`
}

// UserListErrorResponse defines model for UserListErrorResponse.
type UserListErrorResponse struct {
	Message *string        `json:"message,omitempty"`
	Users   *[]interface{} `json:"users,omitempty"`
}

// UserListSuccessResponse defines model for UserListSuccessResponse.
type UserListSuccessResponse struct {
	Message *string `json:"message,omitempty"`
	Users   *[]User `json:"users,omitempty"`
}

// UserLogin defines model for UserLogin.
type UserLogin struct {
	Email    openapitypes.Email `json:"email"`
	Password string             `json:"password"`
}

// UserUpdate defines model for UserUpdate.
type UserUpdate struct {
	Email    *openapitypes.Email `json:"email,omitempty"`
	Name     *string             `json:"name,omitempty"`
	Password string              `json:"password"`
}
