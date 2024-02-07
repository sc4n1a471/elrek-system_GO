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
	userID, _ := CheckAuth(ctx, false)
	if userID == "" {
		return
	}

	serviceID := openapitypes.UUID(uuid.MustParse(ctx.Param("id")))
	dynamicPrices, actionResponses := GetDynamicPrices(serviceID)
	if !actionResponses.Success {
		SendMessageOnly(actionResponses.Message, ctx, 500)
		return
	}
	ctx.JSON(200, dynamicPrices)
}

func GetDynamicPrices(serviceID openapitypes.UUID) ([]models.DynamicPrice, ActionResponse) {
	var dynamicPrices []models.DynamicPrice

	result := DB.Find(&dynamicPrices, "service_id = ?", serviceID)
	if result.Error != nil {
		return nil, ActionResponse{false, "Could not get dynamic prices: " + result.Error.Error()}
	}
	return dynamicPrices, ActionResponse{true, "SUCCESS"}
}

// NOT USED CURRENTLY AS ENDPOINT
//func CreateDynamicPricesWrapperEndpoint(ctx *gin.Context) {
//	userID, _ := CheckAuth(ctx, true)
//	if userID == "" {
//		return
//	}
//
//	tx := DB.Begin()
//	result := CreateDynamicPricesWrapperJson(ctx, tx, userID)
//	if !result.Success {
//		tx.Rollback()
//		SendMessageOnly(result.Message, ctx, 500)
//		return
//	}
//	tx.Commit()
//	SendMessageOnly("Dynamic prices were created successfully", ctx, 200)
//}
//
//func CreateDynamicPricesWrapperJson(ctx *gin.Context, tx *gorm.DB, userID string) ActionResponse {
//	var dynamicPrices []models.DynamicPriceCreateUpdate
//	err := ctx.BindJSON(&dynamicPrices)
//	if err != nil {
//		return ActionResponse{false, "Dynamic price parse error: " + err.Error()}
//	}
//
//	return createDynamicPrices(tx, dynamicPrices, userID)
//}

func createDynamicPrices(
	tx *gorm.DB,
	newDynamicPrices []models.DynamicPriceCreateUpdate,
	userID string,
	serviceID openapitypes.UUID) ActionResponse {

	for _, dynamicPrice := range newDynamicPrices {
		var dynamicPriceFull models.DynamicPrice
		dynamicPriceFull.ServiceID = serviceID
		dynamicPriceFull.Price = dynamicPrice.Price
		dynamicPriceFull.Attendees = dynamicPrice.Attendees
		dynamicPriceFull.UserID = openapitypes.UUID(uuid.MustParse(userID))
		dynamicPriceFull.Active = true
		dynamicPriceFull.ID = openapitypes.UUID(uuid.New())

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
	userID string,
	serviceID openapitypes.UUID) ActionResponse {

	var strippedDynamicPrices []models.DynamicPriceCreateUpdate
	for _, dynamicPrice := range existingDynamicPrices {
		var strippedDynamicPrice models.DynamicPriceCreateUpdate
		strippedDynamicPrice.OwnerID = openapitypes.UUID(uuid.MustParse(userID))
		strippedDynamicPrice.Attendees = dynamicPrice.Attendees
		strippedDynamicPrice.Price = dynamicPrice.Price

		strippedDynamicPrices = append(strippedDynamicPrices, strippedDynamicPrice)
	}

	dpResult := createDynamicPrices(tx, strippedDynamicPrices, userID, serviceID)
	if !dpResult.Success {
		return dpResult
	}
	return ActionResponse{true, "SUCCESS"}
}

func updateDynamicPrices(
	tx *gorm.DB,
	updatableDyPrices []models.DynamicPrice,
	userID string,
	serviceID openapitypes.UUID) ActionResponse {

	var newDynamicPrices []models.DynamicPriceCreateUpdate

	for _, updatableDyPrice := range updatableDyPrices {
		var currentDynamicPrice models.DynamicPrice
		result := DB.First(&currentDynamicPrice, "id = ?", updatableDyPrice.ID)
		if result.Error != nil {
			return ActionResponse{false, "Could not get existing dynamic price: " + result.Error.Error()}
		}

		var newDynamicPrice models.DynamicPriceCreateUpdate
		newDynamicPrice.OwnerID = updatableDyPrice.UserID

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
	result := DB.Find(&existingDynamicPrices, "service_id = ?", serviceID)
	if result.Error != nil {
		return ActionResponse{false, "Could not get existing dynamic prices before deletion: " + result.Error.Error()}
	}

	for _, existingDynamicPrice := range existingDynamicPrices {
		found := false
		for _, updatableDyPrice := range updatableDyPrices {
			if existingDynamicPrice.ID == updatableDyPrice.ID {
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

	return createDynamicPrices(tx, newDynamicPrices, userID, serviceID)
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

func deleteDynamicPricesByServiceID(tx *gorm.DB, serviceID openapitypes.UUID) ActionResponse {
	var dynamicPrices []models.DynamicPrice
	result := DB.Find(&dynamicPrices, "service_id = ?", serviceID)
	if result.Error != nil {
		return ActionResponse{false, "Could not get dynamic prices: " + result.Error.Error()}
	}

	return deleteDynamicPrices(tx, dynamicPrices)
}
