package api

import (
	"elrek-system_GO/controllers"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Api() {
	err := controllers.SetupDB()
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	router := gin.Default()

	// AUTH
	router.POST("/login", controllers.Login)

	// USERS
	router.GET("/users", controllers.GetUsers)
	router.GET("/users/:id", controllers.GetUser)
	router.POST("/users", controllers.CreateUser)
	router.PATCH("/users/:id", controllers.UpdateUser)
	router.DELETE("/users/:id", controllers.DeleteUser)

	router.Run("localhost:3000")
	err = http.ListenAndServe(":3000", router)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
}
