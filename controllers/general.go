package controllers

import (
	"elrek-system_GO/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"os"
)

var DB *gorm.DB
var Error error

type ActionResponse struct {
	Success bool
	Message string
}

func SetupDB() error {
	fmt.Println("Connecting to database...")

	//loc, err := time.LoadLocation("Europe/Budapest")
	//if err != nil {
	//	return err
	//}

	dsn := os.Getenv("DB_USERNAME") +
		":" +
		os.Getenv("DB_PASSWORD") +
		"@tcp(" +
		os.Getenv("DB_HOST") +
		":" +
		os.Getenv("DB_PORT") +
		")/" +
		os.Getenv("DB_NAME") +
		"?parseTime=true" +
		"&loc=Local"

	DB, Error = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if Error != nil {
		return Error
	}

	err := DB.AutoMigrate(
		&models.User{},
		&models.Service{},
		&models.DynamicPrice{},
		&models.Pass{},
		&models.PassInUse{},
		&models.Income{})
	if err != nil {
		fmt.Print(err.Error())
		return err
	}
	fmt.Println("Successfully connected to database!")

	return nil
}

func SendMessageOnly(message string, ctx *gin.Context, statusCode int) {
	ctx.IndentedJSON(statusCode, models.MessageOnlyResponse{
		Message: message,
	})
}

func SendData(message interface{}, ctx *gin.Context) {
	response := map[string]interface{}{
		"message": message,
	}

	ctx.IndentedJSON(http.StatusOK, response)
}
