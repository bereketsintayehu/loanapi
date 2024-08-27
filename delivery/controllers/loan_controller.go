package controllers

import (
	"github.com/gin-gonic/gin"
	"loan/domain"
	"loan/usecase"
	"net/http"
	"strconv"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"

)

type LoanController struct {
	loanUsecase *usecase.LoanUsecase
}

func NewLoanController(lu *usecase.LoanUsecase) *LoanController {
	return &LoanController{
		loanUsecase: lu,
	}
}

// CreateLoan handles the creation of a new loan.
func (lc *LoanController) CreateLoan(c *gin.Context) {
	var loan domain.Loan

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: admin ID not found"})
		return
	}

	if err := c.BindJSON(&loan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if uerr := lc.loanUsecase.CreateLoan(loan, userID.(string)); uerr.Message != "" {
		c.JSON(uerr.StatusCode, gin.H{"error": uerr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Loan created successfully"})
}

// ViewLoanStatus retrieves and returns the status of a specific loan.
func (lc *LoanController) ViewLoanStatus(c *gin.Context) {
	loanID := c.Param("loanID")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: admin ID not found"})
		return
	}

	loan, uerr := lc.loanUsecase.ViewLoanStatus(loanID, userID.(string))
	if uerr.Message != "" {
		c.JSON(uerr.StatusCode, gin.H{"error": uerr.Message})
		return
	}


	if loan == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Loan not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"loan": loan})
}

// ViewAllLoans retrieves and returns all loans with pagination.
func (lc *LoanController) ViewAllLoans(c *gin.Context) {
	Role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: role not found"})
		return
	}

	if Role != "admin" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: admin ID not found"})
		return
	}

	status := c.Query("status")
	order := c.Query("order")

	const defaultLimit, defaultOffset = 10, 0

	limit, err := strconv.ParseInt(c.DefaultQuery("limit", strconv.Itoa(defaultLimit)), 10, 64)
	if err != nil {
		limit = defaultLimit
	}

	offset, err := strconv.ParseInt(c.DefaultQuery("offset", strconv.Itoa(defaultOffset)), 10, 64)
	if err != nil {
		offset = defaultOffset
	}

	loans, total, uerr := lc.loanUsecase.ViewAllLoans(status, order, int(limit), int(offset), adminID.(string))
	if uerr.Message != "" {
		c.JSON(uerr.StatusCode, gin.H{"error": uerr.Message})
		return
	}

	totalPages := (total + limit - 1) / limit 


	c.JSON(http.StatusOK, gin.H{"loans": loans, "current_page": offset, "per_page": limit, "total": total, "total_pages": totalPages})
}

func (lc *LoanController) DeleteLoan(c *gin.Context) {
	loanID := c.Param("loanID")

	Role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: role not found"})
		return
	}

	if Role != "admin" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: admin ID not found"})
		return
	}

	if uerr := lc.loanUsecase.DeleteLoan(loanID, adminID.(string)); uerr.Message != "" {
		c.JSON(uerr.StatusCode, gin.H{"error": uerr.Message})
		return
	}


	c.JSON(http.StatusOK, gin.H{"message": "Loan deleted successfully"})
}

func (lc *LoanController) PatchLoanStatus(c *gin.Context) {
	loanID := c.Param("loanID")

	Role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: role not found"})
		return
	}

	if Role != "admin" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: admin ID not found"})
		return
	}

	var request struct {
		Status          string `json:"status"`
		RejectionReason string `json:"rejection_reason,omitempty"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate the status
	if request.Status != "approved" && request.Status != "rejected" && request.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
		return
	}

	// Set up variables to pass to the use case
	var approvalDate *time.Time
	var rejectionReason *string

	if request.Status == "approved" {
		now := time.Now()
		approvalDate = &now
		rejectionReason = nil
	} else if request.Status == "rejected" {
		if request.RejectionReason == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Rejection reason must be provided when status is rejected"})
			return
		}
		rejectionReason = &request.RejectionReason
		approvalDate = nil
	} else {
		approvalDate = nil
		rejectionReason = nil
	}

	// Convert loanID to ObjectID
	loanObjectID, err := primitive.ObjectIDFromHex(loanID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
		return
	}

	// Convert adminID to ObjectID
	adminObjectID, err := primitive.ObjectIDFromHex(adminID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid admin ID"})
		return
	}

	// Update the loan status in the repository through the use case
	err = lc.loanUsecase.UpdateLoanStatus(loanObjectID, request.Status, adminObjectID, approvalDate, rejectionReason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update loan status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Loan status updated successfully"})
}
