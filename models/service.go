package models

import "time"
import openapitypes "github.com/oapi-codegen/runtime/types"

// Service defines model for Service.
type Service struct {
	IsActive      bool              `json:"is_active,omitempty"`
	Comment       *string           `json:"comment,omitempty"`
	CreatedAt     time.Time         `json:"created_at,omitempty"`
	Id            openapitypes.UUID `json:"id,omitempty"`
	Name          string            `json:"name,omitempty"`
	OwnerId       openapitypes.UUID `json:"owner_id,omitempty"`
	PrevServiceId openapitypes.UUID `json:"prev_service_id,omitempty"`
	Price         int               `json:"price,omitempty"`
}

type ServiceWDP struct {
	IsActive      bool              `json:"is_active,omitempty"`
	Comment       *string           `json:"comment,omitempty"`
	CreatedAt     time.Time         `json:"created_at,omitempty"`
	Id            openapitypes.UUID `json:"id,omitempty"`
	Name          string            `json:"name,omitempty"`
	OwnerId       openapitypes.UUID `json:"owner_id,omitempty"`
	PrevServiceId openapitypes.UUID `json:"prev_service_id,omitempty"`
	Price         int               `json:"price,omitempty"`
	DynamicPrices *[]DynamicPrice   `json:"dynamic_prices,omitempty"`
}

// ServiceCreate defines model for ServiceCreate.
type ServiceCreate struct {
	Name                     string                      `json:"name"`
	Price                    int                         `json:"price"`
	Comment                  *string                     `json:"comment,omitempty"`
	DynamicPriceCreateUpdate *[]DynamicPriceCreateUpdate `json:"dynamic_price,omitempty"`
}

// ServiceList defines model for ServiceList.
type ServiceList = []ServiceWDP

// ServiceUpdate defines model for ServiceUpdate.
type ServiceUpdate struct {
	Name          *string                     `json:"name,omitempty"`
	Price         *int                        `json:"price,omitempty"`
	Comment       *string                     `json:"comment,omitempty"`
	DynamicPrices *[]DynamicPriceCreateUpdate `json:"dynamic_prices,omitempty"`
}
