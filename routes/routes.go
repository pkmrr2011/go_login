package routes

import (
	"login/controllers"
	"login/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(router *gin.Engine, db *mongo.Database) {

	userController := controllers.NewUserController(db)
	adminController := controllers.NewAdminController()

	userRoute := router.Group("/user")
	userRoute.Use(middleware.Authorization)

	{
		userRoute.GET("/", userController.GetUsers)
		userRoute.POST("/", userController.CreateUser)
		userRoute.GET("/:id", userController.GetUser)
		userRoute.PUT("/:id", userController.UpdateUser)
		userRoute.DELETE("/:id", userController.DeleteUser)
	}

	adminRoute := router.Group("/admin")
	{
		adminRoute.GET("/login", adminController.AdminLogin)
	}
}
