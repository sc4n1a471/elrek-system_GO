package models

import (
	"time"
)
import openapitypes "github.com/oapi-codegen/runtime/types"

// User defines model for User.
type User struct {
	CreatedAt *time.Time          `json:"created_at,omitempty"`
	Email     *openapitypes.Email `json:"email,omitempty"`
	Id        *openapitypes.UUID  `json:"id,omitempty,primaryKey"`
	IsActive  bool                `json:"is_active,omitempty,default:true"`
	IsAdmin   bool                `json:"is_admin,omitempty,default:false"`
	Name      *string             `json:"name,omitempty"`
	OwnerId   *openapitypes.UUID  `json:"owner_id,omitempty"`
	Password  []byte              `json:"-"`
	UpdatedAt *time.Time          `json:"updated_at,omitempty"`
}

// UserLoginResponse defines model for UserLoginResponse.
type UserLoginResponse struct {
	Email *openapitypes.Email `json:"email,omitempty"`
	Id    *openapitypes.UUID  `json:"id,omitempty"`
	Name  *string             `json:"name,omitempty"`
	Token *string             `json:"token,omitempty"`
}

// UserResponse defines model for UserResponse.
type UserResponse struct {
	Email    *openapitypes.Email `json:"email,omitempty"`
	Id       *openapitypes.UUID  `json:"id,omitempty"`
	IsActive bool                `json:"is_active"`
	IsAdmin  bool                `json:"is_admin"`
	Name     *string             `json:"name,omitempty"`
}

// UserCreate defines model for UserCreate.
type UserCreate struct {
	Email    openapitypes.Email `json:"email"`
	Name     string             `json:"name"`
	Password string             `json:"password"`
	IsAdmin  bool               `json:"is_admin,omitempty,default:false"`
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
	IsAdmin  *bool               `json:"is_admin,omitempty,default:false"`
}
