package models

import openapitypes "github.com/oapi-codegen/runtime/types"

// DynamicPrice defines model for Dynamic_price.
type DynamicPrice struct {
	Active    *bool              `json:"active,omitempty"`
	Attendees *int32             `json:"attendees,omitempty"`
	Id        *openapitypes.UUID `json:"id,omitempty"`
	OwnerId   *openapitypes.UUID `json:"owner_id,omitempty"`
	Price     *float32           `json:"price,omitempty"`
	ServiceId *openapitypes.UUID `json:"service_id,omitempty"`
}

// DynamicPriceCreate defines model for Dynamic_priceCreate.
type DynamicPriceCreate struct {
	Attendees string            `json:"attendees"`
	OwnerId   openapitypes.UUID `json:"owner_id"`
	Price     float32           `json:"price"`
	ServiceId openapitypes.UUID `json:"service_id"`
}

// DynamicPriceList defines model for Dynamic_priceList.
type DynamicPriceList = []DynamicPrice

// DynamicPriceListErrorResponse defines model for Dynamic_priceListErrorResponse.
type DynamicPriceListErrorResponse struct {
	DynamicPrices *[]interface{} `json:"dynamic_prices,omitempty"`
	Message       *string        `json:"message,omitempty"`
}

// DynamicPriceListSuccessResponse defines model for Dynamic_priceListSuccessResponse.
type DynamicPriceListSuccessResponse struct {
	DynamicPrices *DynamicPriceList `json:"dynamic_prices,omitempty"`
	Message       *string           `json:"message,omitempty"`
}

// DynamicPriceUpdate defines model for Dynamic_priceUpdate.
type DynamicPriceUpdate struct {
	Attendees *string            `json:"attendees,omitempty"`
	Price     *float32           `json:"price,omitempty"`
	ServiceId *openapitypes.UUID `json:"service_id,omitempty"`
}
