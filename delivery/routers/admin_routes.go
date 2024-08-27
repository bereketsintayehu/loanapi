package routers

import (
	"loan/config/db"
	"loan/delivery/controllers"
	"loan/infrastracture"
	"loan/repository"
	"loan/usecase"

	"github.com/gin-gonic/gin"
)

func SetUpAdmin(router *gin.Engine) {
	userRepo := repository.NewUserRepositoryImpl(db.UserCollection)

	tokenGen := infrastracture.NewTokenGenerator()
	passwordSvc := infrastracture.NewPasswordService()

	userUsecase := usecase.NewUserUsecase(userRepo, tokenGen, passwordSvc)

	adminController := controllers.NewUserController(userUsecase)

	admin := router.Group("/admin")
	admin.Use(infrastracture.AuthMiddleware())
	{
		admin.GET("/users", adminController.GetUsers)
		admin.DELETE("/users/:id", adminController.DeleteUser)
	}
}
