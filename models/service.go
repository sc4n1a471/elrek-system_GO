package models

import "time"
import openapitypes "github.com/oapi-codegen/runtime/types"

// Service defines model for Service.
type Service struct {
	Active        bool              `json:"active,omitempty"`
	Comment       *string           `json:"comment,omitempty"`
	CreatedAt     time.Time         `json:"created_at,omitempty"`
	Id            openapitypes.UUID `json:"id,omitempty"`
	Name          string            `json:"name,omitempty"`
	OwnerId       openapitypes.UUID `json:"owner_id,omitempty"`
	PrevServiceId string            `json:"prev_service_id,omitempty"`
	Price         float32           `json:"price,omitempty"`
}

// ServiceCreate defines model for ServiceCreate.
type ServiceCreate struct {
	Name                     string                     `json:"name"`
	Price                    float32                    `json:"price"`
	Comment                  *string                    `json:"comment,omitempty"`
	DynamicPriceCreateUpdate []DynamicPriceCreateUpdate `json:"-"`
}

// ServiceList defines model for ServiceList.
type ServiceList = []Service

// ServiceUpdate defines model for ServiceUpdate.
type ServiceUpdate struct {
	Name          *string                     `json:"name,omitempty"`
	Price         *float32                    `json:"price,omitempty"`
	Comment       *string                     `json:"comment,omitempty"`
	DynamicPrices *[]DynamicPriceCreateUpdate `json:"-"`
}
