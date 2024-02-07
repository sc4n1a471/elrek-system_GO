package models

import "time"
import openapitypes "github.com/oapi-codegen/runtime/types"

// PassInUse defines model for PassInUse.
type PassInUse struct {
	Active     *bool              `json:"active,omitempty"`
	Comment    *string            `json:"comment,omitempty"`
	CreatedAt  *time.Time         `json:"created_at,omitempty"`
	ID         *openapitypes.UUID `json:"id,omitempty"`
	Occasions  *int               `json:"occasions,omitempty"`
	UserID     *openapitypes.UUID `json:"user_id,omitempty" gorm:"size:255"`
	Pass       Pass               `json:"pass,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PassID     *openapitypes.UUID `json:"pass_id,omitempty" gorm:"size:255"`
	PayerID    *openapitypes.UUID `json:"payer_id,omitempty" gorm:"size:255"`
	ValidFrom  *time.Time         `json:"valid_from,omitempty"`
	ValidUntil *time.Time         `json:"valid_until,omitempty"`
}

// PassInUseCreate defines model for PassInUseCreate.
type PassInUseCreate struct {
	Comment    *string           `json:"comment,omitempty"`
	OwnerID    openapitypes.UUID `json:"owner_id"`
	PassID     openapitypes.UUID `json:"pass_id"`
	UserID     openapitypes.UUID `json:"user_id"`
	ValidFrom  time.Time         `json:"valid_from"`
	ValidUntil time.Time         `json:"valid_until"`
}

// PassInUseList defines model for PassInUseList.
type PassInUseList = []PassInUse

// PassInUseUpdate defines model for PassInUseUpdate.
type PassInUseUpdate struct {
	Comment    *string    `json:"comment,omitempty"`
	Occasions  *int       `json:"occasions,omitempty"`
	ValidFrom  *time.Time `json:"valid_from,omitempty"`
	ValidUntil *time.Time `json:"valid_until,omitempty"`
}
