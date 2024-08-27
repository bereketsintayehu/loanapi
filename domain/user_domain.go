package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


type User struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	Username  string              `bson:"username" json:"username"`
	Email     string              `bson:"email" json:"email"`
	Name      string              `bson:"name" json:"name"`
	Password  string              `bson:"password" json:"password"`
	Bio       string              `bson:"bio,omitempty" json:"bio,omitempty"`
	Role      string           `bson:"role" json:"role"`
	CreatedAt primitive.Timestamp `bson:"createdAt" json:"createdAt"`
	UpdatedAt primitive.Timestamp `bson:"updatedAt" json:"updatedAt"`

	ActivationToken string         `bson:"activation_token,omitempty" json:"activation_token,omitempty"`
	TokenCreatedAt  time.Time      `bson:"token_created_at"`
	IsActive        bool           `bson:"is_active"`
	RefreshTokens   []RefreshToken `bson:"refresh_tokens" json:"refresh_tokens"`

	PasswordResetToken string `bson:"password_reset_token,omitempty" json:"password_reset_token,omitempty"`
}

type RefreshToken struct {
	Token     string    `bson:"token" json:"token"`
	DeviceID  string    `bson:"device_id" json:"device_id"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

type LogInResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

type TokenGenerator interface {
	GenerateToken(user User) (string, error)
	GenerateRefreshToken(user User) (string, error)
	RefreshToken(token string) (string, error)
}

type TokenVerifier interface {
	VerifyToken(token string) (*User, error)
	VerifyRefreshToken(token string) (*User, error)
}

type PasswordService interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type UserUsecase interface {
	// for every user
	Login(user *User, deviceID string) (LogInResponse, *CustomError)
	RefreshToken(userID, deviceID, token string) (LogInResponse, *CustomError)
	Register(user User) *CustomError
	GetUserByUsernameOrEmail(username, email string) (User, *CustomError)
	GetNewVerificationEmail(email string) *CustomError
	ActivateAccount(token, email string) *CustomError
	UpdatePassword(PasswordUpdateRequest) *CustomError
	PasswordReset(email string) *CustomError

	GetMyProfile(userID string) (User, *CustomError)
	GetUsers() ([]User, *CustomError)
	DeleteUser(userID string) (User, *CustomError)
}

type UserRepository interface {
	Login(user *User) (*User, error)
	Register(user User) error
	GetUserByUsernameOrEmail(username, email string) (User, error)
	ActivateAccount(token, email string) error
	UpdateUser(user *User) error
	GetUserByID(id string) (User, error)
	GetUserByResetToken(token string) (User, error)
	GetUserByEmail(email string) (User, error)
	GetUserByUsername(username string) (User, error)
	UpdatePassword(token, email, password string) error
	GetMyProfile(userID string) (User, error)
	GetUsers() ([]User, error)
	DeleteUser(userID string) (User, error)
}
