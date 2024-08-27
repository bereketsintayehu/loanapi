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

type LoanRepository struct {
	collection *mongo.Collection
}

func NewLoanRepository(coll *mongo.Collection) domain.LoanRepository {
	return &LoanRepository{
		collection: coll,
	}
}

func (lr *LoanRepository) CreateLoan(loan domain.Loan) error {
	loan.ID = primitive.NewObjectID().Hex()
	loan.CreatedAt = time.Now()
	loan.UpdatedAt = loan.CreatedAt

	_, err := lr.collection.InsertOne(context.Background(), loan)
	return err
}

func (lr *LoanRepository) GetLoanOfUserByID(loanID, userID primitive.ObjectID) (*domain.Loan, error) {
	var loan domain.Loan
	err := lr.collection.FindOne(context.Background(), bson.M{"_id": loanID, "user_id": userID}).Decode(&loan)
	if err != nil {
		return nil, err
	}
	return &loan, nil
}

func (lr *LoanRepository) GetAllLoans(status string, order string, limit int, offset int) ([]*domain.Loan, int64, error) {
	filter := bson.M{}
	if status != "all" && status != "" {
		filter["status"] = status
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	if order == "asc" {
		findOptions.SetSort(bson.D{{Key: "created_at", Value: 1}})
	} else {
		findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})
	}

	var loans []*domain.Loan
	cursor, err := lr.collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &loans); err != nil {
		return nil, 0, err
	}

	total, err := lr.collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, err
	}

	return loans, total, nil
}

func (lr *LoanRepository) UpdateLoanStatus(loanID primitive.ObjectID, status string, adminID primitive.ObjectID, approvalDate *time.Time, rejectionReason *string) error {
	update := bson.M{
		"$set": bson.M{
			"status":    status,
			"admin_id":  adminID.Hex(),
			"updated_at": time.Now(),
		},
	}

	if status == "approved" && approvalDate != nil {
		update["$set"].(bson.M)["approval_date"] = *approvalDate
	}
	if status == "rejected" && rejectionReason != nil {
		update["$set"].(bson.M)["rejection_reason"] = *rejectionReason
	}

	_, err := lr.collection.UpdateOne(context.Background(), bson.M{"_id": loanID}, update)
	return err
}

func (lr *LoanRepository) DeleteLoan(loanID primitive.ObjectID) error {
	_, err := lr.collection.DeleteOne(context.Background(), bson.M{"_id": loanID})
	return err
}