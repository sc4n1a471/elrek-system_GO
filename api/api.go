package api

import (
	"elrek-system_GO/models"
	"fmt"

	//"elrek-system_GO/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"os"
)

var DB *gorm.DB
var Error error

func Api() {
	dsn := os.Getenv("DB_USERNAME") +
		":" +
		os.Getenv("DB_PASSWORD") +
		"@tcp(" +
		os.Getenv("DB_HOST") +
		":" +
		os.Getenv("DB_PORT") +
		")/" +
		os.Getenv("DB_NAME") +
		"?parseTime=true"

	DB, Error = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if Error != nil {
		return
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
		return
	}

	router := gin.Default()
	//router.GET("/cars/:license_plate", getCar)
	//router.GET("/cars", getCars)
	//router.GET("/cars_all_data", getCarsAllData)
	//router.POST("/cars", createCar)
	//router.PUT("/cars", updateCar)
	//router.DELETE("/cars/:license_plate", deleteCar)
	//
	//router.POST("/license_plate", createLicensePlate)
	//router.PUT("/license_plate/:license_plate", updateLicensePLate)
	//
	//router.GET("/inspections/:license_plate", getInspections)
	//router.POST("/inspections", createInspections)
	//router.DELETE("/inspections/:license_plate", deleteInspections)
	//
	//router.GET("/coordinates", getCoordinates)

	//router.Run("localhost:3000")
	err = http.ListenAndServe(":3000", router)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
}

func sendError(error string, ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusConflict, models.MessageOnlyResponse{
		Message: &error,
	})
}

func sendData(message interface{}, ctx *gin.Context) {
	response := map[string]interface{}{
		"message": message,
	}

	ctx.IndentedJSON(http.StatusOK, response)
}
