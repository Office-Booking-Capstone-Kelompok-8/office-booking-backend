package impl

import (
	"context"
	"log"
	"office-booking-backend/internal/user/dto"
	"office-booking-backend/internal/user/repository"
	"office-booking-backend/internal/user/service"
)

type UserServiceImpl struct {
	userRepository repository.UserRepository
}

func NewUserServiceImpl(userRepository repository.UserRepository) service.UserService {
	return &UserServiceImpl{
		userRepository: userRepository,
	}
}

func (u *UserServiceImpl) GetFullUserByID(ctx context.Context, id string) (*dto.UserResponse, error) {
	user, err := u.userRepository.GetFullUserByID(ctx, id)
	if err != nil {
		log.Println("Error while getting user by id: ", err)
		return nil, err
	}

	fullUser := dto.NewUserResponse(user)
	return fullUser, nil
}
