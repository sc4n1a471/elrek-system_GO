package models

import (
	"time"

	openapitypes "github.com/oapi-codegen/runtime/types"
)

// User defines model for User.
type User struct {
	CreatedAt    time.Time          `json:"created_at,omitempty"`
	Email        openapitypes.Email `json:"email" gorm:"unique"`
	ID           openapitypes.UUID  `json:"id" gorm:"primaryKey,unique,size:255"`
	IsActive     bool               `json:"is_active" gorm:"default:true"`
	IsAdmin      bool               `json:"is_admin" gorm:"default:false"`
	Name         string             `json:"name,omitempty"`
	OwnerID      openapitypes.UUID  `json:"owner_id,omitempty" gorm:"size:255"`
	Password     []byte             `json:"-"`
	UpdatedAt    time.Time          `json:"updated_at,omitempty"`
	Services     []Service          `json:"services,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Passes       []Pass             `json:"passes,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ActivePasses []ActivePass       `json:"active_passes,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	BoughtPasses []ActivePass       `json:"bought_passes,omitempty" gorm:"foreignKey:PayerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	MyUsers      []User             `json:"my_users,omitempty" gorm:"foreignKey:OwnerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// UserResponse defines model for UserResponse.
type UserResponse struct {
	Email    openapitypes.Email `json:"email,omitempty"`
	ID       openapitypes.UUID  `json:"id,omitempty"`
	IsActive bool               `json:"isActive"`
	IsAdmin  bool               `json:"isAdmin"`
	Name     string             `json:"name,omitempty"`
}

// UserCreate defines model for UserCreate.
type UserCreate struct {
	Email    openapitypes.Email `json:"email"`
	Name     string             `json:"name,unique"`
	Password string             `json:"password"`
	IsAdmin  bool               `json:"isAdmin,omitempty,default:false"`
}

// UserLogin defines model for UserLogin.
type UserLogin struct {
	Email    openapitypes.Email `json:"email"`
	Password string             `json:"password"`
}

// UserLoginResponse defines model for UserLoginResponse.
type UserLoginResponse struct {
	Email   openapitypes.Email `json:"email,omitempty"`
	ID      openapitypes.UUID  `json:"id,omitempty"`
	Name    string             `json:"name,omitempty"`
	IsAdmin bool               `json:"isAdmin"`
}

// UserUpdate defines model for UserUpdate.
type UserUpdate struct {
	Email    *openapitypes.Email `json:"email,omitempty"`
	Name     *string             `json:"name,omitempty"`
	Password string              `json:"password"`
	IsAdmin  *bool               `json:"is_admin,omitempty,default:false"`
}
