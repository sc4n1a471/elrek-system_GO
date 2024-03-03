package models

import "time"
import openapitypes "github.com/oapi-codegen/runtime/types"

// Income defines model for Income.
type Income struct {
	IsActive    bool               `json:"is_active,omitempty"`
	Amount      int                `json:"amount,omitempty"`
	Comment     *string            `json:"comment,omitempty"`
	CreatedAt   time.Time          `json:"created_at,omitempty"`
	ID          openapitypes.UUID  `json:"id,omitempty"`
	UserID      openapitypes.UUID  `json:"user_id,omitempty"`
	PassInUseID *openapitypes.UUID `json:"pass_in_use_id,omitempty" gorm:"size:255"`
	PassInUse   *PassInUse         `json:"pass_in_use,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ServiceID   *openapitypes.UUID `json:"service_id,omitempty" gorm:"size:255"`
	Service     *Service           `json:"service,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UpdatedAt   time.Time          `json:"updated_at,omitempty"`
	PayerID     openapitypes.UUID  `json:"payer_id,omitempty"`
	Name        *string            `json:"name,omitempty"`
	IsPaid      bool               `json:"is_paid,omitempty"`
}

// IncomeCreate defines model for IncomeCreate.
type IncomeCreate struct {
	Name        *string            `json:"name"`
	Amount      int                `json:"amount"`
	Comment     *string            `json:"comment,omitempty"`
	UserID      openapitypes.UUID  `json:"user_id"`
	PassInUseID *openapitypes.UUID `json:"pass_in_use_id,omitempty"`
	ServiceID   *openapitypes.UUID `json:"service_id,omitempty"`
	PayerID     openapitypes.UUID  `json:"payer_id"`
	CreatedAt   *time.Time         `json:"created_at,omitempty"`
	IsPaid      *bool              `json:"is_paid,omitempty"`
}

// IncomeCreateMultipleUsers defines model for IncomeCreateMultipleUsers.
type IncomeCreateMultipleUsers struct {
	PayerIDs     []openapitypes.UUID  `json:"payer_ids"`
	ServiceIDs   *[]openapitypes.UUID `json:"service_ids"`
	PassInUseIDs *[]openapitypes.UUID `json:"pass_in_use_ids"`
	Comment      *string              `json:"comment,omitempty"`
	CreatedAt    *time.Time           `json:"created_at,omitempty"`
	Amount       *int                 `json:"amount,omitempty"`
	IsPaid       *bool                `json:"is_paid,omitempty"`
	Name         *string              `json:"name,omitempty"`
}

// IncomeUpdate defines model for IncomeUpdate.
type IncomeUpdate struct {
	Name        *string            `json:"name,omitempty"`
	Amount      *int               `json:"amount,omitempty"`
	Comment     *string            `json:"comment,omitempty"`
	PassInUseID *openapitypes.UUID `json:"pass_in_use_id,omitempty"`
	ServiceID   *openapitypes.UUID `json:"service_id,omitempty"`
	PayerID     *openapitypes.UUID `json:"payer_id,omitempty"`
	CreatedAt   *time.Time         `json:"created_at,omitempty"`
	IsPaid      *bool              `json:"is_paid,omitempty"`
}
