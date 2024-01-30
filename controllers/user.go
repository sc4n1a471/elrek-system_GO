package controllers

import (
	"elrek-system_GO/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
)

func GetUsers(ctx *gin.Context) {
	var queryParameters map[string][]string = ctx.Request.URL.Query()
	var isActive bool = true
	if len(queryParameters) != 0 {
		isActive = queryParameters["is_active"][0] == "true"
	}

	var users []models.User
	result := DB.Where("is_active = ?", isActive).Find(&users)
	if result.Error != nil {
		SendMessageOnly("Could not get users: "+result.Error.Error(), ctx, 500)
		return
	}

	if len(users) == 0 {
		SendMessageOnly("No users found", ctx, 404)
		return
	}

	var userResponses []models.UserResponse
	for _, user := range users {
		var userResponse models.UserResponse
		userResponse.Email = user.Email
		userResponse.Id = user.Id
		userResponse.IsActive = user.IsActive
		userResponse.IsAdmin = user.IsAdmin
		userResponse.Name = user.Name
		userResponses = append(userResponses, userResponse)
	}

	ctx.JSON(200, userResponses)
}

func GetUser(ctx *gin.Context) {
	var user models.User
	id := ctx.Param("id")
	result := DB.First(&user, "id = ?", id)
	if result.Error != nil {
		SendMessageOnly("Could not get user: "+result.Error.Error(), ctx, 500)
		return
	}

	var userResponse models.UserResponse
	userResponse.Email = user.Email
	userResponse.Id = user.Id
	userResponse.IsActive = user.IsActive
	userResponse.IsAdmin = user.IsAdmin
	userResponse.Name = user.Name

	ctx.JSON(200, userResponse)
}

func CreateUser(ctx *gin.Context) {
	var userCreate models.UserCreate
	if err := ctx.BindJSON(&userCreate); err != nil {
		SendMessageOnly("Parse error: "+err.Error(), ctx, 400)
		return
	}

	var user models.User
	generatedUUID := openapitypes.UUID(uuid.New())
	user.Id = &generatedUUID
	user.Email = &userCreate.Email
	user.Name = &userCreate.Name
	user.OwnerId = &generatedUUID // TODO: change this to the actual owner id
	user.Password = &userCreate.Password
	user.IsAdmin = userCreate.IsAdmin
	user.IsActive = true

	tx := DB.Begin()

	result := tx.Where("email = ?", user.Email).First(&user)
	if result.Error == nil {
		tx.Rollback()
		SendMessageOnly("User with this email already exists", ctx, 400)
		return
	}

	result = tx.Create(&user)
	if result.Error != nil {
		tx.Rollback()
		SendMessageOnly("Could not create user: "+result.Error.Error(), ctx, 500)
		return
	}

	tx.Commit()
	SendMessageOnly("User was created successfully", ctx, 200)
}

func UpdateUser(ctx *gin.Context) {
	var userUpdate models.UserUpdate
	if err := ctx.BindJSON(&userUpdate); err != nil {
		SendMessageOnly("Parse error: "+err.Error(), ctx, 400)
		return
	}

	var user models.User
	id := ctx.Param("id")
	result := DB.First(&user, "id = ?", id)
	if result.Error != nil {
		SendMessageOnly("Could not get existing user: "+result.Error.Error(), ctx, 500)
		return
	}

	if userUpdate.Email != nil {
		user.Email = userUpdate.Email
	}

	if userUpdate.Name != nil {
		user.Name = userUpdate.Name
	}

	if userUpdate.Password != "" {
		user.Password = &userUpdate.Password
	}

	if userUpdate.IsAdmin != nil {
		user.IsAdmin = *userUpdate.IsAdmin
	}

	result = DB.Save(&user)
	if result.Error != nil {
		SendMessageOnly("Could not update user: "+result.Error.Error(), ctx, 500)
		return
	}

	SendMessageOnly("User was updated successfully", ctx, 200)
}

func DeleteUser(ctx *gin.Context) {
	var user models.User
	id := ctx.Param("id")
	result := DB.First(&user, "id = ?", id)
	if result.Error != nil {
		SendMessageOnly("Could not get existing user: "+result.Error.Error(), ctx, 500)
		return
	}

	user.IsActive = false

	result = DB.Save(&user)
	if result.Error != nil {
		SendMessageOnly("Could not delete user: "+result.Error.Error(), ctx, 500)
		return
	}

	SendMessageOnly("User was deleted successfully", ctx, 200)
}
