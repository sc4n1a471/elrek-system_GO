package controllers

import (
	"elrek-system_GO/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
	"gorm.io/gorm"
)

func GetDynamicPricesWrapper(ctx *gin.Context) {
	userId, _ := CheckAuth(ctx, false)
	if userId == "" {
		return
	}

	serviceId := openapitypes.UUID(uuid.MustParse(ctx.Param("id")))
	dynamicPrices, actionResponses := GetDynamicPrices(serviceId)
	if !actionResponses.Success {
		SendMessageOnly(actionResponses.Message, ctx, 500)
		return
	}
	ctx.JSON(200, dynamicPrices)
}

func GetDynamicPrices(serviceId openapitypes.UUID) ([]models.DynamicPrice, ActionResponse) {
	var dynamicPrices []models.DynamicPrice

	result := DB.Find(&dynamicPrices, "service_id = ?", serviceId)
	if result.Error != nil {
		return nil, ActionResponse{false, "Could not get dynamic prices: " + result.Error.Error()}
	}
	return dynamicPrices, ActionResponse{true, "SUCCESS"}
}

// NOT USED CURRENTLY AS ENDPOINT
//func CreateDynamicPricesWrapperEndpoint(ctx *gin.Context) {
//	userId, _ := CheckAuth(ctx, true)
//	if userId == "" {
//		return
//	}
//
//	tx := DB.Begin()
//	result := CreateDynamicPricesWrapperJson(ctx, tx, userId)
//	if !result.Success {
//		tx.Rollback()
//		SendMessageOnly(result.Message, ctx, 500)
//		return
//	}
//	tx.Commit()
//	SendMessageOnly("Dynamic prices were created successfully", ctx, 200)
//}
//
//func CreateDynamicPricesWrapperJson(ctx *gin.Context, tx *gorm.DB, userId string) ActionResponse {
//	var dynamicPrices []models.DynamicPriceCreateUpdate
//	err := ctx.BindJSON(&dynamicPrices)
//	if err != nil {
//		return ActionResponse{false, "Dynamic price parse error: " + err.Error()}
//	}
//
//	return createDynamicPrices(tx, dynamicPrices, userId)
//}

func createDynamicPrices(
	tx *gorm.DB,
	newDynamicPrices []models.DynamicPriceCreateUpdate,
	userId string,
	serviceId openapitypes.UUID) ActionResponse {

	for _, dynamicPrice := range newDynamicPrices {
		var dynamicPriceFull models.DynamicPrice
		dynamicPriceFull.ServiceId = serviceId
		dynamicPriceFull.Price = dynamicPrice.Price
		dynamicPriceFull.Attendees = dynamicPrice.Attendees
		dynamicPriceFull.OwnerId = openapitypes.UUID(uuid.MustParse(userId))
		dynamicPriceFull.Active = true
		dynamicPriceFull.Id = openapitypes.UUID(uuid.New())

		result := tx.Create(&dynamicPriceFull)
		if result.Error != nil {
			return ActionResponse{false, "Could not create dynamic price: " + result.Error.Error()}
		}
	}
	return ActionResponse{true, "SUCCESS"}
}

func createDynamicPricesFromFullData(
	tx *gorm.DB,
	existingDynamicPrices []models.DynamicPrice,
	userId string,
	serviceId openapitypes.UUID) ActionResponse {

	var strippedDynamicPrices []models.DynamicPriceCreateUpdate
	for _, dynamicPrice := range existingDynamicPrices {
		var strippedDynamicPrice models.DynamicPriceCreateUpdate
		strippedDynamicPrice.OwnerId = openapitypes.UUID(uuid.MustParse(userId))
		strippedDynamicPrice.Attendees = dynamicPrice.Attendees
		strippedDynamicPrice.Price = dynamicPrice.Price

		strippedDynamicPrices = append(strippedDynamicPrices, strippedDynamicPrice)
	}

	dpResult := createDynamicPrices(tx, strippedDynamicPrices, userId, serviceId)
	if !dpResult.Success {
		return dpResult
	}
	return ActionResponse{true, "SUCCESS"}
}

func updateDynamicPrices(
	tx *gorm.DB,
	updatableDyPrices []models.DynamicPrice,
	userId string,
	serviceId openapitypes.UUID) ActionResponse {

	var newDynamicPrices []models.DynamicPriceCreateUpdate

	for _, updatableDyPrice := range updatableDyPrices {
		var currentDynamicPrice models.DynamicPrice
		result := DB.First(&currentDynamicPrice, "id = ?", updatableDyPrice.Id)
		if result.Error != nil {
			return ActionResponse{false, "Could not get existing dynamic price: " + result.Error.Error()}
		}

		var newDynamicPrice models.DynamicPriceCreateUpdate
		newDynamicPrice.OwnerId = updatableDyPrice.OwnerId

		if updatableDyPrice.Attendees != currentDynamicPrice.Attendees {
			newDynamicPrice.Attendees = updatableDyPrice.Attendees
		}
		if updatableDyPrice.Price != currentDynamicPrice.Price {
			newDynamicPrice.Price = updatableDyPrice.Price
		}

		newDynamicPrices = append(newDynamicPrices, newDynamicPrice)
	}

	var deletableDynamicPrices []models.DynamicPrice
	var existingDynamicPrices []models.DynamicPrice
	result := DB.Find(&existingDynamicPrices, "service_id = ?", serviceId)
	if result.Error != nil {
		return ActionResponse{false, "Could not get existing dynamic prices before deletion: " + result.Error.Error()}
	}

	for _, existingDynamicPrice := range existingDynamicPrices {
		found := false
		for _, updatableDyPrice := range updatableDyPrices {
			if existingDynamicPrice.Id == updatableDyPrice.Id {
				found = true
				break
			}
		}
		if !found {
			deletableDynamicPrices = append(deletableDynamicPrices, existingDynamicPrice)
		}
	}

	deleteResult := deleteDynamicPrices(tx, deletableDynamicPrices)
	if !deleteResult.Success {
		return deleteResult
	}

	return createDynamicPrices(tx, newDynamicPrices, userId, serviceId)
}

func deleteDynamicPrices(tx *gorm.DB, deletableDynamicPrices []models.DynamicPrice) ActionResponse {
	fmt.Println("Deleting ", len(deletableDynamicPrices), " dynamic prices")
	for _, deletableDynamicPrice := range deletableDynamicPrices {
		deletableDynamicPrice.Active = false

		result := tx.Save(&deletableDynamicPrice)
		if result.Error != nil {
			return ActionResponse{false, "Could not delete dynamic price: " + result.Error.Error()}
		}
	}
	return ActionResponse{true, "SUCCESS"}
}

func deleteDynamicPricesByServiceId(tx *gorm.DB, serviceId openapitypes.UUID) ActionResponse {
	var dynamicPrices []models.DynamicPrice
	result := DB.Find(&dynamicPrices, "service_id = ?", serviceId)
	if result.Error != nil {
		return ActionResponse{false, "Could not get dynamic prices: " + result.Error.Error()}
	}

	return deleteDynamicPrices(tx, dynamicPrices)
}
