package controllers

import (
	"elrek-system_GO/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
	"gorm.io/gorm"
)

func GetServices(ctx *gin.Context) {
	userId, _ := CheckAuth(ctx, true)
	if userId == "" {
		return
	}

	var services []models.Service
	result := DB.Where("owner_id = ?", userId).Find(&services)
	if result.Error != nil {
		SendMessageOnly("Could not get services: "+result.Error.Error(), ctx, 500)
		return
	}

	ctx.JSON(200, services)
}

func GetService(ctx *gin.Context) {
	userId, _ := CheckAuth(ctx, false)
	if userId == "" {
		return
	}

	var service models.Service
	id := ctx.Param("id")

	result := DB.First(&service, "id = ?", id)
	if result.Error != nil {
		SendMessageOnly("Could not get service: "+result.Error.Error(), ctx, 500)
		return
	}

	ctx.JSON(200, service)
}

func CreateServiceWrapper(ctx *gin.Context) {
	userId, _ := CheckAuth(ctx, true)
	if userId == "" {
		return
	}

	tx := DB.Begin()
	result := createService(ctx, tx, userId, "ORIGINAL")
	if !result.Success {
		SendMessageOnly(result.Message, ctx, 500)
		tx.Rollback()
		return
	}
	tx.Commit()
	SendMessageOnly("Service was created successfully", ctx, 200)
}

func createService(ctx *gin.Context, tx *gorm.DB, userId string, prevServiceId string) ActionResponse {
	var serviceCreate models.ServiceCreate
	err := ctx.BindJSON(&serviceCreate)
	if err != nil {
		return ActionResponse{false, "Could not create service: " + err.Error()}
	}

	var service models.Service
	service.OwnerId = openapitypes.UUID(uuid.MustParse(userId))
	service.Id = openapitypes.UUID(uuid.New())
	service.Name = serviceCreate.Name
	service.Price = serviceCreate.Price
	service.Comment = serviceCreate.Comment
	service.Active = true
	service.PrevServiceId = prevServiceId

	result := tx.Create(&service)
	if result.Error != nil {
		return ActionResponse{false, "Could not create service: " + result.Error.Error()}
	}

	dpResult := createDynamicPrices(tx, serviceCreate.DynamicPriceCreateUpdate, userId, service.Id)
	if !dpResult.Success {
		return dpResult
	}

	return ActionResponse{true, "SUCCESS"}
}

func UpdateService(ctx *gin.Context) {
	userId, _ := CheckAuth(ctx, true)
	if userId == "" {
		return
	}

	var serviceUpdate models.ServiceUpdate
	err := ctx.BindJSON(&serviceUpdate)
	if err != nil {
		SendMessageOnly("Parse error: "+err.Error(), ctx, 400)
		return
	}

	var service models.Service
	id := ctx.Param("id")

	result := DB.First(&service, "id = ?", id)
	if result.Error != nil {
		SendMessageOnly("Could not get existing service: "+result.Error.Error(), ctx, 500)
		return
	}

	tx := DB.Begin()
	deleteResult := deleteService(tx, service)
	if !deleteResult.Success {
		SendMessageOnly(deleteResult.Message, ctx, 500)
		tx.Rollback()
		return
	}

	update := true
	if serviceUpdate.Name != nil {
		service.Name = *serviceUpdate.Name
		update = false
	}

	if serviceUpdate.Price != nil {
		service.Price = *serviceUpdate.Price
		update = false
	}

	if serviceUpdate.DynamicPrices != nil {
		update = false
	}

	if serviceUpdate.Comment != nil {
		service.Comment = serviceUpdate.Comment
	}

	if update {
		result = tx.Save(&service)
	} else {
		service.PrevServiceId = service.Id.String()
		service.Id = openapitypes.UUID(uuid.New())

		result = tx.Create(&service)

		dpResult := createDynamicPrices(tx, *serviceUpdate.DynamicPrices, userId, service.Id)
		if !dpResult.Success {
			SendMessageOnly(dpResult.Message, ctx, 500)
			tx.Rollback()
			return
		}
	}

	if result.Error != nil {
		tx.Rollback()
		SendMessageOnly("Could not update service: "+result.Error.Error(), ctx, 500)
		return
	}

	tx.Commit()
	SendMessageOnly("Service was updated successfully", ctx, 200)
}

func DeleteServiceWrapper(ctx *gin.Context) {
	userId, _ := CheckAuth(ctx, true)
	if userId == "" {
		return
	}

	id := ctx.Param("id")
	tx := DB.Begin()

	result := DeleteService(tx, id)
	if !result.Success {
		SendMessageOnly(result.Message, ctx, 500)
		tx.Rollback()
		return
	}

	tx.Commit()
	SendMessageOnly("Service was deleted successfully", ctx, 200)
}

func DeleteService(tx *gorm.DB, serviceId string) ActionResponse {
	var service models.Service

	result := DB.First(&service, "id = ?", serviceId)
	if result.Error != nil {
		return ActionResponse{false, "Could not get existing service: " + result.Error.Error()}
	}

	return deleteService(tx, service)
}

func deleteService(tx *gorm.DB, service models.Service) ActionResponse {
	service.Active = false
	result := tx.Save(&service)

	if result.Error != nil {
		return ActionResponse{false, "Could not delete service: " + result.Error.Error()}
	}

	dpResult := deleteDynamicPricesByServiceId(tx, service.Id)
	if !dpResult.Success {
		return dpResult
	}
	return ActionResponse{true, "SUCCESS"}
}
