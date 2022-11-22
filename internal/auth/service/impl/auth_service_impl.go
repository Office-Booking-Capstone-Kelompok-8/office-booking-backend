package impl

import (
	"context"
	"log"
	"office-booking-backend/internal/auth/dto"
	"office-booking-backend/internal/auth/repository"
	"office-booking-backend/internal/auth/service"
	"office-booking-backend/pkg/utils/password"
)

const DefaultPasswordCost = 10

type AuthServiceImpl struct {
	repository repository.AuthRepository
	password   password.PasswordService
}

func NewAuthServiceImpl(repository repository.AuthRepository, password password.PasswordService) service.AuthService {
	return &AuthServiceImpl{
		repository: repository,
		password:   password,
	}
}

func (a *AuthServiceImpl) RegisterUser(ctx context.Context, user *dto.SignupRequest) error {
	hashedPassword, err := a.password.GenerateFromPassword([]byte(user.Password), DefaultPasswordCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	userEntity := user.ToEntity()

	err = a.repository.RegisterUser(ctx, userEntity)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
