package models

import (
	"time"

	openapitypes "github.com/oapi-codegen/runtime/types"
)

// activePass defines model for activePass.
type ActivePass struct {
	IsActive   bool              `json:"active,omitempty"`
	Comment    *string           `json:"comment,omitempty"`
	CreatedAt  time.Time         `json:"created_at,omitempty"`
	ID         openapitypes.UUID `json:"id,omitempty"`
	Occasions  int               `json:"occasions"`
	UserID     openapitypes.UUID `json:"user_id,omitempty" gorm:"size:255"`
	Pass       Pass              `json:"pass,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PassID     openapitypes.UUID `json:"pass_id,omitempty" gorm:"size:255"`
	PayerID    openapitypes.UUID `json:"payer_id,omitempty" gorm:"size:255"`
	User       *User             `json:"user,omitempty" gorm:"foreignKey:PayerID;;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ValidFrom  *time.Time        `json:"valid_from,omitempty"`
	ValidUntil *time.Time        `json:"valid_until,omitempty"`
}

// activePassCreate defines model for activePassCreate.
type ActivePassCreate struct {
	Comment    *string           `json:"comment,omitempty"`
	UserID     openapitypes.UUID `json:"user_id"`
	PassID     openapitypes.UUID `json:"pass_id"`
	PayerID    openapitypes.UUID `json:"payer_id"`
	ValidFrom  time.Time         `json:"valid_from"`
	ValidUntil *time.Time        `json:"valid_until"`
}

// activePassList defines model for activePassList.
type ActivePassList = []ActivePass

// activePassUpdate defines model for activePassUpdate.
type ActivePassUpdate struct {
	Comment    *string    `json:"comment,omitempty"`
	Occasions  *int       `json:"occasions,omitempty"`
	ValidFrom  *time.Time `json:"valid_from,omitempty"`
	ValidUntil *time.Time `json:"valid_until,omitempty"`
}
