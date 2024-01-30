package models

import "time"
import openapitypes "github.com/oapi-codegen/runtime/types"

// Income defines model for Income.
type Income struct {
	Active      *bool              `json:"active,omitempty"`
	Amount      *float32           `json:"amount,omitempty"`
	Comment     *string            `json:"comment,omitempty"`
	CreatedAt   *time.Time         `json:"created_at,omitempty"`
	Id          *openapitypes.UUID `json:"id,omitempty"`
	OwnerId     *openapitypes.UUID `json:"owner_id,omitempty"`
	PassInUseId *openapitypes.UUID `json:"pass_in_use_id,omitempty"`
	ServiceId   *openapitypes.UUID `json:"service_id,omitempty"`
	SumId       *openapitypes.UUID `json:"sum_id,omitempty"`
	UpdatedAt   *time.Time         `json:"updated_at,omitempty"`
	UserId      *openapitypes.UUID `json:"user_id,omitempty"`
}

// IncomeCreate defines model for IncomeCreate.
type IncomeCreate struct {
	Amount      float32            `json:"amount"`
	Comment     *string            `json:"comment,omitempty"`
	OwnerId     openapitypes.UUID  `json:"owner_id"`
	PassInUseId *openapitypes.UUID `json:"pass_in_use_id,omitempty"`
	ServiceId   *openapitypes.UUID `json:"service_id,omitempty"`
	SumId       *openapitypes.UUID `json:"sum_id,omitempty"`
	UserId      openapitypes.UUID  `json:"user_id"`
}

// IncomeList defines model for IncomeList.
type IncomeList = []Income

// IncomeUpdate defines model for IncomeUpdate.
type IncomeUpdate struct {
	Amount      *float32           `json:"amount,omitempty"`
	Comment     *string            `json:"comment,omitempty"`
	PassInUseId *openapitypes.UUID `json:"pass_in_use_id,omitempty"`
	ServiceId   *openapitypes.UUID `json:"service_id,omitempty"`
	SumId       *openapitypes.UUID `json:"sum_id,omitempty"`
	UserId      *openapitypes.UUID `json:"user_id,omitempty"`
}
