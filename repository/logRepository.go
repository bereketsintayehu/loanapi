package repository

import (
	"context"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"loan/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"

)

type LogRepository struct {
	collection *mongo.Collection
}

func NewLogRepository(coll *mongo.Collection) domain.LogRepository {
	return &LogRepository{
		collection: coll,
	}
}

func (lr *LogRepository) CreateLog(log domain.Log) error {
	log.ID = primitive.NewObjectID().Hex()
	log.Timestamp = time.Now()

	_, err := lr.collection.InsertOne(context.Background(), log)
	return err
}

func (lr *LogRepository) GetLogs(event string, order string, limit int, offset int) ([]*domain.Log, int64, error) {
	filter := bson.M{}
	if event != "" {
		filter["event"] = event
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	if order == "asc" {
		findOptions.SetSort(bson.D{{Key: "timestamp", Value: 1}})
	} else {
		findOptions.SetSort(bson.D{{Key: "timestamp", Value: -1}})
	}

	var logs []*domain.Log
	cursor, err := lr.collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &logs); err != nil {
		return nil, 0, err
	}

	total, err := lr.collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
