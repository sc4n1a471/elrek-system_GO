package models

import "time"
import openapitypes "github.com/oapi-codegen/runtime/types"

// PassInUse defines model for PassInUse.
type PassInUse struct {
	Active     *bool              `json:"active,omitempty"`
	Comment    *string            `json:"comment,omitempty"`
	CreatedAt  *time.Time         `json:"created_at,omitempty"`
	Id         *openapitypes.UUID `json:"id,omitempty"`
	Occasions  *int               `json:"occasions,omitempty"`
	OwnerId    *openapitypes.UUID `json:"owner_id,omitempty"`
	PassId     *openapitypes.UUID `json:"pass_id,omitempty"`
	UserId     *openapitypes.UUID `json:"user_id,omitempty"`
	ValidFrom  *time.Time         `json:"valid_from,omitempty"`
	ValidUntil *time.Time         `json:"valid_until,omitempty"`
}

// PassInUseCreate defines model for PassInUseCreate.
type PassInUseCreate struct {
	Comment    *string           `json:"comment,omitempty"`
	OwnerId    openapitypes.UUID `json:"owner_id"`
	PassId     openapitypes.UUID `json:"pass_id"`
	UserId     openapitypes.UUID `json:"user_id"`
	ValidFrom  time.Time         `json:"valid_from"`
	ValidUntil time.Time         `json:"valid_until"`
}

// PassInUseList defines model for PassInUseList.
type PassInUseList = []PassInUse

// PassInUseListErrorResponse defines model for PassInUseListErrorResponse.
type PassInUseListErrorResponse struct {
	Message     *string        `json:"message,omitempty"`
	PassesInUse *[]interface{} `json:"passes_in_use,omitempty"`
}

// PassInUseListSuccessResponse defines model for PassInUseListSuccessResponse.
type PassInUseListSuccessResponse struct {
	Message     *string        `json:"message,omitempty"`
	PassesInUse *PassInUseList `json:"passes_in_use,omitempty"`
}

// PassInUseUpdate defines model for PassInUseUpdate.
type PassInUseUpdate struct {
	Comment    *string    `json:"comment,omitempty"`
	Occasions  *int       `json:"occasions,omitempty"`
	ValidFrom  *time.Time `json:"valid_from,omitempty"`
	ValidUntil *time.Time `json:"valid_until,omitempty"`
}
