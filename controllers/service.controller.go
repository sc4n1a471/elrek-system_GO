package controllers

import (
	"elrek-system_GO/models"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
	"gorm.io/gorm"
)

// MARK: GET
func GetServices(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	var services []models.Service
	result := DB.Where("user_id = ? and is_active = ?", userID, true).Preload("DynamicPrices").Order("created_at desc").Find(&services)
	if result.Error != nil {
		SendMessageOnly("Could not get services: "+result.Error.Error(), ctx, 500)
		return
	}

	ctx.JSON(200, services)
}

func GetService(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, false)
	if userID == "" {
		return
	}

	var service models.Service
	id := ctx.Param("id")

	result := DB.Preload("DynamicPrices").First(&service, "id = ?", id)
	if result.Error != nil {
		SendMessageOnly("Could not get service: "+result.Error.Error(), ctx, 500)
		return
	}

	ctx.JSON(200, service)
}

func GetPrevServices(rootServiceID openapitypes.UUID) ([]models.Service, error) {
	var services []models.Service
	currentServiceID := rootServiceID
	defaultUUID := openapitypes.UUID{}

	for currentServiceID != defaultUUID {
		var service models.Service
		result := DB.Where("id = ?", currentServiceID).Find(&service)
		if result.Error != nil {
			return nil, fmt.Errorf("Could not get service: " + result.Error.Error())
		}

		if rootServiceID != service.ID {
			services = append(services, service)
		}

		currentServiceID = service.PrevServiceID
	}

	return services, nil
}

// MARK: CREATE
func CreateServiceWrapper(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	tx := DB.Begin()
	result := createService(ctx, tx, userID, openapitypes.UUID{})
	if !result.Success {
		SendMessageOnly(result.Message, ctx, 500)
		tx.Rollback()
		return
	}
	tx.Commit()
	SendMessageOnly("Service was created successfully", ctx, 201)
}

func createService(ctx *gin.Context, tx *gorm.DB, userID string, prevServiceID openapitypes.UUID) ActionResponse {
	var serviceCreate models.ServiceCreate
	err := ctx.BindJSON(&serviceCreate)
	if err != nil {
		return ActionResponse{false, "Could not create service: " + err.Error()}
	}

	var service models.Service
	service.UserID = openapitypes.UUID(uuid.MustParse(userID))
	service.ID = openapitypes.UUID(uuid.New())
	service.Name = serviceCreate.Name
	service.Price = serviceCreate.Price
	service.Comment = serviceCreate.Comment
	service.IsActive = true
	service.PrevServiceID = prevServiceID

	result := tx.Create(&service)
	if result.Error != nil {
		return ActionResponse{false, "Could not create service: " + result.Error.Error()}
	}

	if serviceCreate.DynamicPriceCreateUpdate != nil {
		dpResult := createDynamicPrices(tx, *serviceCreate.DynamicPriceCreateUpdate, userID, service.ID)
		if !dpResult.Success {
			return dpResult
		}
	}

	return ActionResponse{true, "SUCCESS"}
}

// MARK: UPDATE
func UpdateService(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
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

	if service.UserID.String() != userID {
		SendMessageOnly("You are not allowed to update this service", ctx, 403)
		return
	}

	tx := DB.Begin()

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
		// MARK: Create new service instead of updating
		deleteResult := deleteService(tx, service)
		if !deleteResult.Success {
			SendMessageOnly(deleteResult.Message, ctx, 500)
			tx.Rollback()
			return
		}

		service.PrevServiceID = service.ID
		service.ID = openapitypes.UUID(uuid.New())

		result = tx.Create(&service)

		if serviceUpdate.DynamicPrices != nil {
			dpResult := createDynamicPrices(tx, *serviceUpdate.DynamicPrices, userID, service.ID)
			if !dpResult.Success {
				SendMessageOnly(dpResult.Message, ctx, 500)
				tx.Rollback()
				return
			}
		} else {
			dynamicPrices, dpResult := getDynamicPrices(service.PrevServiceID)
			if !dpResult.Success {
				SendMessageOnly(dpResult.Message, ctx, 500)
				tx.Rollback()
				return
			}

			dpResult = createDynamicPricesFromFullData(tx, dynamicPrices, userID, service.ID)
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

// MARK: DELETE
func DeleteServiceWrapper(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
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

func DeleteService(tx *gorm.DB, serviceID string) ActionResponse {
	var service models.Service

	result := tx.First(&service, "id = ?", serviceID)
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

	dpResult := deleteDynamicPricesByServiceID(tx, service.ID)
	if !dpResult.Success {
		return dpResult
	}
	return ActionResponse{true, "SUCCESS"}
}
