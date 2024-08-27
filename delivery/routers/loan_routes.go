package routers

import (
	"loan/config/db"
	"loan/delivery/controllers"
	"loan/infrastracture"
	"loan/repository"
	"loan/usecase"

	"github.com/gin-gonic/gin"
)

func SetUpLoan(router *gin.Engine) {
	loanRepo := repository.NewLoanRepository(db.LoanCollection)
	logRepo := repository.NewLogRepository(db.LogCollection)

	loanUsecase := usecase.NewLoanUsecase(loanRepo, logRepo)

	loanController := controllers.NewLoanController(loanUsecase)

	loan := router.Group("/")
	loan.Use(infrastracture.AuthMiddleware())
	{
		loan.POST("/loans", loanController.CreateLoan)
		loan.GET("/loans/:loanID", loanController.ViewLoanStatus)
		loan.GET("/admin/loans", loanController.ViewAllLoans)
		loan.PATCH("/admin/loans/:loanID/:status", loanController.PatchLoanStatus)
		loan.DELETE("/admin/loans/:loanID", loanController.DeleteLoan)
	}
}
