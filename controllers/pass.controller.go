package controllers

import (
	"elrek-system_GO/models"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
	"gorm.io/gorm"
)

// MARK: GET
func GetPasses(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	var passes []models.Pass
	result := DB.Where("user_id = ? and is_active = ?", userID, true).Preload("Services").Order("created_at desc").Find(&passes)
	if result.Error != nil {
		SendMessageOnly("Could not get passes: "+result.Error.Error(), ctx, 500)
		return
	}

	ctx.JSON(200, passes)
}

func GetPassWrapper(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, false)
	if userID == "" {
		return
	}

	id := ctx.Param("id")
	pass, err := getPass(id)
	if err != nil {
		SendMessageOnly("Could not get pass: "+err.Error(), ctx, 500)
		return
	}

	if pass.UserID.String() != userID {
		SendMessageOnly("You are not allowed to get this pass", ctx, 403)
		return
	}

	ctx.JSON(200, pass)
}
func getPass(passID string) (models.Pass, error) {
	var pass models.Pass

	result := DB.Preload("Services").First(&pass, "id = ?", passID)
	if result.Error != nil {
		return models.Pass{}, result.Error
	}
	return pass, nil
}

// MARK: CREATE
func CreatePassWrapper(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	var passCreate models.PassCreate
	err := ctx.BindJSON(&passCreate)
	if err != nil {
		SendMessageOnly("Parse error: "+err.Error(), ctx, 400)
		return
	}

	slog.Info("CreatePassWrapper", "passCreate: ", passCreate)

	tx := DB.Begin()
	result := createPass(passCreate, userID, tx)
	if !result.Success {
		tx.Rollback()
		SendMessageOnly(result.Message, ctx, 500)
		return
	}
	tx.Commit()
	SendMessageOnly("Pass was created successfully", ctx, 201)
}

// MARK: createPass
func createPass(passCreate models.PassCreate, userID string, tx *gorm.DB) ActionResponse {
	var pass models.Pass

	pass.UserID = openapitypes.UUID(uuid.MustParse(userID))
	pass.Name = passCreate.Name
	pass.Price = passCreate.Price
	pass.IsActive = true
	pass.Comment = passCreate.Comment
	pass.PrevPassID = openapitypes.UUID{}
	pass.OccasionLimit = passCreate.OccasionLimit
	pass.Duration = passCreate.Duration
	pass.ID = openapitypes.UUID(uuid.New())

	slog.Info("createPass", "pass: ", pass)

	var services []models.Service
	result := DB.Find(&services, passCreate.ServiceIDs)
	if result.Error != nil {
		return ActionResponse{Success: false, Message: "Could not get services for pass creation: " + result.Error.Error()}
	}

	result = tx.Save(&pass)
	if result.Error != nil {
		return ActionResponse{Success: false, Message: "Could not save pass: " + result.Error.Error()}
	}

	// https://github.com/harranali/gorm-relationships-examples/tree/main/many-to-many
	// next assign the languages to the user
	// db.Model(&user).Association("Languages").Append(&languages)
	err := tx.Model(&pass).Association("Services").Append(&services)
	if err != nil {
		return ActionResponse{Success: false, Message: "Could not link pass and services: " + result.Error.Error()}
	}
	return ActionResponse{Success: true, Message: "SUCCESS"}
}

// MARK: UPDATE
func UpdatePass(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, false)
	if userID == "" {
		return
	}

	var passUpdate models.PassUpdate
	err := ctx.BindJSON(&passUpdate)
	if err != nil {
		SendMessageOnly("Parse error: "+err.Error(), ctx, 400)
		return
	}

	var pass models.Pass
	id := ctx.Param("id")

	slog.Info("UpdatePass", "passUpdate: ", pass, "id: ", id)

	result := DB.Preload("Services").First(&pass, "id = ?", id)
	if result.Error != nil {
		SendMessageOnly("Could not get existing pass: "+result.Error.Error(), ctx, 500)
		return
	}

	if pass.UserID.String() != userID {
		SendMessageOnly("You are not allowed to update this pass", ctx, 403)
		return
	}

	update := false
	if passUpdate.Name != nil {
		update = true
		pass.Name = *passUpdate.Name
	}
	if passUpdate.Price != nil {
		update = false
		pass.Price = *passUpdate.Price
	}
	if passUpdate.Comment != nil {
		update = true
		pass.Comment = passUpdate.Comment
	}
	if passUpdate.Duration != nil {
		update = false
		pass.Duration = passUpdate.Duration
	}
	if passUpdate.OccasionLimit != nil {
		update = false
		pass.OccasionLimit = passUpdate.OccasionLimit
	}
	if passUpdate.ServiceIDs != nil {
		var updatedServiceIDs []string
		for _, serviceID := range *passUpdate.ServiceIDs {
			updatedServiceIDs = append(updatedServiceIDs, serviceID.String())
		}
		update = false
		result := DB.Find(&pass.Services, updatedServiceIDs)
		if result.Error != nil {
			SendMessageOnly("Could not get service: "+result.Error.Error(), ctx, 500)
			return
		}
	}

	tx := DB.Begin()

	if update {
		result = tx.Save(&pass)
		if result.Error != nil {
			tx.Rollback()
			SendMessageOnly("Could not update pass: "+result.Error.Error(), ctx, 500)
			return
		}
	} else {
		deleteResult := deletePass(pass.ID, tx)
		if !deleteResult.Success {
			SendMessageOnly(deleteResult.Message, ctx, 500)
			tx.Rollback()
			return
		}

		pass.PrevPassID = pass.ID
		pass.ID = openapitypes.UUID(uuid.New())
		result = tx.Create(&pass)

		if result.Error != nil {
			tx.Rollback()
			SendMessageOnly("Could not create new pass: "+result.Error.Error(), ctx, 500)
			return
		}

		err := tx.Model(&pass).Association("Services").Append(&pass.Services)
		if err != nil {
			tx.Rollback()
			SendMessageOnly("Could not link pass and services: "+result.Error.Error(), ctx, 500)
		}
	}

	tx.Commit()
	SendMessageOnly("Pass was updated successfully", ctx, 200)
}

// MARK: DELETE
func DeletePassWrapper(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	tx := DB.Begin()
	passID := openapitypes.UUID(uuid.MustParse(ctx.Param("id")))
	deleteResult := deletePass(passID, tx)
	if !deleteResult.Success {
		tx.Rollback()
		SendMessageOnly(deleteResult.Message, ctx, 500)
		return
	}

	tx.Commit()
	SendMessageOnly("Pass was deleted successfully", ctx, 200)
}

func deletePass(passID openapitypes.UUID, tx *gorm.DB) ActionResponse {
	var pass models.Pass

	result := tx.First(&pass, "id = ?", passID)
	if result.Error != nil {
		return ActionResponse{Success: false, Message: "Could not get existing pass: " + result.Error.Error()}
	}

	pass.IsActive = false

	result = tx.Save(&pass)
	if result.Error != nil {
		return ActionResponse{Success: false, Message: "Could not delete pass: " + result.Error.Error()}
	}
	return ActionResponse{Success: true, Message: "SUCCESS"}
}
