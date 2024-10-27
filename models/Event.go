package models

import (
	"time"

	openapitypes "github.com/oapi-codegen/runtime/types"
)

type Event struct {
	ID          openapitypes.UUID   `json:"id"`
	UserID      openapitypes.UUID   `json:"userID"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
	Name        string              `json:"name"`
	Datetime    time.Time           `json:"datetime"`
	RootEventID *openapitypes.UUID  `json:"rootEventID"`
	RootEvent   *Event              `json:"rootEvent,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Comment     *string             `json:"comment,omitempty"`
	Capacity    *int                `json:"capacity"`
	AttendeeIDs []openapitypes.UUID `json:"attendeeIDs"`
	Attendees   []User              `json:"attendees,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	IsActive    bool                `json:"isActive"`
	CloseTime   time.Time           `json:"closeTime"`
	LocationID  *openapitypes.UUID  `json:"locationID,omitempty"`
	Location    *Location           `json:"location,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Price       int                 `json:"price"`
	ServiceID   *openapitypes.UUID  `json:"serviceID,omitempty"`
	Service     *Service            `json:"service,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	HostIDs     []openapitypes.UUID `json:"hostIDs"`
	Hosts       []User              `json:"hosts" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type EventCreate struct {
	Name        string             `json:"name"`
	Datetime    time.Time          `json:"datetime"`
	RootEventID *openapitypes.UUID `json:"rootEventID"`
	Comment     *string            `json:"comment"`
	Capacity    *int               `json:"capacity"`
	IsActive    bool               `json:"isActive"`
	CloseTime   time.Time          `json:"closeTime"`
	LocationID  *openapitypes.UUID `json:"locationID"`
	Price       int                `json:"price"`
	ServiceID   *openapitypes.UUID `json:"serviceID"`

	HostIDs []openapitypes.UUID `json:"hostIDs"`
}

type EventUpdate struct {
	Name        *string            `json:"name,omitempty"`
	Datetime    *time.Time         `json:"datetime,omitempty"`
	RootEventID *openapitypes.UUID `json:"rootEventID,omitempty"`
	Comment     *string            `json:"comment,omitempty"`
	Capacity    *int               `json:"capacity,omitempty"`
	IsActive    *bool              `json:"isActive,omitempty"`
	CloseTime   *time.Time         `json:"closeTime,omitempty"`
	LocationID  *openapitypes.UUID `json:"locationID,omitempty"`
	Price       *int               `json:"price,omitempty"`
	ServiceID   *openapitypes.UUID `json:"serviceID,omitempty"`

	HostIDs *[]openapitypes.UUID `json:"hostIDs,omitempty"`
}
