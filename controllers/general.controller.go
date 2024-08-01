package controllers

import (
	"elrek-system_GO/models"
	"fmt"
	"log/slog"
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

// MARK: DB
func SetupDB() error {
	godotenv.Load(".env")
	godotenv.Load("../.env")

	slog.Info("Connecting to database: " + os.Getenv("DB_HOST") + "...")

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

	// MARK: Migration
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
	slog.Info("Successfully connected to database!")

	return nil
}

// MARK: Senders
func SendMessageOnly(message string, ctx *gin.Context, statusCode int) {
	slog.Info("Message sent", "message", message)
	ctx.IndentedJSON(statusCode, models.MessageOnlyResponse{
		Message: message,
	})
}

func SendData(message interface{}, ctx *gin.Context) {
	response := map[string]interface{}{
		"message": message,
	}
	slog.Info("Data sent", "response", response)

	ctx.IndentedJSON(http.StatusOK, response)
}
