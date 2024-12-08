package models

import (
	"time"

	openapitypes "github.com/oapi-codegen/runtime/types"
)

// User defines model for User.
type User struct {
	CreatedAt time.Time          `json:"createdAt,omitempty"`
	Email     openapitypes.Email `json:"email" gorm:"unique"`
	ID        openapitypes.UUID  `json:"id" gorm:"primaryKey,unique,size:255"`
	IsActive  bool               `json:"isActive" gorm:"default:true"`
	IsAdmin   bool               `json:"isAdmin" gorm:"default:false"`
	Name      string             `json:"name,omitempty"`
	OwnerID   openapitypes.UUID  `json:"ownerID,omitempty" gorm:"size:255"`
	Password  []byte             `json:"-"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty"`
	Services  []Service          `json:"services,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Passes    []Pass             `json:"passes,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// ActivePasses is the list of active passes that is connected to the admin user, not the bought passes by the user
	ActivePasses []ActivePass `json:"activePasses,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// BoughtPasses is the list of passes bought by the user, it connects to the ActivePass model, by this field we can get the list of passes bought by the user.
	BoughtPasses []ActivePass `json:"boughtPasses" gorm:"foreignKey:PayerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	MyUsers      []User       `json:"myUsers,omitempty" gorm:"foreignKey:OwnerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
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
	Name     string             `json:"name"`
	Password string             `json:"password"`
	IsAdmin  bool               `json:"isAdmin"`
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

type UserRegister struct {
	Email    openapitypes.Email `json:"email"`
	Name     string             `json:"name"`
	Password string             `json:"password"`
	OwnerID  openapitypes.UUID  `json:"ownerID"`
}

// UserUpdate defines model for UserUpdate.
type UserUpdate struct {
	Email    *openapitypes.Email `json:"email,omitempty"`
	Name     *string             `json:"name,omitempty"`
	Password string              `json:"password"`
	IsAdmin  *bool               `json:"isAdmin"`
}
