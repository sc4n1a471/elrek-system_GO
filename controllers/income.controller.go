package controllers

import (
	"elrek-system_GO/models"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
	"gorm.io/gorm"
)

// MARK: GET /incomes
func GetIncomes(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	var incomes []models.Income
	result := DB.Where("user_id = ? and is_active = ?", userID, true).
		Preload("ActivePass.Pass").
		Preload("User").
		Preload("Service").
		Order("created_at desc").
		Find(&incomes)
	if result.Error != nil {
		SendMessageOnly("Could not get incomes: "+result.Error.Error(), ctx, 500)
		return
	}

	ctx.JSON(200, incomes)
}

func GetIncome(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	var income models.Income
	id := ctx.Param("id")

	result := DB.Preload("ActivePass").
		Preload("Service").
		First(&income, "id = ?", id)
	if result.Error != nil {
		SendMessageOnly("Could not get income: "+result.Error.Error(), ctx, 500)
		return
	}

	ctx.JSON(200, income)
}

// MARK: POST /incomes
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

	slog.Info("CreateIncomeWrapper", "incomeCreate: ", incomeCreate)

	result := createIncome(tx, incomeCreate, openapitypes.UUID(uuid.MustParse(userID)), 1)
	if !result.Success {
		SendMessageOnly(result.Message, ctx, 500)
		tx.Rollback()
		return
	}

	tx.Commit()
	SendMessageOnly("Income was created successfully", ctx, 201)
}

// MARK: Create multiple incomes
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

	slog.Info("CreateIncomeMultipleUsersWrapper", "incomeCreateMultipleUsers: ", incomeCreateMultipleUsers)

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
		} else if incomeCreateMultipleUsers.ActivePassIDs != nil {
			for _, ActivePassID := range *incomeCreateMultipleUsers.ActivePassIDs {
				incomeCreate.ActivePassID = &ActivePassID

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
			SendMessageOnly("ServiceIDs or ActivePassIDs must be provided", ctx, 400)
			tx.Rollback()
			return
		}
	}

	tx.Commit()
	SendMessageOnly("Incomes were created successfully", ctx, 201)
}

// MARK: createIncome
// createIncome creates an income using the provided data (incomeCreate) and the provided user ID (userID
func createIncome(tx *gorm.DB, incomeCreate models.IncomeCreate, userID openapitypes.UUID, numOfAttendees int) ActionResponse {

	// ========== Creating income model ==========
	income := models.Income{
		IsActive:     true,
		Comment:      incomeCreate.Comment,
		ID:           openapitypes.UUID(uuid.New()),
		UserID:       userID,
		ActivePassID: incomeCreate.ActivePassID,
		ServiceID:    incomeCreate.ServiceID,
		PayerID:      incomeCreate.PayerID,
		Name:         incomeCreate.Name,
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

		if incomeCreate.ActivePassID != nil && incomeCreate.Name == nil {
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
			if incomeCreate.Name == nil {
				income.Name = &service.Name
			}

			// MARK: Using active pass
			useResult := useActivePass(tx, incomeCreate.PayerID, service.ID)
			slog.Info("createIncome", "useResult", useResult)

			if !useResult.Success {
				if useResult.Message == NoActivePassWasFound || useResult.Message == NoValidActivePassWasFound {
					income.Amount = service.Price

					var dynamicPrices []models.DynamicPrice
					dynamicPrices, dpResult := getDynamicPrices(service.ID)

					slog.Info("createIncome", "dynamicPrices", dynamicPrices)
					slog.Info("createIncome", "dpResult", dpResult)

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
								slog.Info("Dynamic price used: ", "income.Amount", income.Amount)
								break
							}
						}
					}
				} else {
					return ActionResponse{
						Success: false,
						Message: "Could not use active pass: " + useResult.Message,
					}
				}
			} else {
				// Has valid activePass
				income.Amount = 0
				income.IsPaid = true
			}
		} else {
			return ActionResponse{
				Success: false,
				Message: "Amount must be provided if active pass or service is not provided",
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

	if income.ActivePassID != nil {
		var activePass models.ActivePass
		result = tx.First(&activePass, "id = ?", income.ActivePassID)

		err := tx.Model(&income).Association("ActivePass").Append(&activePass)
		if err != nil {
			return ActionResponse{
				Success: false,
				Message: "Could not associate active pass with income: " + err.Error(),
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

// MARK: PATCH /incomes/:id
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
	if incomeUpdate.IsPaid != nil {
		income.IsPaid = *incomeUpdate.IsPaid
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

// MARK: DELETE /incomes/:id

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
