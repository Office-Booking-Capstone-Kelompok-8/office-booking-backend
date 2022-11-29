package impl

import (
	"context"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"io"
	"log"
	service2 "office-booking-backend/internal/reservation/service"
	"office-booking-backend/internal/user/dto"
	"office-booking-backend/internal/user/repository"
	"office-booking-backend/internal/user/service"
	"office-booking-backend/pkg/entity"
	err2 "office-booking-backend/pkg/errors"
	"office-booking-backend/pkg/utils/imagekit"
)

type UserServiceImpl struct {
	userRepository     repository.UserRepository
	reservationService service2.ReservationService
	imgKitService      imagekit.ImgKitService
}

func NewUserServiceImpl(userRepository repository.UserRepository, resevationService service2.ReservationService, imgKitService imagekit.ImgKitService) service.UserService {
	return &UserServiceImpl{
		userRepository:     userRepository,
		reservationService: resevationService,
		imgKitService:      imgKitService,
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

func (u *UserServiceImpl) GetAllUsers(ctx context.Context, q string, limit int, page int) (*dto.BriefUsersResponse, int64, error) {
	offset := (page - 1) * limit
	users, total, err := u.userRepository.GetAllUsers(ctx, q, limit, offset)
	if err != nil {
		log.Println("Error while getting users: ", err)
		return nil, 0, err
	}

	briefUsers := dto.NewBriefUsersResponse(users)
	return briefUsers, total, nil
}

func (u *UserServiceImpl) UpdateUserByID(ctx context.Context, id string, user *dto.UserUpdateRequest) error {
	userEntity := user.ToEntity()
	userEntity.ID = id
	userEntity.Detail.UserID = id

	if userEntity.Email != "" {
		userEntity.IsVerified = false
	}

	errGroup, c := errgroup.WithContext(ctx)
	errGroup.Go(func() error {
		err := u.userRepository.UpdateUserByID(c, userEntity)
		if err != nil {
			log.Println("Error while updating user by id: ", err)
			return err
		}
		return nil
	})

	errGroup.Go(func() error {
		err := u.userRepository.UpdateUserDetailByID(c, &userEntity.Detail)
		if err != nil {
			log.Println("Error while updating user detail by id: ", err)
			return err
		}
		return nil
	})

	return errGroup.Wait()
}

func (u *UserServiceImpl) DeleteUserByID(ctx context.Context, id string) error {
	userReservation, err := u.reservationService.CountUserActiveReservations(ctx, id)
	if err != nil {
		return err
	}

	if userReservation > 0 {
		return err2.ErrUserHasReservation
	}

	err = u.userRepository.DeleteUserByID(ctx, id)
	if err != nil {
		log.Println("Error while deleting user by id: ", err)
		return err
	}

	return nil
}

func (u *UserServiceImpl) UploadUserAvatar(ctx context.Context, id string, file io.Reader) error {
	fileName := uuid.New().String()

	uploadResult, err := u.imgKitService.UploadFile(ctx, file, fileName, "avatars")
	if err != nil {
		log.Println("Error while uploading user avatar: ", err)
		return err2.ErrPictureServiceFailed
	}

	user := new(entity.UserDetail)
	user.Picture = *dto.NewPictureEntity(uploadResult)
	user.UserID = id

	err = u.userRepository.UpdateUserDetailByID(ctx, user)
	if err != nil {
		log.Println("Error while updating user detail by id: ", err)
		return err
	}

	return nil
}
