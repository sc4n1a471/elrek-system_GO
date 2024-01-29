package models

import "time"
import openapitypes "github.com/oapi-codegen/runtime/types"

// Service defines model for Service.
type Service struct {
	Active        *bool              `json:"active,omitempty"`
	Comment       *string            `json:"comment,omitempty"`
	CreatedAt     *time.Time         `json:"created_at,omitempty"`
	Id            *openapitypes.UUID `json:"id,omitempty"`
	Name          *string            `json:"name,omitempty"`
	OwnerId       *openapitypes.UUID `json:"owner_id,omitempty"`
	PrevServiceId *openapitypes.UUID `json:"prev_service_id,omitempty"`
	Price         *float32           `json:"price,omitempty"`
}

// ServiceCreate defines model for ServiceCreate.
type ServiceCreate struct {
	Name    string            `json:"name"`
	OwnerId openapitypes.UUID `json:"owner_id"`
	Price   float32           `json:"price"`
}

// ServiceList defines model for ServiceList.
type ServiceList = []Service

// ServiceListErrorResponse defines model for ServiceListErrorResponse.
type ServiceListErrorResponse struct {
	Message  *string        `json:"message,omitempty"`
	Services *[]interface{} `json:"services,omitempty"`
}

// ServiceListSuccessResponse defines model for ServiceListSuccessResponse.
type ServiceListSuccessResponse struct {
	Message  *string      `json:"message,omitempty"`
	Services *ServiceList `json:"services,omitempty"`
}

// ServiceUpdate defines model for ServiceUpdate.
type ServiceUpdate struct {
	Name  *string  `json:"name,omitempty"`
	Price *float32 `json:"price,omitempty"`
}
