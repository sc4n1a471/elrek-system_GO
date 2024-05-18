package models

import (
	"time"

	openapitypes "github.com/oapi-codegen/runtime/types"
)

// Service defines model for Service.
type Service struct {
	IsActive      bool              `json:"isActive,omitempty"`
	Comment       *string           `json:"comment,omitempty"`
	CreatedAt     time.Time         `json:"createdAt,omitempty"`
	ID            openapitypes.UUID `json:"id,omitempty" gorm:"size:255"`
	Name          string            `json:"name,omitempty"`
	UserID        openapitypes.UUID `json:"userID" gorm:"size:255"`
	PrevServiceID openapitypes.UUID `json:"prevServiceID,omitempty"`
	Price         int               `json:"price,omitempty"`
	DynamicPrices *[]DynamicPrice   `json:"dynamic-prices,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// ServiceCreate defines model for ServiceCreate.
type ServiceCreate struct {
	Name                     string                      `json:"name"`
	Price                    int                         `json:"price"`
	Comment                  *string                     `json:"comment,omitempty"`
	DynamicPriceCreateUpdate *[]DynamicPriceCreateUpdate `json:"dynamicPrices,omitempty"`
}

// ServiceList defines model for ServiceList.
type ServiceList = []Service

// ServiceUpdate defines model for ServiceUpdate.
type ServiceUpdate struct {
	Name          *string                     `json:"name,omitempty"`
	Price         *int                        `json:"price,omitempty"`
	Comment       *string                     `json:"comment,omitempty"`
	DynamicPrices *[]DynamicPriceCreateUpdate `json:"dynamicPrices,omitempty"`
}
