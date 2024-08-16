package controllers

import (
	"database/sql"
	"elrek-system_GO/models"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	openapitypes "github.com/oapi-codegen/runtime/types"
)

// MARK: Data statistics
func GetStatistics(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, false)
	if userID == "" {
		return
	}

	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	// #region basic all-time statistics
	// MARK: count all users
	var userCount int64
	err := DB.Model(&models.User{}).Where("owner_id = ?", userID).Count(&userCount).Error
	if err != nil {
		slog.Error(err.Error())
		SendMessageOnly("Internal server error: could not get users", ctx, http.StatusInternalServerError)
		return
	}

	// MARK: count all active passes
	var activePassCount int64
	err = DB.
		Model(&models.Pass{}).
		Where("is_active = ?", true).
		Count(&activePassCount).
		Error
	if err != nil {
		slog.Error(err.Error())
		SendMessageOnly("Internal server error: could not get active passes", ctx, http.StatusInternalServerError)
		return
	}

	// MARK: sum all incomes
	var incomeSum sql.NullInt64
	err = DB.
		Model(&models.Income{}).
		Select("SUM(amount)").
		Where("created_at >= ? AND created_at <= ? AND user_id = ?", startDate, endDate, userID).
		Scan(&incomeSum).
		Error
	if err != nil {
		slog.Error(err.Error())
		SendMessageOnly("Internal server error: could not get income sum", ctx, http.StatusInternalServerError)
		return
	}
	// #endregion

	currentTime := time.Now()
	// #region yearly statistics
	var everyYearIncomeSum []models.EveryYearIncomeSum
	for i := currentTime.Year(); i >= 2021; i-- {
		yearStartDate := time.Date(i, time.January, 1, 0, 0, 0, 0, currentTime.Location())
		yearEndDate := time.Date(i+1, time.January, 1, 0, 0, 0, 0, currentTime.Location())

		var yearlyIncomeSum models.EveryYearIncomeSum
		err = DB.
			Model(&models.Income{}).
			Select("SUM(amount)").
			Where("created_at >= ? AND created_at <= ? AND user_id = ?", yearStartDate, yearEndDate, userID).
			Scan(&yearlyIncomeSum.Sum).
			Error
		if err != nil {
			slog.Error(err.Error())
			SendMessageOnly("Internal server error: could not get yearly income sum", ctx, http.StatusInternalServerError)
			return
		}
		yearlyIncomeSum.Year = i

		everyYearIncomeSum = append(everyYearIncomeSum, yearlyIncomeSum)
	}

	// #endregion

	// #region monthly statistics
	var everyMonthIncomeSum []models.EveryMonthIncomeSum
	for i := 1; i <= 12; i++ {
		monthStartDate := time.Date(currentTime.Year(), time.Month(i), 1, 0, 0, 0, 0, currentTime.Location())
		monthEndDate := time.Date(currentTime.Year(), time.Month(i)+1, 1, 0, 0, 0, 0, currentTime.Location())

		var monthlyIncomeSum models.EveryMonthIncomeSum
		err = DB.
			Model(&models.Income{}).
			Select("SUM(amount)").
			Where("created_at >= ? AND created_at <= ? AND user_id = ?", monthStartDate, monthEndDate, userID).
			Scan(&monthlyIncomeSum.Sum).
			Error
		if err != nil {
			slog.Error(err.Error())
			SendMessageOnly("Internal server error: could not get monthly income sum", ctx, http.StatusInternalServerError)
			return
		}
		monthlyIncomeSum.Month = monthStartDate.Month().String()

		everyMonthIncomeSum = append(everyMonthIncomeSum, monthlyIncomeSum)
	}
	// #endregion

	// MARK: sum all paid/unpaid incomes
	var paidIncomeSum sql.NullInt64
	err = DB.
		Model(&models.Income{}).
		Select("SUM(amount)").
		Where("created_at >= ? AND created_at <= ? AND is_paid = ? AND user_id = ?", startDate, endDate, true, userID).
		Scan(&paidIncomeSum).
		Error
	if err != nil {
		slog.Error(err.Error())
		SendMessageOnly("Internal server error: could not get paid income sum", ctx, http.StatusInternalServerError)
		return
	}

	var unpaidIncomeSum sql.NullInt64
	err = DB.
		Model(&models.Income{}).
		Select("SUM(amount)").
		Where("created_at >= ? AND created_at <= ? AND is_paid = ? AND user_id = ?", startDate, endDate, false, userID).
		Scan(&unpaidIncomeSum).
		Error
	if err != nil {
		slog.Error(err.Error())
		SendMessageOnly("Internal server error: could not get unpaid income sum", ctx, http.StatusInternalServerError)
		return
	}

	// MARK: count all paid/unpaid incomes
	var paidIncomeCount int64
	err = DB.
		Model(&models.Income{}).
		Where("created_at >= ? AND created_at <= ? AND is_paid = ?", startDate, endDate, true).
		Count(&paidIncomeCount).
		Error
	if err != nil {
		slog.Error(err.Error())
		SendMessageOnly("Internal server error: could not get paid income count", ctx, http.StatusInternalServerError)
		return
	}

	var unpaidIncomeCount int64
	err = DB.
		Model(&models.Income{}).
		Where("created_at >= ? AND created_at <= ? AND is_paid = ?", startDate, endDate, false).
		Count(&unpaidIncomeCount).
		Error
	if err != nil {
		slog.Error(err.Error())
		SendMessageOnly("Internal server error: could not get unpaid income count", ctx, http.StatusInternalServerError)
		return
	}

	// MARK: sum all incomes by service
	var incomeByService []models.IncomeByService
	err = DB.
		Model(&models.Income{}).
		Select("services.name as name, SUM(amount) as sum").
		Joins("left join services on services.id = incomes.service_id").
		Where("services.created_at >= ? AND services.created_at <= ? AND incomes.user_id = ?", startDate, endDate, userID).
		Group("services.name").
		Scan(&incomeByService).
		Error
	if err != nil {
		slog.Error(err.Error())
		SendMessageOnly("Internal server error: could not get income by service", ctx, http.StatusInternalServerError)
		return
	}

	// MARK: sum all incomes by user
	var incomeByUser []models.IncomeByUser
	err = DB.
		Model(&models.Income{}).
		Select("users.name as name, SUM(amount) as sum").
		Joins("left join users on users.id = incomes.payer_id").
		Where("users.created_at >= ? AND users.created_at <= ? AND incomes.user_id = ?", startDate, endDate, userID).
		Group("users.name").
		Scan(&incomeByUser).
		Error
	if err != nil {
		slog.Error(err.Error())
		SendMessageOnly("Internal server error: could not get income by user", ctx, http.StatusInternalServerError)
		return
	}

	// MARK: sum all incomes by activePass
	var incomeByActivePass []models.IncomeByActivePass
	err = DB.
		Model(&models.Income{}).
		Select("passes.name as name, SUM(amount) as sum").
		Joins("left join active_passes on active_passes.id = incomes.active_pass_id").
		Joins("left join passes on passes.id = active_passes.pass_id").
		Where("active_passes.created_at >= ? AND active_passes.created_at <= ? AND active_pass_id != ? AND incomes.user_id = ?", startDate, endDate, openapitypes.UUID{}, userID).
		Group("passes.name").
		Scan(&incomeByActivePass).
		Error
	if err != nil {
		slog.Error(err.Error())
		SendMessageOnly("Internal server error: could not get income by active pass", ctx, http.StatusInternalServerError)
		return
	}

	statistics := models.Statistics{
		UserCount:           userCount,
		ActivePassCount:     activePassCount,
		IncomeSum:           incomeSum,
		PaidIncomeSum:       paidIncomeSum,
		UnpaidIncomeSum:     unpaidIncomeSum,
		EveryYearIncomeSum:  everyYearIncomeSum,
		EveryMonthIncomeSum: everyMonthIncomeSum,
		IncomesByService:    incomeByService,
		IncomesByUser:       incomeByUser,
		IncomesByActivePass: incomeByActivePass,
	}

	ctx.JSON(http.StatusOK, statistics)
}
