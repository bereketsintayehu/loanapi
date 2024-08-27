package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Loan struct {
    ID              string    `json:"id"`
    UserID          primitive.ObjectID    `json:"user_id"`
    Amount          float64   `json:"amount"`
    InterestRate    float64   `json:"interest_rate"`
    Term            int       `json:"term"`
    Status          string    `json:"status"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
    Reason          string    `json:"reason,omitempty"`
    AdminID         primitive.ObjectID    `json:"admin_id,omitempty"`
    ApprovalDate    *time.Time `json:"approval_date,omitempty"`
    RejectionReason string    `json:"rejection_reason,omitempty"`
}

type LoanUsecase interface {
	CreateLoan(loan Loan, userID string) *CustomError
	ViewLoanStatus(loanID, userID string) (*Loan, *CustomError)
	ViewAllLoans(status string, order string, limit int, offset int, adminID string) ([]*Loan, int64, *CustomError)
	UpdateLoanStatus(loanID primitive.ObjectID, status string, adminID primitive.ObjectID, approvalDate *time.Time, rejectionReason *string) error
	DeleteLoan(loanID, adminId string) *CustomError
}

type LoanRepository interface {
	CreateLoan(loan Loan) error
	GetLoanOfUserByID(loanID, userID primitive.ObjectID) (*Loan, error)
	GetAllLoans(status string, order string, limit int, offset int) ([]*Loan, int64, error)
	UpdateLoanStatus(loanID primitive.ObjectID, status string, adminID primitive.ObjectID, approvalDate *time.Time, rejectionReason *string) error
	DeleteLoan(loanID primitive.ObjectID) error
}
