package routers

import (
	"github.com/gin-gonic/gin"

	"loan/config/db"
	"loan/delivery/controllers"
	"loan/infrastracture"
	"loan/repository"
	"loan/usecase"
)

func SetUpUser(router *gin.Engine) {
	userRepo := repository.NewUserRepositoryImpl(db.UserCollection)

	tokenGen := infrastracture.NewTokenGenerator()
	passwordService := infrastracture.NewPasswordService()

	userUsecase := usecase.NewUserUsecase(userRepo, tokenGen, passwordService)

	userController := controllers.NewUserController(userUsecase)

	user := router.Group("/users")
	{
		user.POST("/refresh-token", userController.RefreshToken)
		user.POST("/register", userController.Register)
		user.GET("/verify-email/:token/:email", userController.ActivateAccount)
		user.POST("/verify-email", userController.GetNewVerificationEmail)
		user.POST("/login", userController.Login)
		user.POST("/refresh", userController.RefreshToken)
		user.POST("/password-reset", userController.PasswordReset)
		user.POST("/password-update", userController.UpdatePassword)
	}
	user.Use(infrastracture.AuthMiddleware())

	{
		user.GET("/profile", userController.GetMyProfile)
	}
}
