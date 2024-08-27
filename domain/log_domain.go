package domain

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Log struct {
	ID        string             `json:"id"`
	Timestamp time.Time          `json:"timestamp"`
	UserID    primitive.ObjectID          `json:"user_id,omitempty"`
	Event     string             `json:"event"`
	Details   string             `json:"details,omitempty"`
}

type LogUsecase interface {
	CreateLog(log Log) *CustomError
	ViewLogs(event string, order string, limit int, offset int) ([]*Log, int64, *CustomError)
}

type LogRepository interface {
	CreateLog(log Log) error
	GetLogs(event string, order string, limit int, offset int) ([]*Log, int64, error)
}
