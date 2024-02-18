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

	router := SetupRouter()

	router.Run("localhost:3000")
	err = http.ListenAndServe(":3000", router)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
}

func SetupRouter() *gin.Engine {
	router := gin.Default()
	//gin.SetMode(gin.ReleaseMode)

	// PING
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// AUTH
	router.POST("/login", controllers.Login)
	router.POST("/logout", controllers.Logout)

	// USERS
	router.GET("/users", controllers.GetUsers)
	router.GET("/users/:id", controllers.GetUser)
	router.POST("/users", controllers.CreateUserWrapper)
	router.POST("/users/create_admin", controllers.CreateAdminUserWrapper)
	router.PATCH("/users/:id", controllers.UpdateUser)
	router.DELETE("/users/:id", controllers.DeleteUser)
	router.DELETE("/users/permanently/:id", controllers.DeleteUserPermanently)

	// SERVICES
	router.GET("/services", controllers.GetServices)
	router.GET("/services/:id", controllers.GetService)
	router.POST("/services", controllers.CreateServiceWrapper)
	router.PATCH("/services/:id", controllers.UpdateService)
	router.DELETE("/services/:id", controllers.DeleteServiceWrapper)

	// DYNAMIC PRICES
	router.GET("/dynamic_prices/:id", controllers.GetDynamicPricesWrapper)
	//router.POST("/dynamic_prices", controllers.CreateDynamicPricesWrapperEndpoint)
	//router.PATCH("/dynamic_prices/:id", controllers.UpdateDynamicPrices)
	//router.DELETE("/dynamic_prices/:id", controllers.DeleteDynamicPrices)

	// PASSES
	router.GET("/passes", controllers.GetPasses)
	router.GET("/passes/:id", controllers.GetPass)
	router.POST("/passes", controllers.CreatePassWrapper)
	router.PATCH("/passes/:id", controllers.UpdatePass)
	router.DELETE("/passes/:id", controllers.DeletePassWrapper)

	// PASSES IN USE
	router.GET("/passes_in_use", controllers.GetPassesInUse)
	router.GET("/passes_in_use/:id", controllers.GetPassInUse)
	router.POST("/passes_in_use", controllers.CreatePassInUse)
	router.PATCH("/passes_in_use/:id", controllers.UpdatePassInUse)
	router.DELETE("/passes_in_use/:id", controllers.DeletePassInUse)
	router.GET("/passes_in_use/:id/validity", controllers.CheckPassInUseValidityWrapper)
	router.GET("/passes_in_use/:id/use", controllers.UsePassInUseWrapper)

	return router
}
