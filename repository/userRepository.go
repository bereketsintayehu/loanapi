package repository

import (
	"context"
	"errors"
	"loan/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserRepositoryImpl struct {
	collection *mongo.Collection
}

func NewUserRepositoryImpl(coll *mongo.Collection) domain.UserRepository {
	return &UserRepositoryImpl{collection: coll}
}

func (ur *UserRepositoryImpl) Login(user *domain.User) (*domain.User, error) {
	var existingUser domain.User
	err := ur.collection.FindOne(context.Background(), map[string]string{"email": user.Email}).Decode(&existingUser)
	if err != nil {
		return &domain.User{}, err
	}
	return &existingUser, nil

}

func (ur *UserRepositoryImpl) GetUserByID(id string) (domain.User, error) {
	var user domain.User
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.User{}, err
	}

	err = ur.collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (ur *UserRepositoryImpl) DeleteRefreshToken(user *domain.User, token string) error {
	objID, err := primitive.ObjectIDFromHex(user.ID.Hex())
	if err != nil {
		return err
	}
	_, err = ur.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$pull": bson.M{"refresh_tokens": bson.M{"token": token}}},
	)
	return err
}

func (ur *UserRepositoryImpl) UpdateUser(user *domain.User) error {

	_, err := ur.collection.UpdateOne(context.Background(), map[string]string{"email": user.Email}, bson.M{"$set": user})
	return err
}

func (ur *UserRepositoryImpl) DeleteAllRefreshTokens(user *domain.User) error {
	_, err := ur.collection.UpdateOne(context.Background(), map[string]string{"username": user.Username}, bson.M{"$set": bson.M{"refresh_tokens": []domain.RefreshToken{}}})
	return err
}

func (ur *UserRepositoryImpl) IsVerified(userID primitive.ObjectID) (bool, error) {
	var user domain.User
	err := ur.collection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return false, err
	}
	return user.IsActive, nil
}

func (ur *UserRepositoryImpl) Register(user domain.User) error {
	_, err := ur.collection.InsertOne(context.Background(), user)
	return err
}

func (ur *UserRepositoryImpl) GetUserByUsernameOrEmail(username, email string) (domain.User, error) {
	var user domain.User
	err := ur.collection.FindOne(context.Background(), bson.M{"username": username, "email": email}).Decode(&user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (ur *UserRepositoryImpl) GetUserByUsername(username string) (domain.User, error) {
	var user domain.User
	err := ur.collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (ur *UserRepositoryImpl) ActivateAccount(token, email string) error {
	var user domain.User

	err := ur.collection.FindOne(context.Background(), bson.M{"email": email, "activation_token": token}).Decode(&user)

	if err != nil {
		return errors.New("invalid token or user not found")
	}

	if time.Since(user.TokenCreatedAt) > 24*time.Hour {
		return errors.New("token has expired")
	}

	_, err = ur.collection.UpdateOne(context.Background(), bson.M{"activation_token": token}, bson.M{"$set": bson.M{"is_active": true}, "$unset": bson.M{"activation_token": ""}, "$currentDate": bson.M{"updated_at": true}})
	if err != nil {
		return errors.New("failed to activate account")
	}

	return nil

}

func (ur *UserRepositoryImpl) GetUserByEmail(email string) (domain.User, error) {
	var user domain.User
	err := ur.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (ur *UserRepositoryImpl) GetUserByResetToken(token string) (domain.User, error) {
	var user domain.User

	err := ur.collection.FindOne(context.Background(), bson.M{"password_reset_token": token}).Decode(&user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

// password Hashing

func (uc *UserRepositoryImpl) HashPasswordRepo(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (uc *UserRepositoryImpl) CheckPasswordHashRepo(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

func (ur *UserRepositoryImpl) UpdatePassword(token, email, password string) error {
	filter := bson.M{
		"email":                email,
		"password_reset_token": token,
		"token_created_at":     bson.M{"$gte": time.Now().Add(-24 * time.Hour)},
	}

	update := bson.M{
		"$set": bson.M{
			"password": password,
		},
		"$unset": bson.M{
			"password_reset_token": "",
		},
	}

	_, err := ur.collection.UpdateOne(context.Background(), filter, update)
	return err
}
