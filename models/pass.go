package models

import "time"
import openapitypes "github.com/oapi-codegen/runtime/types"

// Pass defines model for Pass.
type Pass struct {
	Active        *bool              `json:"active,omitempty"`
	Comment       *string            `json:"comment,omitempty"`
	CreatedAt     *time.Time         `json:"created_at,omitempty"`
	Duration      *string            `json:"duration,omitempty"`
	Id            *openapitypes.UUID `json:"id,omitempty"`
	Name          *string            `json:"name,omitempty"`
	OccasionLimit *int               `json:"occasion_limit,omitempty"`
	OwnerId       *openapitypes.UUID `json:"owner_id,omitempty"`
	PrevPassId    *openapitypes.UUID `json:"prev_pass_id,omitempty"`
	Price         *float32           `json:"price,omitempty"`
	ServiceId     *openapitypes.UUID `json:"service_id,omitempty"`
}

// PassCreate defines model for PassCreate.
type PassCreate struct {
	Comment       *string           `json:"comment,omitempty"`
	Duration      *string           `json:"duration,omitempty"`
	Name          string            `json:"name"`
	OccasionLimit *int              `json:"occasion_limit,omitempty"`
	OwnerId       openapitypes.UUID `json:"owner_id"`
	Price         float32           `json:"price"`
	ServiceId     openapitypes.UUID `json:"service_id"`
}

// PassList defines model for PassList.
type PassList = []Pass

// PassUpdate defines model for PassUpdate.
type PassUpdate struct {
	Comment       *string            `json:"comment,omitempty"`
	Duration      *string            `json:"duration,omitempty"`
	Name          *string            `json:"name,omitempty"`
	OccasionLimit *int               `json:"occasion_limit,omitempty"`
	Price         *float32           `json:"price,omitempty"`
	ServiceId     *openapitypes.UUID `json:"service_id,omitempty"`
}
