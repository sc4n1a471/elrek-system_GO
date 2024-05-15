package models

import (
	"time"

	openapitypes "github.com/oapi-codegen/runtime/types"
)

// Income defines model for Income.
type Income struct {
	IsActive     bool               `json:"is_active,omitempty"`
	Amount       int                `json:"amount,omitempty"`
	Comment      *string            `json:"comment,omitempty"`
	CreatedAt    time.Time          `json:"created_at,omitempty"`
	ID           openapitypes.UUID  `json:"id,omitempty"`
	UserID       openapitypes.UUID  `json:"user_id,omitempty"`
	ActivePassID *openapitypes.UUID `json:"active_pass_id,omitempty" gorm:"size:255"`
	ActivePass   *ActivePass        `json:"active_pass,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ServiceID    *openapitypes.UUID `json:"service_id,omitempty" gorm:"size:255"`
	Service      *Service           `json:"service,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UpdatedAt    time.Time          `json:"updated_at,omitempty"`
	PayerID      openapitypes.UUID  `json:"payer_id,omitempty" gorm:"size:255"`
	User         User               `json:"user,omitempty" gorm:"foreignKey:PayerID;;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name         *string            `json:"name,omitempty"`
	IsPaid       bool               `json:"is_paid,omitempty"`
}

// IncomeCreate defines model for IncomeCreate.
type IncomeCreate struct {
	Name         *string            `json:"name"`
	Amount       int                `json:"amount"`
	Comment      *string            `json:"comment,omitempty"`
	UserID       openapitypes.UUID  `json:"user_id"`
	ActivePassID *openapitypes.UUID `json:"active_pass_id,omitempty"`
	ServiceID    *openapitypes.UUID `json:"service_id,omitempty"`
	PayerID      openapitypes.UUID  `json:"payer_id"`
	CreatedAt    *time.Time         `json:"created_at,omitempty"`
	IsPaid       *bool              `json:"is_paid,omitempty"`
}

// IncomeCreateMultipleUsers defines model for IncomeCreateMultipleUsers.
type IncomeCreateMultipleUsers struct {
	PayerIDs      []openapitypes.UUID  `json:"payerIDs"`
	ServiceIDs    *[]openapitypes.UUID `json:"serviceIDs"`
	ActivePassIDs *[]openapitypes.UUID `json:"activePassIDs"`
	Comment       *string              `json:"comment,omitempty"`
	CreatedAt     *time.Time           `json:"createdAt,omitempty"`
	Amount        *int                 `json:"amount,omitempty"`
	IsPaid        *bool                `json:"isPaid,omitempty"`
	Name          *string              `json:"name,omitempty"`
}

// IncomeUpdate defines model for IncomeUpdate.
type IncomeUpdate struct {
	Name         *string            `json:"name,omitempty"`
	Amount       *int               `json:"amount,omitempty"`
	Comment      *string            `json:"comment,omitempty"`
	ActivePassID *openapitypes.UUID `json:"active_pass_id,omitempty"`
	ServiceID    *openapitypes.UUID `json:"service_id,omitempty"`
	PayerID      *openapitypes.UUID `json:"payer_id,omitempty"`
	CreatedAt    *time.Time         `json:"created_at,omitempty"`
	IsPaid       *bool              `json:"is_paid,omitempty"`
}
