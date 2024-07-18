package controllers

import (
	"elrek-system_GO/models"
	"fmt"
	"log/slog"
	"time"

	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
	"gorm.io/gorm"
)

// =========== Checking validity and using active pass ===========

const (
	NoValidActivePassWasFound     string = "No valid active pass was found"
	NoActivePassWasFound          string = "No active pass found for payer"
	ActivePassWasUsedSuccessfully        = "Pass in use was used successfully"
)

// MARK: Check validity
func CheckactivePassValidityWrapper(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, false)
	if userID == "" {
		return
	}

	tx := DB.Begin()
	id := ctx.Param("id")

	slog.Info("CheckactivePassValidityWrapper", "id", id)

	valid, err := checkActivePassValidity(tx, openapitypes.UUID(uuid.MustParse(id)), openapitypes.UUID{})
	if err != nil {
		tx.Rollback()
		SendMessageOnly(err.Error(), ctx, 500)
		return
	}

	tx.Commit()
	ctx.JSON(200, valid)
}

func checkActivePassValidity(tx *gorm.DB, ActivePassID openapitypes.UUID, serviceID openapitypes.UUID) (bool, error) {
	slog.Info("============ checkactivePassValidity begin ============", "ActivePassID: ", ActivePassID, "serviceID: ", serviceID)
	var activePass models.ActivePass
	result := DB.Preload("Pass").First(&activePass, "id = ?", ActivePassID)
	if result.Error != nil {
		return false, result.Error
	}

	now := time.Now()
	occasionLimit := activePass.Pass.OccasionLimit

	if now.Before(*activePass.ValidFrom) {
		return false, nil
	}

	if occasionLimit == nil {
		slog.Info("occasionLimit is nil")

		if activePass.ValidUntil == nil {
			slog.Info("activePass.ValidUntil is nil, too...")
			return false, errors.New("The Active Pass does not have a valid occasion limit or a valid until date")
		}

		if now.After(*activePass.ValidUntil) {
			activePass.IsActive = false
			result := tx.Save(&activePass)
			if result.Error != nil {
				return false, result.Error
			}
			return false, nil
		}
	} else {
		slog.Info("occasionLimit is not nil")
		if activePass.ValidUntil == nil {
			slog.Info("activePass.ValidUntil is nil")
			if activePass.Occasions >= *occasionLimit {
				slog.Info("Deactivating activePass...")
				activePass.IsActive = false
				result := tx.Save(&activePass)
				if result.Error != nil {
					return false, result.Error
				}
				return false, nil
			}
		} else {
			slog.Info(
				"activePass.ValidUntil is not nil",
				slog.String("*activePass.ValidUntil.String()", activePass.ValidUntil.String()),
				slog.String("now.After(*activePass.ValidUntil)", fmt.Sprintf("%v", now.After(*activePass.ValidUntil))),
				slog.String("activePass.Occasions", fmt.Sprintf("%v", activePass.Occasions)),
				slog.String("*occasionLimit", fmt.Sprintf("%v", *occasionLimit)),
				slog.String("activePass.Occasions >= *occasionLimit", fmt.Sprintf("%v", activePass.Occasions >= *occasionLimit)),
			)

			if now.After(*activePass.ValidUntil) || activePass.Occasions >= *occasionLimit {
				slog.Info("Deactivating activePass...")
				activePass.IsActive = false
				result := tx.Save(&activePass)
				if result.Error != nil {
					return false, result.Error
				}
				return false, nil
			}
		}
	}

	// MARK: Checking if pass is valid for a specific service
	defaultUUID := openapitypes.UUID{}
	if serviceID != defaultUUID {
		slog.Info("ServiceID is not default ", "searching for:", serviceID)
		pass, err := getPass(activePass.PassID.String())
		slog.Info("Pass found: ", "pass", pass)

		if err != nil {
			return false, err
		}

		prevServices, err := GetPrevServices(serviceID)
		if err != nil {
			return false, err
		}

		for _, service := range pass.Services {

			slog.Info("Checking: ",
				slog.String("service.ID (checking)", fmt.Sprintf("%v", service.ID)),
				slog.String("serviceID (search for)", fmt.Sprintf("%v", serviceID)),
				slog.String("service.ID == serviceID", fmt.Sprintf("%v", service.ID == serviceID)),
			)
			if service.ID == serviceID {
				slog.Info("Service found in pass -> user has valid pass")
				return true, nil
			} else {

				// Check if the user has an older pass that is valid for the used service's predecessor
				for _, prevService := range prevServices {
					slog.Info("Checking prev service: ",
						slog.String("prevService.ID (checking)", fmt.Sprintf("%v", prevService.ID)),
						slog.String("service.ID (search for)", fmt.Sprintf("%v", service.ID)),
						slog.String("prevService.ID == service.ID", fmt.Sprintf("%v", prevService.ID == service.ID)),
					)

					if prevService.ID == service.ID {
						slog.Info("Service found in prevServices -> user has valid pass for previous service")
						return true, nil
					}
				}
			}
		}

		slog.Info("Searched serviceID was not found in pass' valid services array")
		return false, nil
	}
	return true, nil
}

func CheckPayerHasValidActivePassForServiceWrapper(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, false)
	if userID == "" {
		return
	}

	tx := DB.Begin()
	payerID := ctx.Param("payer_id")
	serviceID := ctx.Param("service_id")

	slog.Info("CheckPayerHasValidActivePassForServiceWrapper", "payerID", payerID, "serviceID", serviceID)

	valid, err := checkPayerHasValidActivePassForService(
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
func checkPayerHasValidActivePassForService(tx *gorm.DB, payerID openapitypes.UUID, serviceID openapitypes.UUID) (bool, error) {
	var activePass models.ActivePass
	result := DB.First(&activePass, "payer_id = ? and is_active = ?", payerID, true)
	if result.Error != nil {
		return false, result.Error
	}

	valid, err := checkActivePassValidity(tx, activePass.ID, serviceID)
	if err != nil {
		return false, err
	}

	return valid, nil
}

// MARK: USE ACTIVE PASS
func UseActivePassWrapper(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	id := ctx.Param("id")
	slog.Info("UseActivePassWrapper", "id", id)

	tx := DB.Begin()
	response := useActivePass(tx, openapitypes.UUID(uuid.MustParse(id)), openapitypes.UUID{})
	if !response.Success {
		tx.Rollback()
		SendMessageOnly(response.Message, ctx, 500)
		return
	}

	tx.Commit()
	SendMessageOnly(response.Message, ctx, 200)
}

// serviceID: checks if the pass is valid for a specific service
func useActivePass(tx *gorm.DB, payerID openapitypes.UUID, serviceID openapitypes.UUID) ActionResponse {
	var activePasses []models.ActivePass
	result := tx.Find(&activePasses, "payer_id = ? and is_active = ?", payerID, true)
	if result.Error != nil {
		return ActionResponse{
			Success: false,
			Message: NoActivePassWasFound,
		}
	}

	for _, activePass := range activePasses {
		valid, err := checkActivePassValidity(tx, activePass.ID, serviceID)

		if err != nil {
			return ActionResponse{
				Success: false,
				Message: err.Error(),
			}
		}

		if valid {
			result = tx.Model(&activePass).Update("occasions", activePass.Occasions+1)
			if result.Error != nil {
				return ActionResponse{
					Success: false,
					Message: result.Error.Error(),
				}
			}

			return ActionResponse{
				Success: true,
				Message: ActivePassWasUsedSuccessfully,
			}
		}
	}

	return ActionResponse{
		Success: false,
		Message: NoValidActivePassWasFound,
	}
}

// MARK: GET /active-passes
func GetActivePasses(ctx *gin.Context) {
	userID, isAdmin := CheckAuth(ctx, false)
	if userID == "" {
		return
	}

	var activePasses []models.ActivePass
	var result *gorm.DB
	if isAdmin {
		result = DB.Where("user_id = ? and is_active = ?", userID, true).Preload("Pass").Preload("Pass.Services").Preload("User").Order("created_at desc").Find(&activePasses)
	} else {
		result = DB.Where("payer_id = ? and is_active = ?", userID, true).Preload("Pass").Preload("Pass.Services").Preload("User").Order("created_at desc").Find(&activePasses)
	}

	if result.Error != nil {
		SendMessageOnly("Could not get passes in use: "+result.Error.Error(), ctx, 500)
		return
	}

	ctx.JSON(200, activePasses)
}

func GetActivePass(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, false)
	if userID == "" {
		return
	}

	var activePass models.ActivePass
	id := ctx.Param("id")

	result := DB.Preload("Pass").Preload("User").First(&activePass, "id = ?", id)
	if result.Error != nil {
		SendMessageOnly("Could not get active pass: "+result.Error.Error(), ctx, 500)
		return
	}

	ctx.JSON(200, activePass)
}

func getActivePass(id openapitypes.UUID) (models.ActivePass, error) {
	var activePass models.ActivePass

	result := DB.Preload("Pass").First(&activePass, "id = ?", id)
	if result.Error != nil {
		return models.ActivePass{}, result.Error
	}
	return activePass, nil
}

// MARK: POST /active-passes
func CreateActivePass(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	var activePassCreate models.ActivePassCreate
	err := ctx.BindJSON(&activePassCreate)
	if err != nil {
		SendMessageOnly("Parse error: "+err.Error(), ctx, 400)
		return
	}

	slog.Info("CreateActivePass", "activePassCreate", activePassCreate)

	tx := DB.Begin()

	var pass models.Pass
	result := DB.First(&pass, "id = ?", activePassCreate.PassID)
	if result.Error != nil {
		tx.Rollback()
		SendMessageOnly("Could not get pass: "+result.Error.Error(), ctx, 500)
		return
	}

	var payer models.User
	result = DB.First(&payer, "id = ?", activePassCreate.PayerID)
	if result.Error != nil {
		tx.Rollback()
		SendMessageOnly("Could not get user (payer): "+result.Error.Error(), ctx, 500)
		return
	}

	slog.Info("CreateActivePass", "activePassCreate.ValidFrom", activePassCreate.ValidFrom)
	roundedValidFrom := activePassCreate.ValidFrom.Round(time.Second)
	activePass := models.ActivePass{}

	slog.Info("CreateActivePass",
		slog.String("activePassCreate.ValidUntil", fmt.Sprintf("%v", activePassCreate.ValidUntil)),
		slog.String("activePassCreate.ValidUntil == nil", fmt.Sprintf("%v", activePassCreate.ValidUntil == nil)),
	)

	if activePassCreate.ValidUntil == nil {
		activePass = models.ActivePass{
			IsActive:  true,
			ID:        openapitypes.UUID(uuid.New()),
			Occasions: 0,
			Comment:   activePassCreate.Comment,
			UserID:    openapitypes.UUID(uuid.MustParse(userID)),
			PassID:    activePassCreate.PassID,
			PayerID:   activePassCreate.PayerID,
			ValidFrom: &roundedValidFrom,
		}
	} else {
		roundedValidUntil := activePassCreate.ValidUntil.Round(time.Second)
		activePass = models.ActivePass{
			IsActive:   true,
			ID:         openapitypes.UUID(uuid.New()),
			Occasions:  0,
			Comment:    activePassCreate.Comment,
			UserID:     openapitypes.UUID(uuid.MustParse(userID)),
			PassID:     activePassCreate.PassID,
			PayerID:    activePassCreate.PayerID,
			ValidFrom:  &roundedValidFrom,
			ValidUntil: &roundedValidUntil,
		}
	}

	result = tx.Create(&activePass)
	if result.Error != nil {
		tx.Rollback()
		SendMessageOnly("Could not create active pass: "+result.Error.Error(), ctx, 500)
		return
	}

	err = tx.Model(&activePass).Association("Pass").Append(&pass)
	if err != nil {
		tx.Rollback()
		SendMessageOnly("Could not associate pass with active pass: "+err.Error(), ctx, 500)
		return
	}

	income := models.IncomeCreate{
		Amount:       activePass.Pass.Price,
		ActivePassID: &activePass.ID,
		PayerID:      payer.ID,
	}
	incomeResult := createIncome(tx, income, openapitypes.UUID(uuid.MustParse(userID)), 0)
	if !incomeResult.Success {
		tx.Rollback()
		SendMessageOnly(incomeResult.Message, ctx, 500)
		return
	}

	tx.Commit()
	SendMessageOnly("Active pass was created successfully", ctx, 201)
}

// MARK: PATCH /active-passes/:id
func UpdateActivePass(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	var activePassUpdate models.ActivePassUpdate
	err := ctx.BindJSON(&activePassUpdate)
	if err != nil {
		SendMessageOnly("Parse error: "+err.Error(), ctx, 400)
		return
	}

	slog.Info("UpdateActivePass", "activePassUpdate", activePassUpdate)

	var activePass models.ActivePass
	id := ctx.Param("id")
	result := DB.First(&activePass, "id = ?", id)
	if result.Error != nil {
		SendMessageOnly("Could not get active pass: "+result.Error.Error(), ctx, 500)
		return
	}

	if activePass.UserID.String() != userID {
		SendMessageOnly("You are not allowed to update this active pass", ctx, 403)
		return
	}

	if activePassUpdate.Occasions != nil {
		activePass.Occasions = *activePassUpdate.Occasions
	}
	if activePassUpdate.Comment != nil {
		activePass.Comment = activePassUpdate.Comment
	}
	if activePassUpdate.ValidFrom != nil {
		roundedTime := activePassUpdate.ValidFrom.Round(time.Second)
		activePass.ValidFrom = &roundedTime
	}
	if activePassUpdate.ValidUntil != nil {
		roundedTime := activePassUpdate.ValidUntil.Round(time.Second)
		activePass.ValidUntil = &roundedTime
	}

	result = DB.Save(&activePass)
	if result.Error != nil {
		SendMessageOnly("Could not update active pass: "+result.Error.Error(), ctx, 500)
		return
	}

	SendMessageOnly("Pass in use was updated successfully", ctx, 200)
}

// MARK: DELETE /active-passes/:id
func DeleteActivePass(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	var activePass models.ActivePass
	id := ctx.Param("id")

	slog.Info("DeleteActivePass", "id", id)

	result := DB.First(&activePass, "id = ?", id)
	if result.Error != nil {
		SendMessageOnly("Could not get active pass: "+result.Error.Error(), ctx, 500)
		return
	}

	if activePass.UserID.String() != userID {
		SendMessageOnly("You are not allowed to delete this active pass", ctx, 403)
		return
	}

	tx := DB.Begin()
	activePass.IsActive = false
	result = tx.Save(&activePass)
	if result.Error != nil {
		tx.Rollback()
		SendMessageOnly("Could not delete active pass: "+result.Error.Error(), ctx, 500)
		return
	}

	tx.Commit()
	ctx.JSON(200, nil)
}
