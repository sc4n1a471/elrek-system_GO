package controllers

import (
	"elrek-system_GO/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
	"gorm.io/gorm"
	"time"
)

// =========== Checking validity and using pass in use ===========

const (
	NoValidPIUWasFound     string = "No valid pass in use was found"
	NoPIUWasFound          string = "No pass in use found for payer"
	PIUWasUsedSuccessfully        = "Pass in use was used successfully"
)

func CheckPassInUseValidityWrapper(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, false)
	if userID == "" {
		return
	}

	tx := DB.Begin()
	id := ctx.Param("id")
	valid, err := checkPassInUseValidity(tx, openapitypes.UUID(uuid.MustParse(id)), openapitypes.UUID{})
	if err != nil {
		tx.Rollback()
		SendMessageOnly(err.Error(), ctx, 500)
		return
	}

	tx.Commit()
	ctx.JSON(200, valid)
}
func checkPassInUseValidity(tx *gorm.DB, passInUseID openapitypes.UUID, serviceID openapitypes.UUID) (bool, error) {
	fmt.Println("============ checkPassInUseValidity begin ============")
	fmt.Println("PassInUseID: ", passInUseID)
	var passInUse models.PassInUse
	result := DB.Preload("Pass").First(&passInUse, "id = ?", passInUseID)
	if result.Error != nil {
		return false, result.Error
	}

	now := time.Now()
	occasionLimit := passInUse.Pass.OccasionLimit

	if now.Before(*passInUse.ValidFrom) {
		return false, nil
	}

	if occasionLimit == nil {
		fmt.Println("occasionLimit is nil")
		if now.After(*passInUse.ValidUntil) {
			passInUse.IsActive = false
			result := tx.Save(&passInUse)
			if result.Error != nil {
				return false, result.Error
			}
			return false, nil
		}
	} else {
		fmt.Println("occasionLimit is not nil")
		fmt.Println("now.After(*passInUse.ValidUntil)", *passInUse.ValidUntil, ": ", now.After(*passInUse.ValidUntil))
		fmt.Println("passInUse.Occasions (", passInUse.Occasions, ") >= *occasionLimit (", *occasionLimit, "):", passInUse.Occasions >= *occasionLimit)
		if now.After(*passInUse.ValidUntil) || passInUse.Occasions >= *occasionLimit {
			fmt.Println("Deactivating passInUse...")
			passInUse.IsActive = false
			result := tx.Save(&passInUse)
			if result.Error != nil {
				return false, result.Error
			}
			return false, nil
		}
	}

	defaultUUID := openapitypes.UUID{}
	if serviceID != defaultUUID {
		fmt.Println("ServiceID is not default, searching for: ", serviceID)
		pass, err := getPass(passInUse.PassID.String())
		fmt.Println("Pass found: ", pass)
		if err != nil {
			return false, err
		}

		for _, service := range pass.Services {

			fmt.Println("Checking: ", service.ID, " -> ", service.ID == serviceID)
			if service.ID == serviceID {
				fmt.Println("Service found in pass -> user has valid pass")
				return true, nil
			}
		}

		fmt.Println("Searched serviceID was not found in pass' valid services array")
		return false, nil
	}
	return true, nil
}

func CheckPayerHasValidPIUForServiceWrapper(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, false)
	if userID == "" {
		return
	}

	tx := DB.Begin()
	payerID := ctx.Param("payer_id")
	serviceID := ctx.Param("service_id")
	valid, err := checkPayerHasValidPIUForService(
		tx,
		openapitypes.UUID(uuid.MustParse(payerID)),
		openapitypes.UUID(uuid.MustParse(serviceID)),
	)

	if err != nil {
		tx.Rollback()
		SendMessageOnly(err.Error(), ctx, 500)
		return
	}

	tx.Commit()
	ctx.JSON(200, valid)

}
func checkPayerHasValidPIUForService(tx *gorm.DB, payerID openapitypes.UUID, serviceID openapitypes.UUID) (bool, error) {
	var passInUse models.PassInUse
	result := DB.First(&passInUse, "payer_id = ? and is_active = ?", payerID, true)
	if result.Error != nil {
		return false, result.Error
	}

	valid, err := checkPassInUseValidity(tx, passInUse.ID, serviceID)
	if err != nil {
		return false, err
	}

	return valid, nil
}

func UsePassInUseWrapper(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	id := ctx.Param("id")

	tx := DB.Begin()
	response := usePassInUse(tx, openapitypes.UUID(uuid.MustParse(id)), openapitypes.UUID{})
	if !response.Success {
		tx.Rollback()
		SendMessageOnly(response.Message, ctx, 500)
		return
	}

	tx.Commit()
	SendMessageOnly(response.Message, ctx, 200)
}

func usePassInUse(tx *gorm.DB, payerID openapitypes.UUID, serviceID openapitypes.UUID) ActionResponse {
	var passesInUse []models.PassInUse
	result := tx.Find(&passesInUse, "payer_id = ? and is_active = ?", payerID, true)
	if result.Error != nil {
		return ActionResponse{
			Success: false,
			Message: NoPIUWasFound,
		}
	}

	for _, passInUse := range passesInUse {
		valid, err := checkPassInUseValidity(tx, passInUse.ID, serviceID)

		if err != nil {
			return ActionResponse{
				Success: false,
				Message: err.Error(),
			}
		}

		if valid {
			result = tx.Model(&passInUse).Update("occasions", passInUse.Occasions+1)
			if result.Error != nil {
				return ActionResponse{
					Success: false,
					Message: result.Error.Error(),
				}
			}

			return ActionResponse{
				Success: true,
				Message: PIUWasUsedSuccessfully,
			}
		}
	}

	return ActionResponse{
		Success: false,
		Message: NoValidPIUWasFound,
	}
}

// =========== GET /passes_in_use ===========

func GetPassesInUse(ctx *gin.Context) {
	userID, isAdmin := CheckAuth(ctx, false)
	if userID == "" {
		return
	}

	var passesInUse []models.PassInUse
	var result *gorm.DB
	if isAdmin {
		result = DB.Where("user_id = ? and is_active = ?", userID, true).Preload("Pass").Find(&passesInUse)
	} else {
		result = DB.Where("payer_id = ? and is_active = ?", userID, true).Preload("Pass").Find(&passesInUse)
	}

	if result.Error != nil {
		SendMessageOnly("Could not get passes in use: "+result.Error.Error(), ctx, 500)
		return
	}

	ctx.JSON(200, passesInUse)
}

func GetPassInUse(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, false)
	if userID == "" {
		return
	}

	var passInUse models.PassInUse
	id := ctx.Param("id")

	result := DB.Preload("Pass").First(&passInUse, "id = ?", id)
	if result.Error != nil {
		SendMessageOnly("Could not get pass in use: "+result.Error.Error(), ctx, 500)
		return
	}

	ctx.JSON(200, passInUse)
}

func getPassInUse(id openapitypes.UUID) (models.PassInUse, error) {
	var passInUse models.PassInUse

	result := DB.Preload("Pass").First(&passInUse, "id = ?", id)
	if result.Error != nil {
		return models.PassInUse{}, result.Error
	}
	return passInUse, nil
}

// =========== POST /passes_in_use ===========

func CreatePassInUse(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	var passInUseCreate models.PassInUseCreate
	err := ctx.BindJSON(&passInUseCreate)
	if err != nil {
		SendMessageOnly("Parse error: "+err.Error(), ctx, 400)
		return
	}

	tx := DB.Begin()

	var pass models.Pass
	result := DB.First(&pass, "id = ?", passInUseCreate.PassID)
	if result.Error != nil {
		tx.Rollback()
		SendMessageOnly("Could not get pass: "+result.Error.Error(), ctx, 500)
		return
	}

	var payer models.User
	result = DB.First(&payer, "id = ?", passInUseCreate.PayerID)
	if result.Error != nil {
		tx.Rollback()
		SendMessageOnly("Could not get user (payer): "+result.Error.Error(), ctx, 500)
		return
	}

	roundedValidFrom := passInUseCreate.ValidFrom.Round(time.Second)
	roundedValidUntil := passInUseCreate.ValidUntil.Round(time.Second)

	passInUse := models.PassInUse{
		IsActive:   true,
		ID:         openapitypes.UUID(uuid.New()),
		Occasions:  0,
		Comment:    passInUseCreate.Comment,
		UserID:     openapitypes.UUID(uuid.MustParse(userID)),
		PassID:     passInUseCreate.PassID,
		PayerID:    passInUseCreate.PayerID,
		ValidFrom:  &roundedValidFrom,
		ValidUntil: &roundedValidUntil,
	}
	result = tx.Create(&passInUse)
	if result.Error != nil {
		tx.Rollback()
		SendMessageOnly("Could not create pass in use: "+result.Error.Error(), ctx, 500)
		return
	}

	err = tx.Model(&passInUse).Association("Pass").Append(&pass)
	if err != nil {
		tx.Rollback()
		SendMessageOnly("Could not associate pass with pass in use: "+err.Error(), ctx, 500)
		return
	}

	income := models.IncomeCreate{
		Amount:      passInUse.Pass.Price,
		PassInUseID: &passInUse.ID,
		PayerID:     payer.ID,
	}
	incomeResult := createIncome(tx, income, openapitypes.UUID(uuid.MustParse(userID)), 0)
	if !incomeResult.Success {
		tx.Rollback()
		SendMessageOnly(incomeResult.Message, ctx, 500)
		return
	}

	tx.Commit()
	SendMessageOnly("Pass in use was created successfully", ctx, 201)
}

// =========== PATCH /passes_in_use/:id ===========

func UpdatePassInUse(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	var passInUseUpdate models.PassInUseUpdate
	err := ctx.BindJSON(&passInUseUpdate)
	if err != nil {
		SendMessageOnly("Parse error: "+err.Error(), ctx, 400)
		return
	}

	var passInUse models.PassInUse
	id := ctx.Param("id")
	result := DB.First(&passInUse, "id = ?", id)
	if result.Error != nil {
		SendMessageOnly("Could not get pass in use: "+result.Error.Error(), ctx, 500)
		return
	}

	if passInUse.UserID.String() != userID {
		SendMessageOnly("You are not allowed to update this pass in use", ctx, 403)
		return
	}

	if passInUseUpdate.Occasions != nil {
		passInUse.Occasions = *passInUseUpdate.Occasions
	}
	if passInUseUpdate.Comment != nil {
		passInUse.Comment = passInUseUpdate.Comment
	}
	if passInUseUpdate.ValidFrom != nil {
		roundedTime := passInUseUpdate.ValidFrom.Round(time.Second)
		passInUse.ValidFrom = &roundedTime
	}
	if passInUseUpdate.ValidUntil != nil {
		roundedTime := passInUseUpdate.ValidUntil.Round(time.Second)
		passInUse.ValidUntil = &roundedTime
	}

	result = DB.Save(&passInUse)
	if result.Error != nil {
		SendMessageOnly("Could not update pass in use: "+result.Error.Error(), ctx, 500)
		return
	}

	SendMessageOnly("Pass in use was updated successfully", ctx, 200)
}

// =========== DELETE /passes_in_use/:id ===========

func DeletePassInUse(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	var passInUse models.PassInUse
	id := ctx.Param("id")
	result := DB.First(&passInUse, "id = ?", id)
	if result.Error != nil {
		SendMessageOnly("Could not get pass in use: "+result.Error.Error(), ctx, 500)
		return
	}

	if passInUse.UserID.String() != userID {
		SendMessageOnly("You are not allowed to delete this pass in use", ctx, 403)
		return
	}

	tx := DB.Begin()
	passInUse.IsActive = false
	result = tx.Save(&passInUse)
	if result.Error != nil {
		tx.Rollback()
		SendMessageOnly("Could not delete pass in use: "+result.Error.Error(), ctx, 500)
		return
	}

	tx.Commit()
	ctx.JSON(200, nil)
}
