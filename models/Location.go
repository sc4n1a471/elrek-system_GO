package models

import (
	"time"

	openapitypes "github.com/oapi-codegen/runtime/types"
)

type Location struct {
	ID        openapitypes.UUID `json:"id"`
	UserID    openapitypes.UUID `json:"userID"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
	Name      string            `json:"name"`
	Address   *string           `json:"address"`
	Comment   *string           `json:"comment"`
	IsActive  bool              `json:"isActive"`
}

type LocationCreate struct {
	Name     string  `json:"name"`
	Address  *string `json:"address"`
	Comment  *string `json:"comment"`
	IsActive bool    `json:"isActive"`
}

type LocationUpdate struct {
	Name       *string `json:"name,omitempty"`
	Address    *string `json:"address,omitempty"`
	Comment    *string `json:"comment,omitempty"`
	IsActive   *bool   `json:"isActive,omitempty"`
	UpdateOnly bool    `json:"updateOnly"`
}
