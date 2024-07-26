package models

import (
	"time"

	openapitypes "github.com/oapi-codegen/runtime/types"
)

// Income defines model for Income.
type Income struct {
	IsActive     bool               `json:"isActive,omitempty"`
	Amount       int                `json:"amount"`
	Comment      *string            `json:"comment,omitempty"`
	CreatedAt    time.Time          `json:"createdAt,omitempty"`
	ID           openapitypes.UUID  `json:"id,omitempty"`
	UserID       openapitypes.UUID  `json:"userID,omitempty"`
	ActivePassID *openapitypes.UUID `json:"activePassID,omitempty" gorm:"size:255"`
	ActivePass   *ActivePass        `json:"activePass,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ServiceID    *openapitypes.UUID `json:"serviceID,omitempty" gorm:"size:255"`
	Service      *Service           `json:"service,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UpdatedAt    time.Time          `json:"updatedAt,omitempty"`
	PayerID      openapitypes.UUID  `json:"payerID,omitempty" gorm:"size:255"`
	User         User               `json:"user,omitempty" gorm:"foreignKey:PayerID;;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name         *string            `json:"name,omitempty"`
	IsPaid       bool               `json:"isPaid"`
}

// IncomeCreate defines model for IncomeCreate.
type IncomeCreate struct {
	Name         *string            `json:"name"`
	Amount       int                `json:"amount"`
	Comment      *string            `json:"comment,omitempty"`
	UserID       openapitypes.UUID  `json:"userID"`
	ActivePassID *openapitypes.UUID `json:"activePassID,omitempty"`
	ServiceID    *openapitypes.UUID `json:"serviceID,omitempty"`
	PayerID      openapitypes.UUID  `json:"payerID"`
	CreatedAt    *time.Time         `json:"createdAt,omitempty"`
	IsPaid       *bool              `json:"isPaid,omitempty"`
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
	ActivePassID *openapitypes.UUID `json:"activePassID,omitempty"`
	ServiceID    *openapitypes.UUID `json:"serviceID,omitempty"`
	PayerID      *openapitypes.UUID `json:"payerID,omitempty"`
	CreatedAt    *time.Time         `json:"createdAt,omitempty"`
	IsPaid       *bool              `json:"isPaid,omitempty"`
}

// IncomeListResponse defines model for IncomeListResponse, used only for tests.
type IncomeListResponse struct {
	Incomes      []Income `json:"incomes"`
	TotalIncomes int64    `json:"totalIncomes"`
}
