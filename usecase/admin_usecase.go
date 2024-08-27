package usecase

import (
	"loan/domain"
)

func (uc *UserUsecase) GetUsers() ([]domain.User, *domain.CustomError) {
	users, err := uc.UserRepo.GetUsers()
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return users, &domain.CustomError{}
}

func (uc *UserUsecase) DeleteUser(userID string) (domain.User, *domain.CustomError) {
	user, err := uc.UserRepo.DeleteUser(userID)
	if err != nil {
		return domain.User{}, domain.ErrFailedToDeleteUser
	}
	return user, &domain.CustomError{}
}
