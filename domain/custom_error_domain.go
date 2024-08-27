package domain

import "net/http"

type CustomError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func New(message string, statusCode int) *CustomError {
	return &CustomError{
		Message:    message,
		StatusCode: statusCode,
	}
}

func (e *CustomError) Error() string {
	return e.Message
}

var (
	// General Errors
	ErrNotFound       = New("resource not found", http.StatusNotFound)
	ErrInternalServer = New("internal server error", http.StatusInternalServerError)
	ErrBadRequest     = New("bad request", http.StatusBadRequest)
	ErrUnauthorized   = New("unauthorized", http.StatusUnauthorized)
	ErrForbidden      = New("forbidden", http.StatusForbidden)

	// Auth-specific Errors
	ErrUserNotFound         = New("user not found", http.StatusNotFound)
	ErrRefreshTokenNotFound = New("refresh token not found", http.StatusNotFound)
	ErrInvalidRefreshToken  = New("invalid refresh token", http.StatusUnauthorized)
	ErrInvalidAccessToken   = New("invalid access token", http.StatusUnauthorized)
	ErrExpiredAccessToken   = New("expired access token", http.StatusUnauthorized)
	ErrExpiredRefreshToken  = New("expired refresh token", http.StatusUnauthorized)
	ErrInvalidToken         = New("invalid token", http.StatusUnauthorized)
	ErrInvalidRole          = New("invalid role", http.StatusUnauthorized)
	ErrInvalidUserID        = New("invalid user id", http.StatusBadRequest)
	ErrInvalidDeviceID      = New("invalid device id", http.StatusBadRequest)
	ErrDeviceNotFound       = New("device not found", http.StatusNotFound)
	ErrInvalidEmail         = New("invalid email", http.StatusBadRequest)
	ErrInvalidPassword      = New("invalid password", http.StatusBadRequest)

	ErrFailedToUpdateUser    = New("failed to update user", http.StatusInternalServerError)
	ErrMissingRequiredFields = New("missing required fields", http.StatusBadRequest)
	ErrInvalidUpdateRequest  = New("invalid update request", http.StatusBadRequest)
	ErrFailedToSendEmail     = New("failed to send email", http.StatusInternalServerError)
	ErrActivationFailed      = New("account activation failed", http.StatusInternalServerError)
	ErrFailedToDeleteUser    = New("failed to delete user", http.StatusInternalServerError)
	ErrFailedToDeleteAccount = New("failed to delete account", http.StatusInternalServerError)
	ErrFailedToUploadImage   = New("failed to upload image", http.StatusInternalServerError)
	ErrFailedToUpdateProfile = New("failed to update profile", http.StatusInternalServerError)
	ErrAccountNotActivated  = New("account not activated", http.StatusUnauthorized)
	ErrInvalidCredentials  = New("invalid credentials", http.StatusUnauthorized)
	ErrEmailAlreadyUsed = New("email already used", http.StatusBadRequest)
	ErrUsernameAlreadyUsed = New("username already used", http.StatusBadRequest)
)
