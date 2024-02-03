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
	router.POST("/logout", controllers.Logout)

	// USERS
	router.GET("/users", controllers.GetUsers)
	router.GET("/users/:id", controllers.GetUser)
	router.POST("/users", controllers.CreateUser)
	router.PATCH("/users/:id", controllers.UpdateUser)
	router.DELETE("/users/:id", controllers.DeleteUser)

	// SERVICES
	router.GET("/services", controllers.GetServices)
	router.GET("/services/:id", controllers.GetService)
	router.POST("/services", controllers.CreateServiceWrapper)
	router.PATCH("/services/:id", controllers.UpdateService)
	router.DELETE("/services/:id", controllers.DeleteServiceWrapper)

	// DYNAMIC PRICES
	router.GET("/dynamic_prices/:id", controllers.GetDynamicPrices)
	//router.POST("/dynamic_prices", controllers.CreateDynamicPricesWrapperEndpoint)
	//router.PATCH("/dynamic_prices/:id", controllers.UpdateDynamicPrices)
	//router.DELETE("/dynamic_prices/:id", controllers.DeleteDynamicPrices)

	router.Run("localhost:3000")
	err = http.ListenAndServe(":3000", router)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
}
