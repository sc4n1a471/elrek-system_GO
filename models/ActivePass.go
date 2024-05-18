package models

import (
	"time"

	openapitypes "github.com/oapi-codegen/runtime/types"
)

// activePass defines model for activePass.
type ActivePass struct {
	IsActive   bool              `json:"active,omitempty"`
	Comment    *string           `json:"comment,omitempty"`
	CreatedAt  time.Time         `json:"createdAt,omitempty"`
	ID         openapitypes.UUID `json:"id,omitempty"`
	Occasions  int               `json:"occasions"`
	UserID     openapitypes.UUID `json:"userID,omitempty" gorm:"size:255"`
	Pass       Pass              `json:"pass,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PassID     openapitypes.UUID `json:"passID,omitempty" gorm:"size:255"`
	PayerID    openapitypes.UUID `json:"payerID,omitempty" gorm:"size:255"`
	User       *User             `json:"user,omitempty" gorm:"foreignKey:PayerID;;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ValidFrom  *time.Time        `json:"validFrom,omitempty"`
	ValidUntil *time.Time        `json:"validUntil,omitempty"`
}

// activePassCreate defines model for activePassCreate.
type ActivePassCreate struct {
	Comment    *string           `json:"comment,omitempty"`
	UserID     openapitypes.UUID `json:"userID"`
	PassID     openapitypes.UUID `json:"passID"`
	PayerID    openapitypes.UUID `json:"payerID"`
	ValidFrom  time.Time         `json:"validFrom"`
	ValidUntil *time.Time        `json:"validUntil"`
}

// activePassList defines model for activePassList.
type ActivePassList = []ActivePass

// activePassUpdate defines model for activePassUpdate.
type ActivePassUpdate struct {
	Comment    *string    `json:"comment,omitempty"`
	Occasions  *int       `json:"occasions,omitempty"`
	ValidFrom  *time.Time `json:"validFrom,omitempty"`
	ValidUntil *time.Time `json:"validUntil,omitempty"`
}
