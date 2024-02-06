package controllers

import (
	"elrek-system_GO/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
	"gorm.io/gorm"
)

func getDPToService(services *[]models.Service) []models.ServiceWDP {
	var servicesWDP []models.ServiceWDP
	for _, service := range *services {
		var serviceWDP models.ServiceWDP
		serviceWDP.Id = service.Id
		serviceWDP.Name = service.Name
		serviceWDP.Price = service.Price
		serviceWDP.OwnerId = service.OwnerId
		serviceWDP.Comment = service.Comment
		serviceWDP.IsActive = service.IsActive
		serviceWDP.CreatedAt = service.CreatedAt
		serviceWDP.PrevServiceId = service.PrevServiceId

		dynamicPrices, dpResult := GetDynamicPrices(service.Id)
		if !dpResult.Success {
			fmt.Println(dpResult.Message)
			continue
		}
		serviceWDP.DynamicPrices = &dynamicPrices

		servicesWDP = append(servicesWDP, serviceWDP)
	}
	return servicesWDP
}

func GetServices(ctx *gin.Context) {
	userId, _ := CheckAuth(ctx, true)
	if userId == "" {
		return
	}

	var services []models.Service
	result := DB.Where("owner_id = ? and is_active = true", userId).Find(&services)
	if result.Error != nil {
		SendMessageOnly("Could not get services: "+result.Error.Error(), ctx, 500)
		return
	}

	servicesWDP := getDPToService(&services)

	ctx.JSON(200, servicesWDP)
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

	serviceWDP := getDPToService(&[]models.Service{service})[0]

	ctx.JSON(200, serviceWDP)
}

func CreateServiceWrapper(ctx *gin.Context) {
	userId, _ := CheckAuth(ctx, true)
	if userId == "" {
		return
	}

	tx := DB.Begin()
	result := createService(ctx, tx, userId, openapitypes.UUID{})
	if !result.Success {
		SendMessageOnly(result.Message, ctx, 500)
		tx.Rollback()
		return
	}
	tx.Commit()
	SendMessageOnly("Service was created successfully", ctx, 201)
}

func createService(ctx *gin.Context, tx *gorm.DB, userId string, prevServiceId openapitypes.UUID) ActionResponse {
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
	service.IsActive = true
	service.PrevServiceId = prevServiceId

	result := tx.Create(&service)
	if result.Error != nil {
		return ActionResponse{false, "Could not create service: " + result.Error.Error()}
	}

	if serviceCreate.DynamicPriceCreateUpdate != nil {
		dpResult := createDynamicPrices(tx, *serviceCreate.DynamicPriceCreateUpdate, userId, service.Id)
		if !dpResult.Success {
			return dpResult
		}
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

	fmt.Println(serviceUpdate.Name, serviceUpdate.Price, serviceUpdate.DynamicPrices, serviceUpdate.Comment)

	update := true
	if serviceUpdate.Name != nil {
		fmt.Println("New name: ", *serviceUpdate.Name)
		service.Name = *serviceUpdate.Name
		update = false
	}

	if serviceUpdate.Price != nil {
		fmt.Println("New price: ", *serviceUpdate.Price)
		service.Price = *serviceUpdate.Price
		update = false
	}

	if serviceUpdate.DynamicPrices != nil {
		fmt.Println("New dynamic prices: ", *serviceUpdate.DynamicPrices)
		update = false
	}

	if serviceUpdate.Comment != nil {
		fmt.Println("New comment: ", *serviceUpdate.Comment)
		service.Comment = serviceUpdate.Comment
	}

	if update {
		result = tx.Save(&service)
	} else {
		deleteResult := deleteService(tx, service)
		if !deleteResult.Success {
			SendMessageOnly(deleteResult.Message, ctx, 500)
			tx.Rollback()
			return
		}

		service.PrevServiceId = service.Id
		service.Id = openapitypes.UUID(uuid.New())

		result = tx.Create(&service)

		if serviceUpdate.DynamicPrices != nil {
			dpResult := createDynamicPrices(tx, *serviceUpdate.DynamicPrices, userId, service.Id)
			if !dpResult.Success {
				SendMessageOnly(dpResult.Message, ctx, 500)
				tx.Rollback()
				return
			}
		} else {
			dynamicPrices, dpResult := GetDynamicPrices(service.PrevServiceId)
			if !dpResult.Success {
				SendMessageOnly(dpResult.Message, ctx, 500)
				tx.Rollback()
				return
			}

			dpResult = createDynamicPricesFromFullData(tx, dynamicPrices, userId, service.Id)
			if !dpResult.Success {
				SendMessageOnly(dpResult.Message, ctx, 500)
				tx.Rollback()
				return
			}
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

	result := tx.First(&service, "id = ?", serviceId)
	if result.Error != nil {
		return ActionResponse{false, "Could not get existing service: " + result.Error.Error()}
	}

	return deleteService(tx, service)
}

func deleteService(tx *gorm.DB, service models.Service) ActionResponse {
	service.IsActive = false
	result := tx.Save(&service)

	if result.Error != nil {
		return ActionResponse{false, "Could not delete service: " + result.Error.Error()}
	} else {
		fmt.Println("Service part deleted")
	}

	dpResult := deleteDynamicPricesByServiceId(tx, service.Id)
	if !dpResult.Success {
		return dpResult
	}
	return ActionResponse{true, "SUCCESS"}
}
