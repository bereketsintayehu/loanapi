package routers

import (
	"loan/config/db"
	"loan/delivery/controllers"
	"loan/infrastracture"
	"loan/repository"
	"loan/usecase"

	"github.com/gin-gonic/gin"
)

func SetUpLog(router *gin.Engine) {
	logRepo := repository.NewLogRepository(db.LogCollection)

	logUsecase := usecase.NewLogUsecase(logRepo)

	logController := controllers.NewLogController(logUsecase)

	log := router.Group("/")
	log.Use(infrastracture.AuthMiddleware())
	{
		log.GET("/admin/logs", logController.ViewLogs)
	}
}
