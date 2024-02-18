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

func CheckPassInUseValidityWrapper(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, false)
	if userID == "" {
		return
	}

	tx := DB.Begin()
	id := ctx.Param("id")
	valid, err := CheckPassInUseValidity(tx, id)
	if err != nil {
		tx.Rollback()
		SendMessageOnly(err.Error(), ctx, 500)
		return
	}

	tx.Commit()
	ctx.JSON(200, valid)
}
func CheckPassInUseValidity(tx *gorm.DB, id string) (bool, error) {
	var passInUse models.PassInUse
	result := DB.Preload("Pass").First(&passInUse, "id = ?", id)
	if result.Error != nil {
		return false, result.Error
	}

	now := time.Now()
	occasionLimit := passInUse.Pass.OccasionLimit

	fmt.Println(passInUse)
	fmt.Println("Occasion limit: ", occasionLimit)
	fmt.Println("Occasions: ", passInUse.Occasions)
	fmt.Println("passInUse.Occasions >= *occasionLimit: ", passInUse.Occasions >= *occasionLimit)
	fmt.Println("now.After(*passInUse.ValidUntil): ", now.After(*passInUse.ValidUntil), passInUse.ValidUntil)
	fmt.Println("now.Before(*passInUse.ValidFrom): ", now.Before(*passInUse.ValidFrom), passInUse.ValidFrom)

	if now.Before(*passInUse.ValidFrom) {
		return false, nil
	}

	if occasionLimit == nil {
		if now.After(*passInUse.ValidUntil) {
			passInUse.IsActive = false
			result := tx.Save(&passInUse)
			if result.Error != nil {
				return false, result.Error
			}
			return false, nil
		}
	} else {
		if now.After(*passInUse.ValidUntil) || passInUse.Occasions >= *occasionLimit {
			passInUse.IsActive = false
			result := tx.Save(&passInUse)
			if result.Error != nil {
				return false, result.Error
			}
			return false, nil
		}
	}
	return true, nil
}

func GetPassesInUse(ctx *gin.Context) {
	userID, isAdmin := CheckAuth(ctx, false)
	if userID == "" {
		return
	}

	var passesInUse []models.PassInUse
	var result *gorm.DB
	if isAdmin {
		result = DB.Where("user_id = ?", userID).Preload("Pass").Find(&passesInUse)
	} else {
		result = DB.Where("payer_id = ?", userID).Preload("Pass").Find(&passesInUse)
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

	// TODO: Also create income

	tx.Commit()
	SendMessageOnly("Pass in use was created successfully", ctx, 201)
}

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

func UsePassInUseWrapper(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	id := ctx.Param("id")

	tx := DB.Begin()
	response := UsePassInUse(tx, id)
	if !response.Success {
		tx.Rollback()
		SendMessageOnly(response.Message, ctx, 500)
		return
	}

	tx.Commit()
	SendMessageOnly(response.Message, ctx, 200)
}

// UsePassInUse increases the occasions attribute by 1 + checks for validity
func UsePassInUse(tx *gorm.DB, id string) ActionResponse {
	var passInUse models.PassInUse
	result := tx.First(&passInUse, "id = ?", id)
	if result.Error != nil {
		return ActionResponse{
			Success: false,
			Message: result.Error.Error(),
		}
	}

	valid, err := CheckPassInUseValidity(tx, id)

	if err != nil {
		return ActionResponse{
			Success: false,
			Message: err.Error(),
		}
	}

	if !valid {
		return ActionResponse{
			Success: true,
			Message: "Pass in use is not valid",
		}
	} else {
		result = tx.Model(&passInUse).Update("occasions", passInUse.Occasions+1)
		if result.Error != nil {
			return ActionResponse{
				Success: false,
				Message: result.Error.Error(),
			}
		}

		return ActionResponse{
			Success: true,
			Message: "Pass in use was used successfully",
		}
	}
}

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
