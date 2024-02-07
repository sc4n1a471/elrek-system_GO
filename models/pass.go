package models

import "time"
import openapitypes "github.com/oapi-codegen/runtime/types"

// Pass defines model for Pass.
type Pass struct {
	Active        bool              `json:"active,omitempty"`
	Comment       *string           `json:"comment,omitempty"`
	CreatedAt     time.Time         `json:"created_at,omitempty"`
	Duration      *string           `json:"duration,omitempty"`
	ID            openapitypes.UUID `json:"id,omitempty"`
	Name          string            `json:"name,omitempty"`
	OccasionLimit *int              `json:"occasion_limit,omitempty"`
	UserID        openapitypes.UUID `json:"user_id,omitempty" gorm:"size:255"`
	PrevPassID    openapitypes.UUID `json:"prev_pass_id,omitempty"`
	Price         int               `json:"price,omitempty"`
	Services      []Service         `json:"services,omitempty" gorm:"many2many:pass_services;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// PassCreate defines model for PassCreate.
type PassCreate struct {
	Comment       *string           `json:"comment,omitempty"`
	Duration      *string           `json:"duration,omitempty"`
	Name          string            `json:"name"`
	OccasionLimit *int              `json:"occasion_limit,omitempty"`
	UserID        openapitypes.UUID `json:"user_id"`
	Price         int               `json:"price"`
	ServiceID     openapitypes.UUID `json:"service_id"`
}

// PassList defines model for PassList.
type PassList = []Pass

// PassUpdate defines model for PassUpdate.
type PassUpdate struct {
	Comment       *string            `json:"comment,omitempty"`
	Duration      *string            `json:"duration,omitempty"`
	Name          *string            `json:"name,omitempty"`
	OccasionLimit *int               `json:"occasion_limit,omitempty"`
	Price         *int               `json:"price,omitempty"`
	ServiceID     *openapitypes.UUID `json:"service_id,omitempty"`
}
