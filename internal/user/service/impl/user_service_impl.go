package impl

import (
	"context"
	"io"
	"log"
	service2 "office-booking-backend/internal/reservation/service"
	"office-booking-backend/internal/user/dto"
	"office-booking-backend/internal/user/repository"
	"office-booking-backend/internal/user/service"
	"office-booking-backend/pkg/custom"
	"office-booking-backend/pkg/entity"
	err2 "office-booking-backend/pkg/errors"
	"office-booking-backend/pkg/utils/imagekit"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type UserServiceImpl struct {
	userRepository     repository.UserRepository
	reservationService service2.ReservationService
	imgKitService      imagekit.ImgKitService
}

func NewUserServiceImpl(userRepository repository.UserRepository, reservationService service2.ReservationService, imgKitService imagekit.ImgKitService) service.UserService {
	return &UserServiceImpl{
		userRepository:     userRepository,
		reservationService: reservationService,
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

func (u *UserServiceImpl) GetAllUsers(ctx context.Context, filter *dto.UserFilterRequest) (*dto.BriefUsersResponse, int64, error) {
	filter.Offset = (filter.Page - 1) * filter.Limit
	users, total, err := u.userRepository.GetAllUsers(ctx, filter)
	if err != nil {
		log.Println("Error while getting users: ", err)
		return nil, 0, err
	}

	briefUsers := dto.NewBriefUsersResponse(users)
	return briefUsers, total, nil
}

func (u *UserServiceImpl) GetRegisteredMemberStat(ctx context.Context) (*dto.RegisteredStatResponseList, error) {
	stat, err := u.userRepository.GetRegisteredMemberStat(ctx)
	if err != nil {
		log.Println("Error while getting registered member stat: ", err)
		return nil, err
	}

	return dto.NewRegisteredStatResponseList(stat), nil
}

func (u *UserServiceImpl) UpdateUserByID(ctx context.Context, id string, user *dto.UserUpdateRequest) error {
	userEntity := user.ToEntity()
	userEntity.ID = id
	userEntity.Detail.UserID = id

	if userEntity.Email != "" {
		userEntity.IsVerified = custom.Bool(false)
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

	picId, err := u.userRepository.DeleteUserByID(ctx, id)
	if err != nil {
		log.Println("Error while deleting user by id: ", err)
		return err
	}

	if picId != "" {
		err = u.imgKitService.DeleteFile(ctx, picId)
		if err != nil {
			log.Println("Error while deleting user avatar: ", err)
			return err
		}
	}

	return nil
}

func (u *UserServiceImpl) UploadUserAvatar(ctx context.Context, id string, file io.Reader) error {
	picture, err := u.userRepository.GetUserProfilePictureID(ctx, id)
	if err != nil {
		log.Println("Error while getting user profile picture id: ", err)
		return err
	}

	errGroup := errgroup.Group{}
	errGroup.Go(func() error {
		if picture.ID == "" {
			return nil
		}

		err = u.userRepository.DeleteUserProfilePicture(ctx, picture.ID)
		if err != nil {
			log.Println("Error while deleting user profile picture: ", err)
			return err
		}

		err = u.imgKitService.DeleteFile(ctx, picture.ID)
		if err != nil {
			log.Println("Error while deleting user avatar: ", err)
			return err
		}

		return nil
	})

	errGroup.Go(func() error {
		pictureKey := uuid.New().String()
		uploadResult, err := u.imgKitService.UploadFile(ctx, file, pictureKey, "avatars")
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
	})

	return errGroup.Wait()
}
