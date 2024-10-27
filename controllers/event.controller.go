package controllers

import (
	"elrek-system_GO/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
	"gorm.io/gorm"
)

// MARK: GetEvents
func GetEvents(ctx *gin.Context) {
	userID, isHost := CheckAuth(ctx, false)
	if userID == "" {
		return
	}

	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")
	pageIndex, _ := strconv.Atoi(ctx.Query("pageIndex"))
	pageSize, _ := strconv.Atoi(ctx.Query("pageSize"))
	isActive := ctx.Query("isActive")

	var events []models.Event
	var query *gorm.DB

	if isHost {
		query = DB.Where("hosts.id = ?", userID).Preload("Location").Preload("Service").Preload("Hosts").Order("datetime desc")
	} else {
		// Get user's owner_id
		var user models.User
		result := DB.First(&user, "id = ?", userID)
		if result.Error != nil {
			SendMessageOnly("Could not get user details: "+result.Error.Error(), ctx, 500)
			return
		}
		query = DB.Where("hosts.id = ?", user.OwnerID).Preload("Location").Preload("Service").Preload("Hosts").Order("datetime desc")
	}

	if startDate != "" {
		query = query.Where("datetime >= ?", startDate)
	}

	if endDate != "" {
		query = query.Where("datetime <= ?", endDate)
	}

	if isActive != "" {
		query = query.Where("is_active = ?", isActive)
	} else {
		query = query.Where("is_active = ?", true)
	}

	result := query.Offset(pageIndex * pageSize).Limit(pageSize).Find(&events)
	if result.Error != nil {
		SendMessageOnly("Could not get events: "+result.Error.Error(), ctx, 500)
	}

	// TODO: Group them daily

	ctx.JSON(200, events)
}

// MARK: GetEvent
func GetEvent(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, false)
	if userID == "" {
		return
	}

	var event models.Event
	id := ctx.Param("id")

	result := DB.Preload("Location").Preload("Service").Preload("Hosts").First(&event, "id = ?", id)
	if result.Error != nil {
		SendMessageOnly("Could not get event: "+result.Error.Error(), ctx, 500)
		return
	}

	ctx.JSON(200, event)
}

// MARK: CreateEvent
func CreateEvent(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	var event models.EventCreate
	if err := ctx.ShouldBindJSON(&event); err != nil {
		SendMessageOnly("Invalid request: "+err.Error(), ctx, 400)
		return
	}

	// Check if user is the owner of the location
	var location models.Location
	result := DB.First(&location, "id = ? and user_id = ?", event.LocationID, userID)
	if result.Error != nil {
		SendMessageOnly("This location does not belong to you: "+result.Error.Error(), ctx, 403)
		return
	}

	// Check if user is the owner of the service
	var service models.Service
	result = DB.First(&service, "id = ? and user_id = ?", event.ServiceID, userID)
	if result.Error != nil {
		SendMessageOnly("This service does not belong to you: "+result.Error.Error(), ctx, 403)
		return
	}

	newEvent := models.Event{
		UserID:      openapitypes.UUID(uuid.MustParse(userID)),
		Name:        event.Name,
		Datetime:    event.Datetime,
		RootEventID: event.RootEventID,
		Comment:     event.Comment,
		Capacity:    event.Capacity,
		IsActive:    event.IsActive,
		CloseTime:   event.CloseTime,
		LocationID:  event.LocationID,
		Price:       service.Price,
		ServiceID:   event.ServiceID,
		HostIDs:     event.HostIDs,
	}

	tx := DB.Begin()

	result = tx.Create(&newEvent)
	if result.Error != nil {
		SendMessageOnly("Could not create event: "+result.Error.Error(), ctx, 500)
		tx.Rollback()
		return
	}

	tx.Commit()
	SendMessageOnly("Event created successfully", ctx, 201)
}

// MARK: UpdateEvent
func UpdateEvent(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	var event models.Event
	id := ctx.Param("id")

	result := DB.First(&event, "id = ?", id)
	if result.Error != nil {
		SendMessageOnly("Could not get event: "+result.Error.Error(), ctx, 500)
		return
	}

	var eventUpdate models.EventUpdate
	if err := ctx.ShouldBindJSON(&eventUpdate); err != nil {
		SendMessageOnly("Invalid request: "+err.Error(), ctx, 400)
		return
	}

	if eventUpdate.Name != nil {
		event.Name = *eventUpdate.Name
	}

	if eventUpdate.Datetime != nil {
		event.Datetime = *eventUpdate.Datetime
	}

	if eventUpdate.RootEventID != nil {
		event.RootEventID = eventUpdate.RootEventID
	}

	if eventUpdate.Comment != nil {
		event.Comment = eventUpdate.Comment
	}

	if eventUpdate.Capacity != nil {
		event.Capacity = eventUpdate.Capacity
	}

	if eventUpdate.IsActive != nil {
		event.IsActive = *eventUpdate.IsActive
	}

	if eventUpdate.CloseTime != nil {
		event.CloseTime = *eventUpdate.CloseTime
	}

	if eventUpdate.LocationID != nil {
		event.LocationID = eventUpdate.LocationID
	}

	if eventUpdate.Price != nil {
		event.Price = *eventUpdate.Price
	}

	if eventUpdate.ServiceID != nil {
		event.ServiceID = eventUpdate.ServiceID
	}

	if eventUpdate.HostIDs != nil {
		event.HostIDs = *eventUpdate.HostIDs
	}

	tx := DB.Begin()

	result = tx.Save(&event)
	if result.Error != nil {
		SendMessageOnly("Could not update event: "+result.Error.Error(), ctx, 500)
		tx.Rollback()
		return
	}

	tx.Commit()
	SendMessageOnly("Event updated successfully", ctx, 200)
}

// MARK: DeleteEvent
func DeleteEvent(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	id := ctx.Param("id")

	tx := DB.Begin()

	var event models.Event
	result := tx.First(&event, "id = ?", id)
	if result.Error != nil {
		SendMessageOnly("Could not get existing event: "+result.Error.Error(), ctx, 500)
		tx.Rollback()
		return
	}

	if event.UserID.String() != userID {
		SendMessageOnly("You are not allowed to delete this event", ctx, 403)
		tx.Rollback()
		return
	}

	event.IsActive = false

	result = tx.Save(&event)
	if result.Error != nil {
		SendMessageOnly("Could not delete event: "+result.Error.Error(), ctx, 500)
		tx.Rollback()
		return
	}

	tx.Commit()
	SendMessageOnly("Event was deleted successfully", ctx, 200)
}
