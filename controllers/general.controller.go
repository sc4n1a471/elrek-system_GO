package controllers

import (
	"elrek-system_GO/models"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var Error error

type ActionResponse struct {
	Success bool
	Message string
}

func SetupDB() error {
	godotenv.Load(".env")
	godotenv.Load("../.env")

	fmt.Println("Connecting to database: " + os.Getenv("DB_HOST") + "...")

	//loc, err := time.LoadLocation("Europe/Budapest")
	//if err != nil {
	//	return err
	//}

	if os.Getenv("DB_HOST") == "" {
		return fmt.Errorf("DB_HOST is not set")
	}

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
		&models.ActivePass{},
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
