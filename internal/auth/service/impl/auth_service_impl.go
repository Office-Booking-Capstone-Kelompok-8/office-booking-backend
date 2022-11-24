package impl

import (
	"context"
	"errors"
	"log"
	"office-booking-backend/internal/auth/dto"
	"office-booking-backend/internal/auth/repository"
	"office-booking-backend/internal/auth/service"
	err2 "office-booking-backend/pkg/errors"
	"office-booking-backend/pkg/utils/password"
)

const DefaultPasswordCost = 10

type AuthServiceImpl struct {
	repository repository.AuthRepository
	password   password.PasswordService
	token      service.TokenService
}

func NewAuthServiceImpl(repository repository.AuthRepository, tokenService service.TokenService, password password.PasswordService) service.AuthService {
	return &AuthServiceImpl{
		repository: repository,
		password:   password,
		token:      tokenService,
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

func (a *AuthServiceImpl) LoginUser(ctx context.Context, user *dto.LoginRequest) (*dto.TokenPair, error) {
	userEntity, err := a.repository.FindUserByEmail(ctx, user.Email)
	if err != nil {
		if errors.Is(err, err2.ErrUserNotFound) {
			return nil, err2.ErrInvalidCredentials
		}

		log.Println("Error while finding user by email: ", err)
		return nil, err
	}

	err = a.password.CompareHashAndPassword([]byte(userEntity.Password), []byte(user.Password))
	if err != nil {
		return nil, err2.ErrInvalidCredentials
	}

	token, err := a.token.NewTokenPair(ctx, userEntity)
	if err != nil {
		log.Println("Error while generating token pair: ", err)
		return nil, err
	}

	return token, nil
}

func (a *AuthServiceImpl) LogoutUser(ctx context.Context, uid string) error {
	return a.token.DeleteTokenPair(ctx, uid)
}

func (a *AuthServiceImpl) RefreshToken(ctx context.Context, token *dto.RefreshTokenRequest) (*dto.TokenPair, error) {
	claims, err := a.token.ParseRefreshToken(token.RefreshToken)
	if err != nil {
		return nil, err
	}

	user, err := a.repository.FindUserByID(ctx, claims.UID)
	if err != nil {
		log.Println("Error while finding user by id: ", err)
		return nil, err
	}

	return a.token.NewTokenPair(ctx, user)
}
