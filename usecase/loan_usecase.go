package usecase

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"loan/domain"
	"sync"
	"time"
)

type LoanUsecase struct {
	loanRepo domain.LoanRepository
	logRepo  domain.LogRepository
}

func NewLoanUsecase(lr domain.LoanRepository, logRepo domain.LogRepository) *LoanUsecase {
	return &LoanUsecase{
		loanRepo: lr,
		logRepo:  logRepo,
	}
}

func handleGoroutineErrors(errChan <-chan error) *domain.CustomError {
	for err := range errChan {
		if err != nil {
			return domain.ErrOperationFailed
		}
	}
	return &domain.CustomError{}
}

func (lu *LoanUsecase) CreateLoan(loan domain.Loan, userID string) *domain.CustomError {
	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	loan.CreatedAt = time.Now()
	loan.UpdatedAt = time.Now()

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return domain.ErrInvalidUserID
	}

	loan.UserID = userObjID

	wg.Add(2)

	go func() {
		defer wg.Done()
		err := lu.loanRepo.CreateLoan(loan)
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		logEvent := domain.Log{
			Event:   "create_loan",
			Details: "Loan created by user: " + userObjID.Hex(),
			UserID:  userObjID,
		}
		err := lu.logRepo.CreateLog(logEvent)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	return handleGoroutineErrors(errChan)
}

func (lu *LoanUsecase) ViewLoanStatus(loanID, userID string) (*domain.Loan, *domain.CustomError) {
	var wg sync.WaitGroup
	errChan := make(chan error, 2)
	var loan *domain.Loan

	loanObjID, err := primitive.ObjectIDFromHex(loanID)
	if err != nil {
		return nil, domain.ErrInvalidLoanID
	}

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, domain.ErrInvalidUserID
	}

	wg.Add(2)

	go func() {
		defer wg.Done()
		loan, err = lu.loanRepo.GetLoanOfUserByID(loanObjID, userObjID)
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		logEvent := domain.Log{
			Event:   "view_loan_status",
			Details: "Viewed loan status for loan ID: " + loanObjID.Hex() + " by user ID: " + userObjID.Hex(),
			UserID:  userObjID,
		}
		err := lu.logRepo.CreateLog(logEvent)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	if err := handleGoroutineErrors(errChan); err.Message != "" {
		return nil, err
	}

	return loan, &domain.CustomError{}
}

func (lu *LoanUsecase) ViewAllLoans(status string, order string, limit int, offset int, adminID string) ([]*domain.Loan, int64, *domain.CustomError) {
	var wg sync.WaitGroup
	errChan := make(chan error, 2)
	var loans []*domain.Loan
	var total int64

	adminObjID, err := primitive.ObjectIDFromHex(adminID)
	if err != nil {
		return nil, 0, domain.ErrInvalidUserID
	}

	wg.Add(2)

	go func() {
		defer wg.Done()
		loans, total, err = lu.loanRepo.GetAllLoans(status, order, limit, offset)
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		logEvent := domain.Log{
			Event:   "view_all_loans",
			Details: "Viewed all loans by Admin ID: " + adminObjID.Hex(),
			UserID:  adminObjID,
		}
		err := lu.logRepo.CreateLog(logEvent)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	if err := handleGoroutineErrors(errChan); err.Message != "" {
		return nil, 0, err
	}

	return loans, total, &domain.CustomError{}
}

func (uc *LoanUsecase) UpdateLoanStatus(loanID primitive.ObjectID, status string, adminID primitive.ObjectID, approvalDate *time.Time, rejectionReason *string) *domain.CustomError {
	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	wg.Add(2)

	go func() {
		defer wg.Done()
		err := uc.loanRepo.UpdateLoanStatus(loanID, status, adminID, approvalDate, rejectionReason)
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		logEvent := domain.Log{
			Event:   "update_loan_status",
			Details: "Loan status updated to " + status + " for Loan ID: " + loanID.Hex() + " by Admin ID: " + adminID.Hex(),
			UserID:  adminID,
		}
		err := uc.logRepo.CreateLog(logEvent)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	return handleGoroutineErrors(errChan)
}

func (lu *LoanUsecase) DeleteLoan(loanID, adminID string) *domain.CustomError {
	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	loanObjID, err := primitive.ObjectIDFromHex(loanID)
	if err != nil {
		return domain.ErrInvalidLoanID
	}

	adminObjID, err := primitive.ObjectIDFromHex(adminID)
	if err != nil {
		return domain.ErrInvalidUserID
	}

	wg.Add(2)

	go func() {
		defer wg.Done()
		err := lu.loanRepo.DeleteLoan(loanObjID)
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		logEvent := domain.Log{
			Event:   "delete_loan",
			Details: "Loan deleted with ID: " + loanObjID.Hex() + " by Admin ID: " + adminObjID.Hex(),
			UserID:  adminObjID,
		}
		err := lu.logRepo.CreateLog(logEvent)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	return handleGoroutineErrors(errChan)
}
