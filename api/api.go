package api

import (
	"elrek-system_GO/controllers"
	"elrek-system_GO/middlewares"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gopkg.in/Graylog2/go-gelf.v2/gelf"
)

func Api() {
	err := controllers.SetupDB()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	router := SetupRouter()

	slog.Info("Starting server on " + os.Getenv("BACKEND_URL"))

	// TODO: https://github.com/samber/slog-graylog
	graylogHost := os.Getenv("GRAYLOG_HOST")
	gelfWriter, err := gelf.NewTCPWriter(graylogHost)
	if err != nil {
		fmt.Println("error gelf.NewWriter:", err)
		return
	}

	router.Use(middlewares.LoggingMiddleware(gelfWriter))

	router.Run(os.Getenv("BACKEND_URL") + ":3000")
	err = http.ListenAndServe(":3000", router)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

func SetupRouter() *gin.Engine {
	godotenv.Load(".env")
	godotenv.Load("../.env")

	router := gin.Default()
	// gin.SetMode(gin.ReleaseMode)

	slog.Info("Setting up CORS with origin: " + os.Getenv("FRONTEND_URL"))

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("FRONTEND_URL")},
		AllowMethods:     []string{"GET", "HEAD", "OPTIONS", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// PING
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// AUTH
	router.POST("/login", controllers.Login)
	router.GET("/check-permissions", controllers.CheckPermissions)
	router.POST("/logout", controllers.Logout)
	router.POST("/register", controllers.Register)

	// USERS
	router.GET("/users", controllers.GetUsers)
	router.GET("/users/:id", controllers.GetUser)
	router.POST("/users", controllers.CreateUserWrapper)
	router.POST("/users/create-admin", controllers.CreateAdminUserWrapper)
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
	router.GET("/dynamic-prices/:id", controllers.GetDynamicPricesWrapper)
	//router.POST("/dynamic-prices", controllers.CreateDynamicPricesWrapperEndpoint)
	//router.PATCH("/dynamic-prices/:id", controllers.UpdateDynamicPrices)
	//router.DELETE("/dynamic-prices/:id", controllers.DeleteDynamicPrices)

	// PASSES
	router.GET("/passes", controllers.GetPasses)
	router.GET("/passes/:id", controllers.GetPassWrapper)
	router.POST("/passes", controllers.CreatePassWrapper)
	router.PATCH("/passes/:id", controllers.UpdatePass)
	router.DELETE("/passes/:id", controllers.DeletePassWrapper)

	// ACTIVE PASS
	router.GET("/active-passes", controllers.GetActivePasses)
	router.GET("/active-passes/:id", controllers.GetActivePass)
	router.POST("/active-passes", controllers.CreateActivePass)
	router.PATCH("/active-passes/:id", controllers.UpdateActivePass)
	router.DELETE("/active-passes/:id", controllers.DeleteActivePass)
	router.GET("/active-passes/:id/validity", controllers.CheckactivePassValidityWrapper)
	//router.GET("/active-passes/:id/use", controllers.UseactivePassWrapper)

	// INCOMES
	router.GET("/incomes", controllers.GetIncomes)
	router.GET("/incomes/:id", controllers.GetIncome)
	router.POST("/incomes", controllers.CreateIncomeWrapper)
	router.POST("/incomes/multiple-users", controllers.CreateIncomeMultipleUsersWrapper)
	router.PATCH("/incomes/:id", controllers.UpdateIncome)
	router.DELETE("/incomes/:id", controllers.DeleteIncome)

	// STATISTICS
	router.GET("/statistics", controllers.GetStatistics)

	// LOCATIONS
	router.GET("/locations", controllers.GetLocations)
	router.GET("/locations/:id", controllers.GetLocation)
	router.GET("/locations/:id/events", controllers.GetLocationEvents)
	router.POST("/locations", controllers.CreateLocation)
	router.PATCH("/locations/:id", controllers.UpdateLocation)
	router.DELETE("/locations/:id", controllers.DeleteLocation)

	// EVENTS
	router.GET("/events", controllers.GetEvents)
	router.GET("/events/:id", controllers.GetEvent)
	router.POST("/events", controllers.CreateEvent)
	router.PATCH("/events/:id", controllers.UpdateEvent)
	router.DELETE("/events/:id", controllers.DeleteEvent)

	return router
}
