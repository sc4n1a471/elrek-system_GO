package models

import (
	openapitypes "github.com/oapi-codegen/runtime/types"
	"sort"
)

// DynamicPrice defines model for Dynamic_price.
type DynamicPrice struct {
	Active    bool              `json:"active,omitempty"`
	Attendees int32             `json:"attendees,omitempty"`
	Id        openapitypes.UUID `json:"id,omitempty"`
	OwnerId   openapitypes.UUID `json:"owner_id,omitempty"`
	Price     float32           `json:"price,omitempty"`
	ServiceId openapitypes.UUID `json:"service_id,omitempty"`
}

// DynamicPriceCreateUpdate defines model for Dynamic_priceCreateUpdate.
// Updating dynamic prices uses the same model
type DynamicPriceCreateUpdate struct {
	Attendees int32             `json:"attendees"`
	OwnerId   openapitypes.UUID `json:"owner_id"`
	Price     float32           `json:"price"`
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

// Compares two DynamicPrice in terms of Attendees and Price
func (d DynamicPrice) isAPEqual(other DynamicPrice) bool {
	return d.Attendees == other.Attendees &&
		d.Price == other.Price
}

// AreDPsEqualInAttPri compares two arrays of DynamicPrice in terms of Attendees and Price
func AreDPsEqualInAttPri(dp1, dp2 []DynamicPrice) bool {
	if len(dp1) != len(dp2) {
		return false
	}

	// sorts both arrays by attendees descending
	sort.Slice(dp1, func(i, j int) bool {
		return dp1[i].Attendees > dp1[j].Attendees
	})
	sort.Slice(dp2, func(i, j int) bool {
		return dp2[i].Attendees > dp2[j].Attendees
	})

	for i := range dp1 {
		if !dp1[i].isAPEqual(dp2[i]) {
			return false
		}
	}
	return true
}
