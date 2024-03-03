package controllers

import (
	"elrek-system_GO/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
	"gorm.io/gorm"
)

// ========== GET /incomes ==========

func GetIncomes(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	var incomes []models.Income
	result := DB.Where("user_id = ? and is_active = ?", userID, true).
		Preload("PassInUse").
		Preload("Service").
		Find(&incomes)
	if result.Error != nil {
		SendMessageOnly("Could not get incomes: "+result.Error.Error(), ctx, 500)
		return
	}

	// Why did I do this?
	//for _, income := range incomes {
	//	if income.PassInUseID != nil {
	//		passInUse, err := getPassInUse(*income.PassInUseID)
	//		if err != nil {
	//			SendMessageOnly("Could not get pass in use: "+err.Error(), ctx, 500)
	//			return
	//		}
	//
	//		income.PassInUse = &passInUse
	//	}
	//}

	ctx.JSON(200, incomes)
}

func GetIncome(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	var income models.Income
	id := ctx.Param("id")

	result := DB.Preload("PassInUse").
		Preload("Service").
		First(&income, "id = ?", id)
	if result.Error != nil {
		SendMessageOnly("Could not get income: "+result.Error.Error(), ctx, 500)
		return
	}

	ctx.JSON(200, income)
}

// ========== POST /incomes ==========

// CreateIncomeWrapper is a wrapper for the endpoint
func CreateIncomeWrapper(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	tx := DB.Begin()

	var incomeCreate models.IncomeCreate
	err := ctx.BindJSON(&incomeCreate)
	if err != nil {
		SendMessageOnly("Could not bind income: "+err.Error(), ctx, 400)
		return
	}

	result := createIncome(tx, incomeCreate, openapitypes.UUID(uuid.MustParse(userID)), 1)
	if !result.Success {
		SendMessageOnly(result.Message, ctx, 500)
		tx.Rollback()
		return
	}

	tx.Commit()
	SendMessageOnly("Income was created successfully", ctx, 201)
}

// CreateIncomeMultipleUsersWrapper is a wrapper for the endpoint which receives multiple users
func CreateIncomeMultipleUsersWrapper(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	var incomeCreateMultipleUsers models.IncomeCreateMultipleUsers
	err := ctx.BindJSON(&incomeCreateMultipleUsers)
	if err != nil {
		SendMessageOnly("Could not bind incomeCreateMultipleUsers: "+err.Error(), ctx, 400)
		return
	}

	tx := DB.Begin()

	var incomeCreate models.IncomeCreate
	incomeCreate.Comment = incomeCreateMultipleUsers.Comment

	if incomeCreateMultipleUsers.CreatedAt != nil {
		incomeCreate.CreatedAt = incomeCreateMultipleUsers.CreatedAt
	}

	if incomeCreateMultipleUsers.Amount != nil {
		incomeCreate.Amount = *incomeCreateMultipleUsers.Amount
	}

	if incomeCreateMultipleUsers.IsPaid != nil {
		incomeCreate.IsPaid = incomeCreateMultipleUsers.IsPaid
	}

	if incomeCreateMultipleUsers.Name != nil {
		incomeCreate.Name = incomeCreateMultipleUsers.Name
	}

	for _, payerID := range incomeCreateMultipleUsers.PayerIDs {
		incomeCreate.PayerID = payerID

		if incomeCreateMultipleUsers.ServiceIDs != nil {
			for _, serviceID := range *incomeCreateMultipleUsers.ServiceIDs {
				incomeCreate.ServiceID = &serviceID

				result := createIncome(
					tx,
					incomeCreate,
					openapitypes.UUID(uuid.MustParse(userID)),
					len(incomeCreateMultipleUsers.PayerIDs))

				if !result.Success {
					SendMessageOnly(result.Message, ctx, 500)
					tx.Rollback()
					return
				}
			}
		} else if incomeCreateMultipleUsers.PassInUseIDs != nil {
			for _, passInUseID := range *incomeCreateMultipleUsers.PassInUseIDs {
				incomeCreate.PassInUseID = &passInUseID

				result := createIncome(
					tx,
					incomeCreate,
					openapitypes.UUID(uuid.MustParse(userID)),
					len(incomeCreateMultipleUsers.PayerIDs))

				if !result.Success {
					SendMessageOnly(result.Message, ctx, 500)
					tx.Rollback()
					return
				}
			}
		} else {
			SendMessageOnly("ServiceIDs or PassInUseIDs must be provided", ctx, 400)
			tx.Rollback()
			return
		}
	}

	tx.Commit()
	SendMessageOnly("Incomes were created successfully", ctx, 201)
}

// createIncome creates an income using the provided data (incomeCreate) and the provided user ID (userID
func createIncome(tx *gorm.DB, incomeCreate models.IncomeCreate, userID openapitypes.UUID, numOfAttendees int) ActionResponse {

	// ========== Creating income model ==========
	income := models.Income{
		IsActive:    true,
		Comment:     incomeCreate.Comment,
		ID:          openapitypes.UUID(uuid.New()),
		UserID:      userID,
		PassInUseID: incomeCreate.PassInUseID,
		ServiceID:   incomeCreate.ServiceID,
		PayerID:     incomeCreate.PayerID,
		Name:        incomeCreate.Name,
	}

	if incomeCreate.CreatedAt != nil {
		income.CreatedAt = *incomeCreate.CreatedAt
	}

	if incomeCreate.IsPaid != nil {
		income.IsPaid = *incomeCreate.IsPaid
	} else {
		income.IsPaid = false
	}

	if incomeCreate.Amount != 0 {
		income.Amount = incomeCreate.Amount

		if incomeCreate.PassInUseID != nil && incomeCreate.Name == nil {
			name := "Bérlet vásárlás"
			income.Name = &name
		}
	} else {
		if incomeCreate.ServiceID != nil {
			var service models.Service
			result := tx.First(&service, "id = ?", incomeCreate.ServiceID)
			if result.Error != nil {
				return ActionResponse{
					Success: false,
					Message: "Could not get service in incomeCreation: " + result.Error.Error(),
				}
			}

			useResult := usePassInUse(tx, incomeCreate.PayerID, service.ID)
			fmt.Println("createIncome / useResult: ", useResult)
			if !useResult.Success {
				if useResult.Message == NoPIUWasFound || useResult.Message == NoValidPIUWasFound {
					income.Amount = service.Price

					var dynamicPrices []models.DynamicPrice
					dynamicPrices, dpResult := getDynamicPrices(service.ID)
					fmt.Println("createIncome / dynamicPrices: ", dynamicPrices)
					fmt.Println("createIncome / dpResult: ", dpResult)
					if !dpResult.Success {
						return ActionResponse{
							Success: false,
							Message: "Could not get dynamic prices: " + dpResult.Message,
						}
					}

					if len(dynamicPrices) > 0 {
						for _, dynamicPrice := range dynamicPrices {
							if dynamicPrice.Attendees >= numOfAttendees {
								income.Amount = dynamicPrice.Price
								fmt.Println("Dynamic price used: ", income.Amount)
								break
							}
						}
					}
				} else {
					return ActionResponse{
						Success: false,
						Message: "Could not use pass in use: " + useResult.Message,
					}
				}
			} else {
				// Has valid passInUse
				income.Amount = 0
			}
		} else {
			return ActionResponse{
				Success: false,
				Message: "Amount must be provided if pass in use or service is not provided",
			}
		}
	}

	// ========== Creating income ==========
	result := tx.Create(&income)
	if result.Error != nil {
		return ActionResponse{
			Success: false,
			Message: "Could not create income: " + result.Error.Error(),
		}
	}

	if income.PassInUseID != nil {
		var passInUse models.PassInUse
		result = tx.First(&passInUse, "id = ?", income.PassInUseID)

		err := tx.Model(&income).Association("PassInUse").Append(&passInUse)
		if err != nil {
			return ActionResponse{
				Success: false,
				Message: "Could not associate pass in use with income: " + err.Error(),
			}
		}
	}

	if income.ServiceID != nil {
		var service models.Service
		result = DB.First(&service, "id = ?", income.ServiceID)

		err := tx.Model(&income).Association("Service").Append(&service)
		if err != nil {
			return ActionResponse{
				Success: false,
				Message: "Could not associate service with income: " + err.Error(),
			}
		}
	}

	return ActionResponse{
		Success: true,
		Message: "Income was created successfully!",
	}
}

// ========== PATCH /incomes/:id ==========

func UpdateIncome(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	var income models.Income
	id := ctx.Param("id")

	result := DB.First(&income, "id = ?", id)
	if result.Error != nil {
		SendMessageOnly("Could not get existing income: "+result.Error.Error(), ctx, 500)
		return
	}

	if income.UserID.String() != userID {
		SendMessageOnly("You are not allowed to update this income", ctx, 403)
		return
	}

	var incomeUpdate models.IncomeUpdate
	err := ctx.BindJSON(&incomeUpdate)
	if err != nil {
		SendMessageOnly("Could not bind income: "+err.Error(), ctx, 400)
		return
	}

	if incomeUpdate.Name != nil {
		income.Name = incomeUpdate.Name
	}
	if incomeUpdate.Amount != nil {
		income.Amount = *incomeUpdate.Amount
	}
	if incomeUpdate.Comment != nil {
		income.Comment = incomeUpdate.Comment
	}
	if incomeUpdate.PayerID != nil {
		income.PayerID = *incomeUpdate.PayerID
	}
	if incomeUpdate.CreatedAt != nil {
		income.CreatedAt = *incomeUpdate.CreatedAt
	}

	tx := DB.Begin()
	result = tx.Save(&income)
	if result.Error != nil {
		SendMessageOnly("Could not update income: "+result.Error.Error(), ctx, 500)
		tx.Rollback()
		return
	}

	tx.Commit()
	SendMessageOnly("Income was updated successfully", ctx, 200)
}

// ========== DELETE /incomes/:id ==========

func DeleteIncome(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	var income models.Income
	id := ctx.Param("id")

	result := DB.First(&income, "id = ?", id)
	if result.Error != nil {
		SendMessageOnly("Could not get existing income: "+result.Error.Error(), ctx, 500)
		return
	}

	if income.UserID.String() != userID {
		SendMessageOnly("You are not allowed to delete this income", ctx, 403)
		return
	}

	income.IsActive = false

	tx := DB.Begin()
	result = tx.Save(&income)
	if result.Error != nil {
		SendMessageOnly("Could not delete income: "+result.Error.Error(), ctx, 500)
		tx.Rollback()
		return
	}

	tx.Commit()
	SendMessageOnly("Income was deleted successfully", ctx, 200)
}
