package controllers

import (
	"elrek-system_GO/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
	"golang.org/x/crypto/bcrypt"
)

func Login(ctx *gin.Context) {
	var userLogin models.UserLogin
	if err := ctx.BindJSON(&userLogin); err != nil {
		SendMessageOnly("Parse error: "+err.Error(), ctx, 400)
		return
	}

	var user models.User
	result := DB.First(&user, "email = ?", userLogin.Email)
	if result.Error != nil {
		SendMessageOnly("Could not get user: "+result.Error.Error(), ctx, 500)
		return
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(userLogin.Password)); err != nil {
		// If the email is present in the DB then compare the Passwords and if incorrect password then return error.
		SendMessageOnly("Wrong password", ctx, 401)
		return
	}

	var userLoginResponse models.UserLoginResponse
	userLoginResponse.Email = user.Email
	userLoginResponse.Id = user.Id
	userLoginResponse.Name = user.Name
	userLoginResponse.Token = nil // TODO: generate token

	ctx.JSON(200, userLoginResponse)
}

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

	password, _ := bcrypt.GenerateFromPassword([]byte(userCreate.Password), 14)
	//GenerateFromPassword returns the bcrypt hash of the password at the given cost i.e. (14 in our case).

	var user models.User
	generatedUUID := openapitypes.UUID(uuid.New())
	user.Id = &generatedUUID
	user.Email = &userCreate.Email
	user.Name = &userCreate.Name
	user.OwnerId = &generatedUUID // TODO: change this to the actual owner id
	user.Password = password
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
		password, _ := bcrypt.GenerateFromPassword([]byte(userUpdate.Password), 14)
		user.Password = password
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
