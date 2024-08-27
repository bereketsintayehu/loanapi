package usecase

import (
	"loan/domain"
	"loan/infrastracture"
	"log"
	"time"
)

type UserUsecase struct {
	UserRepo    domain.UserRepository
	TokenGen    domain.TokenGenerator
	PasswordSvc domain.PasswordService
}

func NewUserUsecase(userRepo domain.UserRepository, tokenGen domain.TokenGenerator, passwordSvc domain.PasswordService) domain.UserUsecase {
	return &UserUsecase{
		UserRepo:    userRepo,
		TokenGen:    tokenGen,
		PasswordSvc: passwordSvc,
	}
}

func (u *UserUsecase) Login(user *domain.User, deviceID string) (domain.LogInResponse, *domain.CustomError) {
	if u.UserRepo == nil || u.PasswordSvc == nil || u.TokenGen == nil {
		log.Fatal("Necessary services are nil")
		return domain.LogInResponse{}, domain.ErrInternalServer
	}

	existingUser, err := u.UserRepo.Login(user)
	if err != nil {
		return domain.LogInResponse{}, domain.ErrInvalidCredentials
	}

	if !u.PasswordSvc.CheckPasswordHash(user.Password, existingUser.Password) {
		return domain.LogInResponse{}, domain.ErrInvalidCredentials
	}

	if !existingUser.IsActive {
		return domain.LogInResponse{}, domain.ErrAccountNotActivated
	}

	refreshToken, err := u.TokenGen.GenerateRefreshToken(*existingUser)
	if err != nil {
		return domain.LogInResponse{}, domain.ErrInternalServer
	}

	newRefreshToken := domain.RefreshToken{
		Token:     refreshToken,
		DeviceID:  deviceID,
		CreatedAt: time.Now(),
	}

	for i, rt := range existingUser.RefreshTokens {
		if rt.DeviceID == deviceID {
			existingUser.RefreshTokens = append(existingUser.RefreshTokens[:i], existingUser.RefreshTokens[i+1:]...)
			break
		}
	}

	existingUser.RefreshTokens = append(existingUser.RefreshTokens, newRefreshToken)

	err = u.UserRepo.UpdateUser(existingUser)
	if err != nil {
		return domain.LogInResponse{}, domain.ErrFailedToUpdateUser
	}

	accessToken, err := u.TokenGen.GenerateToken(*existingUser)
	if err != nil {
		return domain.LogInResponse{}, domain.ErrInternalServer
	}

	return domain.LogInResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken.Token,
	}, &domain.CustomError{}
}

func (u *UserUsecase) ActivateAccount(token, email string) *domain.CustomError {
	if !infrastracture.IsValidEmail(email) {
		return domain.ErrInvalidEmail
	}

	err := u.UserRepo.ActivateAccount(token, email)
	if err != nil {
		return domain.ErrInvalidToken
	}

	return &domain.CustomError{}
}

func (u *UserUsecase) RefreshToken(userID, deviceID, token string) (domain.LogInResponse, *domain.CustomError) {
	user, err := u.UserRepo.GetUserByID(userID)
	if err != nil {
		return domain.LogInResponse{}, domain.ErrNotFound
	}

	for i, rt := range user.RefreshTokens {
		if rt.Token == token && rt.DeviceID == deviceID {
			user.RefreshTokens = append(user.RefreshTokens[:i], user.RefreshTokens[i+1:]...)

			refreshToken, err := u.TokenGen.GenerateRefreshToken(user)
			if err != nil {
				return domain.LogInResponse{}, domain.ErrInternalServer
			}

			newRefreshToken := domain.RefreshToken{
				Token:     refreshToken,
				DeviceID:  deviceID,
				CreatedAt: time.Now(),
			}

			user.RefreshTokens = append(user.RefreshTokens, newRefreshToken)
			err = u.UserRepo.UpdateUser(&user)
			if err != nil {
				return domain.LogInResponse{}, domain.ErrFailedToUpdateUser
			}

			accessToken, err := u.TokenGen.GenerateToken(user)
			if err != nil {
				return domain.LogInResponse{}, domain.ErrInternalServer
			}

			return domain.LogInResponse{
				AccessToken:  accessToken,
				RefreshToken: newRefreshToken.Token,
			}, &domain.CustomError{}
		}
	}

	return domain.LogInResponse{}, domain.ErrInvalidToken
}

func (u *UserUsecase) Register(user domain.User) *domain.CustomError {
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return domain.ErrMissingRequiredFields
	}

	if !infrastracture.IsValidEmail(user.Email) {
		return domain.ErrInvalidEmail
	}

	if !infrastracture.IsValidPassword(user.Password) {
		return domain.ErrInvalidPassword
	}

	_, err := u.UserRepo.GetUserByEmail(user.Email)
	if err == nil {
		return domain.ErrEmailAlreadyUsed
	}

	_, err = u.UserRepo.GetUserByUsername(user.Username)
	if err == nil {
		return domain.ErrUsernameAlreadyUsed
	}

	user.Role = "user"

	// Hash password
	hashedPassword, err := u.PasswordSvc.HashPassword(user.Password)
	if err != nil {
		return domain.ErrInternalServer
	}

	token, err := infrastracture.GenerateActivationToken()
	if err != nil {
		return domain.ErrInternalServer
	}

	user.Password = hashedPassword
	user.ActivationToken = token
	user.TokenCreatedAt = time.Now()

	// Create user account in the database
	err = u.UserRepo.Register(user)
	if err != nil {
		return domain.ErrInternalServer
	}

	// Send activation email or link to the user
	err = infrastracture.SendActivationEmail(user.Email, token)
	if err != nil {
		return domain.ErrFailedToSendEmail
	}

	return &domain.CustomError{}
}

func (u *UserUsecase) GetUserByUsernameOrEmail(username, email string) (domain.User, *domain.CustomError) {
	user, err := u.UserRepo.GetUserByUsernameOrEmail(username, email)
	if err != nil {
		return domain.User{}, domain.ErrNotFound
	}
	return user, &domain.CustomError{}
}

func (u *UserUsecase) PasswordReset(email string) *domain.CustomError {
	user, err := u.UserRepo.GetUserByEmail(email)
	if err != nil {
		return domain.ErrNotFound
	}

	resetToken, err := infrastracture.GenerateActivationToken()
	if err != nil {
		return domain.ErrInternalServer
	}
	user.PasswordResetToken = resetToken
	user.TokenCreatedAt = time.Now()

	err = u.UserRepo.UpdateUser(&user)
	if err != nil {
		return domain.ErrFailedToUpdateUser
	}

	err = infrastracture.SendResetLink(user.Email, resetToken)
	if err != nil {
		return domain.ErrFailedToSendEmail
	}

	return &domain.CustomError{}
}

func (uc *UserUsecase) GetMyProfile(userID string) (domain.User, *domain.CustomError) {
	user, err := uc.UserRepo.GetMyProfile(userID)
	if err != nil {
		return domain.User{}, domain.ErrNotFound
	}
	return user, &domain.CustomError{}
}

func (u *UserUsecase) UpdatePassword(req domain.PasswordUpdateRequest) *domain.CustomError {
	if req.Email == "" || req.Password == "" || req.Token == "" {
		return domain.ErrMissingRequiredFields
	}
	if !infrastracture.IsValidEmail(req.Email) {
		return domain.ErrInvalidEmail
	}
	if !infrastracture.IsValidPassword(req.Password) {
		return domain.ErrInvalidPassword
	}
	err := u.UserRepo.UpdatePassword(req.Token, req.Email, req.Password)
	if err != nil {
		return domain.ErrInvalidToken
	}

	return &domain.CustomError{}
}

func (uc *UserUsecase) GetNewVerificationEmail(email string) *domain.CustomError {
	user, err := uc.UserRepo.GetUserByEmail(email)

	token, err := infrastracture.GenerateActivationToken()
	if err != nil {
		return domain.ErrInternalServer
	}

	user.ActivationToken = token
	user.TokenCreatedAt = time.Now()

	if user.IsActive {
		return &domain.CustomError{}
	}

	err = uc.UserRepo.UpdateUser(&user)
	if err != nil {
		return domain.ErrInternalServer
	}

	err = infrastracture.SendActivationEmail(user.Email, token)
	if err != nil {
		return domain.ErrFailedToSendEmail
	}

	return &domain.CustomError{}
}
