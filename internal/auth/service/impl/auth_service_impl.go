package impl

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"office-booking-backend/internal/auth/dto"
	"office-booking-backend/internal/auth/repository"
	"office-booking-backend/internal/auth/service"
	"office-booking-backend/pkg/config"
	"office-booking-backend/pkg/database/redis"
	"office-booking-backend/pkg/entity"
	err2 "office-booking-backend/pkg/errors"
	"office-booking-backend/pkg/utils/mail"
	"office-booking-backend/pkg/utils/password"
	"office-booking-backend/pkg/utils/random"

	redis2 "github.com/go-redis/redis/v9"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

const DefaultPasswordCost = 10

type AuthServiceImpl struct {
	repository repository.AuthRepository
	token      service.TokenService
	redisRepo  redis.RedisClient
	mail       mail.Client
	password   password.Hash
	generator  random.Generator
	config     *viper.Viper
}

func NewAuthServiceImpl(repository repository.AuthRepository, tokenService service.TokenService, redisRepo redis.RedisClient, mail mail.Client, password password.Hash, generator random.Generator, config *viper.Viper) service.AuthService {
	return &AuthServiceImpl{
		repository: repository,
		token:      tokenService,
		redisRepo:  redisRepo,
		mail:       mail,
		password:   password,
		generator:  generator,
		config:     config,
	}
}

func (a *AuthServiceImpl) RegisterUser(ctx context.Context, user *dto.SignupRequest) (string, error) {
	hashedPassword, err := a.password.GenerateFromPassword([]byte(user.Password), DefaultPasswordCost)
	if err != nil {
		return "", err
	}

	user.Password = string(hashedPassword)
	userEntity := user.ToEntity()

	err = a.repository.RegisterUser(ctx, userEntity)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return userEntity.ID, nil
}

func (a *AuthServiceImpl) RegisterAdmin(ctx context.Context, user *dto.SignupRequest) (string, error) {
	hashedPassword, err := a.password.GenerateFromPassword([]byte(user.Password), DefaultPasswordCost)
	if err != nil {
		return "", err
	}

	user.Password = string(hashedPassword)
	userEntity := user.ToEntity()
	userEntity.Role = config.ADMIN_ROLE

	err = a.repository.RegisterUser(ctx, userEntity)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return userEntity.ID, nil
}

func (a *AuthServiceImpl) LoginUser(ctx context.Context, user *dto.LoginRequest) (*dto.TokenPair, error) {
	userEntity, err := a.repository.GetUserByEmail(ctx, user.Email)
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

	jwtMapClaim := jwt.MapClaims{
		"uid": claims.UID,
		"cat": claims.Category,
		"jti": claims.ID,
	}

	ok, err := a.token.CheckToken(ctx, &jwtMapClaim)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, err2.ErrInvalidToken
	}

	user, err := a.repository.GetUserByID(ctx, claims.UID)
	if err != nil {
		log.Println("Error while finding user by id: ", err)
		return nil, err
	}

	return a.token.NewTokenPair(ctx, user)
}

func createKey(email string, subject string) string {
	hasher := md5.New()
	hasher.Write([]byte(fmt.Sprintf("%s:%s", email, subject)))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (a *AuthServiceImpl) createResetOTP(ctx context.Context, email string) (*string, error) {
	user, err := a.repository.GetUserByEmail(ctx, email)
	if err != nil {
		log.Println("Error while finding user by email: ", err)
		return nil, err
	}

	otp, err := a.generator.GenerateRandomIntString(a.config.GetInt("otp.length"))
	if err != nil {
		log.Println("Error while generating random string: ", err)
		return nil, err
	}

	otpTokenPair, err := json.Marshal(entity.CachedOTP{
		OTP:    otp,
		Key:    uuid.New().String(),
		UserID: user.ID,
	})
	if err != nil {
		log.Println("Error while marshalling otp token pair: ", err)
		return nil, err
	}
	key := createKey(email, config.RESET_PASSWORD_SUBJECT)

	err = a.redisRepo.Set(ctx, key, string(otpTokenPair), a.config.GetDuration("otp.exp"))
	if err != nil {
		log.Println("Error while setting key in redis: ", err)
		return nil, err
	}

	return &otp, nil
}

func (a *AuthServiceImpl) RequestPasswordResetOTP(ctx context.Context, email string) error {
	otp, err := a.createResetOTP(ctx, email)
	if err != nil {
		log.Println("Error while creating otp: ", err)
		return err
	}

	msg := &mail.Mail{
		Subject:  "Reset Password OTP",
		Template: "password-reset",
		Variable: map[string]string{
			"otp": *otp,
		},
		Recipient: email,
	}

	err = a.mail.SendMail(ctx, msg)
	if err != nil {
		log.Println("Error while sending email: ", err)
		return err
	}

	return nil
}

func (a *AuthServiceImpl) VerifyPasswordResetOTP(ctx context.Context, otp *dto.ResetPasswordOTPVerifyRequest) (*string, error) {
	key := createKey(otp.Email, config.RESET_PASSWORD_SUBJECT)
	jsonToken, err := a.redisRepo.Get(ctx, key)
	if err != nil {
		if errors.Is(err, redis2.Nil) {
			return nil, err2.ErrInvalidOTP
		}

		log.Println("Error while getting key from redis: ", err)
		return nil, err
	}

	token := new(entity.CachedOTP)
	err = json.Unmarshal([]byte(jsonToken), token)
	if err != nil {
		log.Println("Error while unmarshalling json token: ", err)
		return nil, err
	}

	if token.OTP != otp.Code {
		return nil, err2.ErrInvalidOTP
	}

	return &token.Key, nil
}

func (a *AuthServiceImpl) createVerifyOTP(ctx context.Context, userID string, email string) (*string, error) {
	otp, err := a.generator.GenerateRandomIntString(a.config.GetInt("otp.length"))
	if err != nil {
		log.Println("Error while generating random string: ", err)
		return nil, err
	}

	otpTokenPair, err := json.Marshal(entity.CachedOTP{
		OTP:    otp,
		UserID: userID,
	})
	if err != nil {
		log.Println("Error while marshalling otp token pair: ", err)
		return nil, err
	}
	key := createKey(email, config.VERIFY_EMAIL_SUBJECT)

	err = a.redisRepo.Set(ctx, key, string(otpTokenPair), a.config.GetDuration("otp.exp"))
	if err != nil {
		log.Println("Error while setting key in redis: ", err)
		return nil, err
	}

	return &otp, nil
}

func (a *AuthServiceImpl) RequestEmailOTP(ctx context.Context, userID string) error {
	user, err := a.repository.GetUserByID(ctx, userID)
	if err != nil {
		log.Println("Error while finding user by id: ", err)
		return err
	}

	if *user.IsVerified {
		return err2.ErrEmailAlreadyVerified
	}

	otp, err := a.createVerifyOTP(ctx, userID, user.Email)
	if err != nil {
		log.Println("Error while creating otp: ", err)
		return err
	}

	msg := &mail.Mail{
		Subject:  "Verify Email OTP",
		Template: "email-verify",
		Variable: map[string]string{
			"otp": *otp,
		},
		Recipient: user.Email,
	}

	err = a.mail.SendMail(ctx, msg)
	if err != nil {
		log.Println("Error while sending email: ", err)
		return err
	}

	return nil
}
func (a *AuthServiceImpl) VerifyEmailOTP(ctx context.Context, userID string, otp *dto.VerifyEmailOTOPVerifyRequest) error {
	user, err := a.repository.GetUserByID(ctx, userID)
	if err != nil {
		log.Println("Error while finding user by id: ", err)
		return err
	}

	key := createKey(user.Email, config.VERIFY_EMAIL_SUBJECT)
	jsonToken, err := a.redisRepo.Get(ctx, key)
	if err != nil {
		if errors.Is(err, redis2.Nil) {
			return err2.ErrInvalidOTP
		}

		log.Println("Error while getting key from redis: ", err)
		return err
	}

	token := new(entity.CachedOTP)
	err = json.Unmarshal([]byte(jsonToken), token)
	if err != nil {
		log.Println("Error while unmarshalling json token: ", err)
		return err
	}

	if token.OTP != otp.Code {
		return err2.ErrInvalidOTP
	}

	err = a.repository.VerifyEmail(ctx, token.UserID)
	if err != nil {
		log.Println("Error while verifying email: ", err)
		return err
	}

	return nil
}

func (a *AuthServiceImpl) ResetPassword(ctx context.Context, password *dto.PasswordResetRequest) error {
	key := createKey(password.Email, config.RESET_PASSWORD_SUBJECT)
	jsonOTP, err := a.redisRepo.Get(ctx, key)
	if err != nil {
		if errors.Is(err, redis2.Nil) {
			return err2.ErrInvalidOTPToken
		}

		log.Println("Error while getting key from redis: ", err)
		return err
	}

	otp := new(entity.CachedOTP)
	err = json.Unmarshal([]byte(jsonOTP), otp)
	if err != nil {
		log.Println("Error while unmarshalling json token: ", err)
		return err
	}

	if otp.Key != password.Key {
		return err2.ErrInvalidOTPToken
	}

	err = a.redisRepo.Del(ctx, key)
	if err != nil {
		log.Println("Error while deleting key from redis: ", err)
		return err
	}

	hashedPassword, err := a.password.GenerateFromPassword([]byte(password.Password), DefaultPasswordCost)
	if err != nil {
		return err
	}

	err = a.repository.ChangePassword(ctx, otp.UserID, string(hashedPassword))
	if err != nil {
		log.Println("Error while changing password: ", err)
		return err
	}

	return nil
}

func (a *AuthServiceImpl) ChangePassword(ctx context.Context, uid string, password *dto.ChangePasswordRequest) error {
	user, err := a.repository.GetUserByID(ctx, uid)
	if err != nil {
		log.Println("Error while finding user by id: ", err)
		return err
	}

	err = a.password.CompareHashAndPassword([]byte(user.Password), []byte(password.OldPassword))
	if err != nil {
		return err2.ErrPasswordNotMatch
	}

	hashedPassword, err := a.password.GenerateFromPassword([]byte(password.NewPassword), DefaultPasswordCost)
	if err != nil {
		return err
	}

	err = a.repository.ChangePassword(ctx, uid, string(hashedPassword))
	if err != nil {
		log.Println("Error while changing password: ", err)
		return err
	}

	return nil
}
