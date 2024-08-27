package usecase

import (
	"loan/domain"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LogUsecase struct {
	logRepo domain.LogRepository
}

func NewLogUsecase(lr domain.LogRepository) *LogUsecase {
	return &LogUsecase{
		logRepo: lr,
	}
}

func (lu *LogUsecase) CreateLog(log domain.Log) *domain.CustomError {
	var wg sync.WaitGroup
	var err error

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = lu.logRepo.CreateLog(log)
	}()

	wg.Wait()

	if err != nil {
		return domain.ErrLogCreationFailed
	}

	return &domain.CustomError{}
}

func (lu *LogUsecase) ViewLogs(event string, order string, limit int, offset int, adminID string) ([]*domain.Log, int64, *domain.CustomError) {
	var wg sync.WaitGroup
	var logs []*domain.Log
	var total int64
	var err error

	wg.Add(2)

	go func() {
		defer wg.Done()
		logs, total, err = lu.logRepo.GetLogs(event, order, limit, offset)
	}()

	adminObjectID, err := primitive.ObjectIDFromHex(adminID)
	if err != nil {
		return nil, 0, domain.ErrInvalidUserID
	}

	go func() {
		defer wg.Done()
		// Log the fetching event
		logEvent := domain.Log{
			Event:   "view_logs",
			Details: "Logs viewed bh Admin ID: " + adminID,
			UserID: adminObjectID,
		}
		err = lu.logRepo.CreateLog(logEvent)
		if err != nil {
			return 
		}
	}()

	wg.Wait()

	if err != nil {
		return nil, 0, domain.ErrFetchingLogsFailed
	}

	return logs, total, &domain.CustomError{}
}