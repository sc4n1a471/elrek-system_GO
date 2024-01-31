package models

import (
	"time"
)
import openapitypes "github.com/oapi-codegen/runtime/types"

// User defines model for User.
type User struct {
	CreatedAt *time.Time         `json:"created_at,omitempty"`
	Email     openapitypes.Email `gorm:"unique"`
	Id        openapitypes.UUID  `gorm:"primaryKey"`
	IsActive  bool               `gorm:"default:true"`
	IsAdmin   bool               `gorm:"default:false"`
	Name      *string            `json:"name,omitempty"`
	OwnerId   openapitypes.UUID  `json:"owner_id,omitempty"`
	Password  []byte             `json:"-"`
	UpdatedAt *time.Time         `json:"updated_at,omitempty"`
}

// UserResponse defines model for UserResponse.
type UserResponse struct {
	Email    openapitypes.Email `json:"email,omitempty"`
	Id       openapitypes.UUID  `json:"id,omitempty"`
	IsActive bool               `json:"is_active"`
	IsAdmin  bool               `json:"is_admin"`
	Name     *string            `json:"name,omitempty"`
}

// UserCreate defines model for UserCreate.
type UserCreate struct {
	Email    openapitypes.Email `json:"email"`
	Name     string             `json:"name,unique"`
	Password string             `json:"password"`
	IsAdmin  bool               `json:"is_admin,omitempty,default:false"`
}

// UserLogin defines model for UserLogin.
type UserLogin struct {
	Email    openapitypes.Email `json:"email"`
	Password string             `json:"password"`
}

// UserLoginResponse defines model for UserLoginResponse.
type UserLoginResponse struct {
	Email   openapitypes.Email `json:"email,omitempty"`
	Id      openapitypes.UUID  `json:"id,omitempty"`
	Name    *string            `json:"name,omitempty"`
	IsAdmin bool               `json:"is_admin"`
}

// UserUpdate defines model for UserUpdate.
type UserUpdate struct {
	Email    *openapitypes.Email `json:"email,omitempty"`
	Name     *string             `json:"name,omitempty"`
	Password string              `json:"password"`
	IsAdmin  *bool               `json:"is_admin,omitempty,default:false"`
}
