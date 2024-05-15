package models

import (
	"sort"

	openapitypes "github.com/oapi-codegen/runtime/types"
)

// DynamicPrice defines model for Dynamic_price.
type DynamicPrice struct {
	Active    bool              `json:"active,omitempty"`
	Attendees int               `json:"attendees,omitempty"`
	ID        openapitypes.UUID `json:"id,omitempty"`
	UserID    openapitypes.UUID `json:"user_id,omitempty"`
	Price     int               `json:"price,omitempty"`
	ServiceID openapitypes.UUID `json:"service_id,omitempty" gorm:"size:255"`
}

// DynamicPriceCreateUpdate defines model for Dynamic_priceCreateUpdate.
// Updating dynamic prices uses the same model
type DynamicPriceCreateUpdate struct {
	Attendees int               `json:"attendees"`
	OwnerID   openapitypes.UUID `json:"user_id"`
	Price     int               `json:"price"`
}

// DynamicPriceList defines model for Dynamic_priceList.
type DynamicPriceList = []DynamicPrice

// DynamicPriceListErrorResponse defines model for Dynamic_priceListErrorResponse.
type DynamicPriceListErrorResponse struct {
	DynamicPrices *[]interface{} `json:"dynamic-prices,omitempty"`
	Message       *string        `json:"message,omitempty"`
}

// DynamicPriceListSuccessResponse defines model for Dynamic_priceListSuccessResponse.
type DynamicPriceListSuccessResponse struct {
	DynamicPrices *DynamicPriceList `json:"dynamic-prices,omitempty"`
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
