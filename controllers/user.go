package controllers

import (
	"elrek-system_GO/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

var SecretKey = "123456789ABCDEF"

func Login(ctx *gin.Context) {
	var userLogin models.UserLogin
	if err := ctx.BindJSON(&userLogin); err != nil {
		SendMessageOnly("Parse error: "+err.Error(), ctx, 400)
		return
	}

	var user *models.User
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

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.ID.String(),                      //issuer contains the ID of the user.
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //Adds time to the token i.e. 24 hours.
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		SendMessageOnly("Could not login: "+err.Error(), ctx, 500)
		return
	}

	_, err = ctx.Cookie("jwt")
	if err != nil {
		//cookie = "he"
		ctx.SetCookie("jwt", token, 3600, "/", "localhost", false, true)
	}
	ctx.SetCookie("jwt", token, 3600, "/", "localhost", false, true)

	var userLoginResponse models.UserLoginResponse
	userLoginResponse.Email = user.Email
	userLoginResponse.ID = user.ID
	userLoginResponse.Name = user.Name
	userLoginResponse.IsAdmin = user.IsAdmin

	ctx.JSON(200, userLoginResponse)
}

// This function checks if the user is authenticated or not. If yes, returns the following information
//  1. bool: true if authenticated else false
//  2. string: the user ID
//  3. bool: true if the user is admin else false
func CheckAuth(ctx *gin.Context, onlyAdmin bool) (string, bool) {
	cookie, err := ctx.Cookie("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil //using the SecretKey which was generated in th Login function
	})

	if err != nil {
		SendMessageOnly("Not logged in", ctx, 401)
		return "", false
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User
	result := DB.First(&user, "id = ?", claims.Issuer)
	if result.Error != nil {
		SendMessageOnly("Could not get user in checking authentication: "+result.Error.Error(), ctx, 500)
		return "", false
	}

	if onlyAdmin && !user.IsAdmin {
		SendMessageOnly("Access denied", ctx, 403)
		return "", false
	}

	return claims.Issuer, user.IsAdmin
}

func Logout(ctx *gin.Context) {
	_, err := ctx.Cookie("jwt")
	if err != nil {
		SendMessageOnly("Not logged in", ctx, 401)
		return
	}

	ctx.SetCookie("jwt", "", -1, "/", "localhost", false, true)
	SendMessageOnly("Logged out successfully", ctx, 200)
}

func GetUsers(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	var queryParameters map[string][]string = ctx.Request.URL.Query()
	var isActive bool = true
	if len(queryParameters) != 0 {
		isActive = queryParameters["is_active"][0] == "true"
	}

	var users []models.User
	result := DB.Where("is_active = ? and owner_id = ?", isActive, userID).Find(&users)
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
		userResponse.ID = user.ID
		userResponse.IsActive = user.IsActive
		userResponse.IsAdmin = user.IsAdmin
		userResponse.Name = user.Name
		userResponses = append(userResponses, userResponse)
	}

	ctx.JSON(200, userResponses)
}

func GetUser(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, false)
	if userID == "" {
		return
	}

	var user models.User
	id := ctx.Param("id")

	if userID != id {
		userID, _ = CheckAuth(ctx, true)
		if userID == "" {
			return
		}

		result := DB.Where("owner_id = ?", userID).First(&user, "id = ?", id)
		if result.Error != nil {
			SendMessageOnly("Could not get user: "+result.Error.Error(), ctx, 500)
			return
		}
	} else {
		result := DB.First(&user, "id = ?", id)
		if result.Error != nil {
			SendMessageOnly("Could not get user: "+result.Error.Error(), ctx, 500)
			return
		}
	}

	var userResponse models.UserResponse
	userResponse.Email = user.Email
	userResponse.ID = user.ID
	userResponse.IsActive = user.IsActive
	userResponse.IsAdmin = user.IsAdmin
	userResponse.Name = user.Name

	ctx.JSON(200, userResponse)
}

func CreateUserWrapper(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	var userCreate models.UserCreate
	if err := ctx.BindJSON(&userCreate); err != nil {
		SendMessageOnly("Parse error: "+err.Error(), ctx, 400)
		return
	}

	tx := DB.Begin()

	createResult := createUser(userCreate, userID, tx, false)
	if !createResult.Success {
		SendMessageOnly(createResult.Message, ctx, 500)
		tx.Rollback()
		return
	}

	tx.Commit()
	SendMessageOnly(createResult.Message, ctx, 201)
}

// CreateAdminUserWrapper Creates an admin account
// ctx *gin.Context The context of the request
func CreateAdminUserWrapper(ctx *gin.Context) {
	userID := openapitypes.UUID{}.String()

	var userCreate models.UserCreate
	if err := ctx.BindJSON(&userCreate); err != nil {
		SendMessageOnly("Parse error: "+err.Error(), ctx, 400)
		return
	}

	tx := DB.Begin()

	createResult := createUser(userCreate, userID, tx, true)
	if !createResult.Success {
		SendMessageOnly(createResult.Message, ctx, 500)
		tx.Rollback()
		return
	}

	tx.Commit()
	SendMessageOnly(createResult.Message, ctx, 201)
}

// createUser Creates a user |
// userCreate models.UserCreate - The user to be created |
// userID string - The ID of the user |
// tx *gorm.DB - The transaction |
// admin bool - If the user is an admin
func createUser(userCreate models.UserCreate, userID string, tx *gorm.DB, admin bool) ActionResponse {
	// Check if user is already present in the DB
	checkResult := tx.First(&models.User{}, "email = ?", userCreate.Email)
	if checkResult.RowsAffected != 0 {
		return ActionResponse{false, "User with this email already exists"}
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(userCreate.Password), 14)
	//GenerateFromPassword returns the bcrypt hash of the password at the given cost i.e. (14 in our case).

	var user models.User
	user.ID = openapitypes.UUID(uuid.New())
	user.Email = userCreate.Email
	user.Name = userCreate.Name
	user.Password = password
	user.IsAdmin = userCreate.IsAdmin
	user.IsActive = true
	var ownerID openapitypes.UUID
	if admin {
		ownerID = user.ID
		user.IsAdmin = true
	} else {
		ownerID = openapitypes.UUID(uuid.MustParse(userID))
	}
	user.OwnerID = ownerID

	result := tx.Create(&user)
	if result.Error != nil {
		return ActionResponse{false, "Could not create user: " + result.Error.Error()}
	}

	return ActionResponse{true, "User was created successfully"}
}

func UpdateUser(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, false)
	if userID == "" {
		return
	}

	var userUpdate models.UserUpdate
	if err := ctx.BindJSON(&userUpdate); err != nil {
		SendMessageOnly("Parse error: "+err.Error(), ctx, 400)
		return
	}

	var user models.User
	id := ctx.Param("id")
	if userID != id {
		userID, _ = CheckAuth(ctx, true)
		if userID == "" {
			return
		}
	}

	result := DB.First(&user, "id = ?", id)
	if result.Error != nil {
		SendMessageOnly("Could not get existing user: "+result.Error.Error(), ctx, 500)
		return
	}

	//if userUpdate.Email != nil {
	//	user.Email = userUpdate.Email
	//}

	if userUpdate.Name != nil {
		user.Name = *userUpdate.Name
	}

	if userUpdate.Password != "" {
		password, _ := bcrypt.GenerateFromPassword([]byte(userUpdate.Password), 14)
		user.Password = password
	}

	if userUpdate.IsAdmin != nil {
		userIDAdmin, _ := CheckAuth(ctx, true)
		if userIDAdmin == "" {
			return
		}
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
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	var user models.User
	id := ctx.Param("id")
	result := DB.First(&user, "id = ? and owner_id = ?", id, userID)
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

func DeleteUserPermanently(ctx *gin.Context) {
	userID, _ := CheckAuth(ctx, true)
	if userID == "" {
		return
	}

	var user models.User
	id := ctx.Param("id")
	result := DB.First(&user, "id = ? and owner_id = ?", id, userID)
	if result.Error != nil {
		SendMessageOnly("Could not get existing user: "+result.Error.Error(), ctx, 500)
		return
	}

	result = DB.Delete(&user)
	if result.Error != nil {
		SendMessageOnly("Could not delete user: "+result.Error.Error(), ctx, 500)
		return
	}

	SendMessageOnly("User was permanently deleted successfully", ctx, 200)
}
