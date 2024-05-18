package models

import (
	"time"

	openapitypes "github.com/oapi-codegen/runtime/types"
)

// Pass defines model for Pass.
type Pass struct {
	IsActive      bool              `json:"isActive,omitempty"`
	Comment       *string           `json:"comment,omitempty"`
	CreatedAt     time.Time         `json:"createdAt,omitempty"`
	Duration      *string           `json:"duration,omitempty"`
	ID            openapitypes.UUID `json:"id,omitempty"`
	Name          string            `json:"name,omitempty"`
	OccasionLimit *int              `json:"occasionLimit,omitempty"`
	UserID        openapitypes.UUID `json:"userID,omitempty" gorm:"size:255"`
	PrevPassID    openapitypes.UUID `json:"prevPassID,omitempty"`
	Price         int               `json:"price,omitempty"`
	Services      []Service         `json:"services,omitempty" gorm:"many2many:pass_services;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// PassCreate defines model for PassCreate.
type PassCreate struct {
	Comment       *string             `json:"comment,omitempty"`
	Duration      *string             `json:"duration,omitempty"`
	Name          string              `json:"name"`
	OccasionLimit *int                `json:"occasionLimit,omitempty"`
	UserID        openapitypes.UUID   `json:"userID"`
	Price         int                 `json:"price"`
	ServiceIDs    []openapitypes.UUID `json:"serviceIDs"`
}

// PassUpdate defines model for PassUpdate.
type PassUpdate struct {
	Comment       *string              `json:"comment,omitempty"`
	Duration      *string              `json:"duration,omitempty"`
	Name          *string              `json:"name,omitempty"`
	OccasionLimit *int                 `json:"occasionLimit,omitempty"`
	Price         *int                 `json:"price,omitempty"`
	ServiceIDs    *[]openapitypes.UUID `json:"serviceIDs,omitempty"`
}
