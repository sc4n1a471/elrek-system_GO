package models

import "time"
import openapitypes "github.com/oapi-codegen/runtime/types"

// Income defines model for Income.
type Income struct {
	Active      *bool              `json:"active,omitempty"`
	Amount      *float32           `json:"amount,omitempty"`
	Comment     *string            `json:"comment,omitempty"`
	CreatedAt   *time.Time         `json:"created_at,omitempty"`
	ID          *openapitypes.UUID `json:"id,omitempty"`
	OwnerID     *openapitypes.UUID `json:"owner_id,omitempty"`
	PassInUseID *openapitypes.UUID `json:"pass_in_use_id,omitempty"`
	ServiceID   *openapitypes.UUID `json:"service_id,omitempty"`
	SumID       *openapitypes.UUID `json:"sum_id,omitempty"`
	UpdatedAt   *time.Time         `json:"updated_at,omitempty"`
	UserID      *openapitypes.UUID `json:"user_id,omitempty"`
}

// IncomeCreate defines model for IncomeCreate.
type IncomeCreate struct {
	Amount      float32            `json:"amount"`
	Comment     *string            `json:"comment,omitempty"`
	OwnerID     openapitypes.UUID  `json:"owner_id"`
	PassInUseID *openapitypes.UUID `json:"pass_in_use_id,omitempty"`
	ServiceID   *openapitypes.UUID `json:"service_id,omitempty"`
	SumID       *openapitypes.UUID `json:"sum_id,omitempty"`
	UserID      openapitypes.UUID  `json:"user_id"`
}

// IncomeList defines model for IncomeList.
type IncomeList = []Income

// IncomeUpdate defines model for IncomeUpdate.
type IncomeUpdate struct {
	Amount      *float32           `json:"amount,omitempty"`
	Comment     *string            `json:"comment,omitempty"`
	PassInUseID *openapitypes.UUID `json:"pass_in_use_id,omitempty"`
	ServiceID   *openapitypes.UUID `json:"service_id,omitempty"`
	SumID       *openapitypes.UUID `json:"sum_id,omitempty"`
	UserID      *openapitypes.UUID `json:"user_id,omitempty"`
}
