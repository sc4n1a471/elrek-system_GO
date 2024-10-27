package controllers

import (
	"elrek-system_GO/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
)

// MARK: GetLocations
func GetLocations(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, false)
	if userID == "" {
		return
	}

	var locations []models.Location
	result := DB.Find(&locations)
	if result.Error != nil {
		SendMessageOnly("Could not get locations: "+result.Error.Error(), ctx, 500)
		return
	}

	ctx.JSON(200, locations)
}

func GetLocation(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, false)
	if userID == "" {
		return
	}

	var location models.Location
	id := ctx.Param("id")

	result := DB.First(&location, "id = ?", id)
	if result.Error != nil {
		SendMessageOnly("Could not get location: "+result.Error.Error(), ctx, 500)
		return
	}

	ctx.JSON(200, location)
}

func GetLocationEvents(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, false)
	if userID == "" {
		return
	}

	var events []models.Event
	id := ctx.Param("id")

	pageIndex, _ := strconv.Atoi(ctx.Query("pageIndex"))
	pageSize, _ := strconv.Atoi(ctx.Query("pageSize"))

	result := DB.Where("location_id = ?", id).Offset(pageIndex * pageSize).Limit(pageSize).Find(&events)
	if result.Error != nil {
		SendMessageOnly("Could not get events: "+result.Error.Error(), ctx, 500)
		return
	}

	ctx.JSON(200, events)
}

// MARK: CreateLocation
func CreateLocation(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	var locationCreate models.LocationCreate
	err := ctx.BindJSON(&locationCreate)
	if err != nil {
		SendMessageOnly("Parse error: "+err.Error(), ctx, 400)
		return
	}

	location := models.Location{
		ID:       openapitypes.UUID(uuid.New()),
		UserID:   openapitypes.UUID(uuid.MustParse(userID)),
		Name:     locationCreate.Name,
		Address:  locationCreate.Address,
		Comment:  locationCreate.Comment,
		IsActive: locationCreate.IsActive,
	}

	tx := DB.Begin()

	result := tx.Create(&location)
	if result.Error != nil {
		tx.Rollback()
		SendMessageOnly("Could not create location: "+result.Error.Error(), ctx, 500)
		return
	}

	tx.Commit()
	SendMessageOnly("Location created successfully", ctx, 201)
}

// MARK: UpdateLocation
func UpdateLocation(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	var locationUpdate models.LocationUpdate
	err := ctx.BindJSON(&locationUpdate)
	if err != nil {
		SendMessageOnly("Parse error: "+err.Error(), ctx, 400)
		return
	}

	var location models.Location
	id := ctx.Param("id")

	result := DB.First(&location, "id = ?", id)
	if result.Error != nil {
		SendMessageOnly("Could not get existing location: "+result.Error.Error(), ctx, 500)
		return
	}

	if location.UserID.String() != userID {
		SendMessageOnly("You are not allowed to update this location", ctx, 403)
		return
	}

	if locationUpdate.Name != nil {
		location.Name = *locationUpdate.Name
	}

	if locationUpdate.Address != nil {
		location.Address = locationUpdate.Address
	}

	if locationUpdate.Comment != nil {
		location.Comment = locationUpdate.Comment
	}

	if locationUpdate.IsActive != nil {
		location.IsActive = *locationUpdate.IsActive
	}

	tx := DB.Begin()
	if locationUpdate.UpdateOnly {
		result = tx.Save(&location)
	} else {
		newLocation := models.Location{
			ID:       openapitypes.UUID(uuid.New()),
			UserID:   location.UserID,
			Name:     location.Name,
			Address:  location.Address,
			Comment:  location.Comment,
			IsActive: location.IsActive,
		}
		result = tx.Create(&newLocation)
	}
	if result.Error != nil {
		tx.Rollback()
		SendMessageOnly("Could not update location: "+result.Error.Error(), ctx, 500)
		return
	}

	tx.Commit()
	SendMessageOnly("Location was updated successfully", ctx, 200)
}

// MARK: DeleteLocation
func DeleteLocation(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	id := ctx.Param("id")

	tx := DB.Begin()

	var location models.Location
	result := tx.First(&location, "id = ?", id)
	if result.Error != nil {
		SendMessageOnly("Could not get existing location: "+result.Error.Error(), ctx, 500)
		tx.Rollback()
		return
	}

	if location.UserID.String() != userID {
		SendMessageOnly("You are not allowed to delete this location", ctx, 403)
		tx.Rollback()
		return
	}

	location.IsActive = false

	result = tx.Save(&location)
	if result.Error != nil {
		SendMessageOnly("Could not delete location: "+result.Error.Error(), ctx, 500)
		tx.Rollback()
		return
	}

	tx.Commit()
	SendMessageOnly("Location deleted successfully", ctx, 200)
}
