package impl

import (
	"context"
	"encoding/json"
	"errors"
	"office-booking-backend/pkg/database/redis"
	"office-booking-backend/pkg/entity"
	"office-booking-backend/pkg/utils/ptr"
	"testing"
	"time"

	redis2 "github.com/go-redis/redis/v9"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestSuiteTokenService struct {
	suite.Suite
	mockRedis *redis.RedisClientMock
	service   *TokenServiceImpl
}

func (s *TestSuiteTokenService) SetupTest() {
	s.mockRedis = &redis.RedisClientMock{}
	s.service = &TokenServiceImpl{
		AccessTokenSecret:  "123",
		AccessTokenExp:     15 * time.Minute,
		RefreshTokenSecret: "123",
		RefreshTokenExp:    15 * time.Minute,
		RedisRepository:    s.mockRedis,
	}
}

func (s *TestSuiteTokenService) TearDownTest() {
	s.mockRedis = nil
	s.service = nil
}

func TestTokenService(t *testing.T) {
	suite.Run(t, new(TestSuiteTokenService))
}

func (s *TestSuiteTokenService) TestNewTokenServiceImpl() {
	s.Run("NewTokenServiceImpl", func() {
		t := NewTokenServiceImpl("123", "123", 15*time.Minute, 15*time.Minute, s.mockRedis)
		s.NotNil(t)
	})
}

func (s *TestSuiteTokenService) TestGenerateAccessToken() {
	s.Run("generateAccessToken", func() {
		accessToken, tokenID, err := s.service.generateAccessToken(&entity.User{
			ID:         "123",
			Role:       1,
			IsVerified: ptr.Bool(true),
		}, time.Now().Add(s.service.AccessTokenExp))
		s.Nil(err)
		s.NotEmpty(accessToken)
		s.NotEmpty(tokenID)
	})
}

func (s *TestSuiteTokenService) TestGenerateRefreshToken() {
	s.Run("generateRefreshToken", func() {
		refreshToken, tokenID, err := s.service.generateRefreshToken(&entity.User{}, time.Now().Add(s.service.RefreshTokenExp))
		s.Nil(err)
		s.NotEmpty(refreshToken)
		s.NotEmpty(tokenID)
	})
}

func (s *TestSuiteTokenService) TestNewTokenPair() {
	user := &entity.User{
		ID:         "123",
		Role:       1,
		IsVerified: ptr.Bool(true),
	}
	s.Run("Success", func() {
		s.mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		tokenPair, err := s.service.NewTokenPair(context.Background(), user)
		s.NoError(err)
		s.NotNil(tokenPair)
	})
	s.TearDownTest()
	s.SetupTest()
	s.Run("Fail: redis error", func() {
		s.mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("redis error"))

		tokenPair, err := s.service.NewTokenPair(context.Background(), user)
		s.Nil(tokenPair)
		s.Error(err)
	})
}

func (s *TestSuiteTokenService) TestDeleteTokenPair() {
	s.Run("Success", func() {
		s.mockRedis.On("Del", mock.Anything, mock.Anything).Return(nil)

		err := s.service.DeleteTokenPair(context.Background(), "123")
		s.NoError(err)
	})
	s.TearDownTest()
	s.SetupTest()
	s.Run("Fail: redis error", func() {
		s.mockRedis.On("Del", mock.Anything, mock.Anything).Return(errors.New("redis error"))

		err := s.service.DeleteTokenPair(context.Background(), "123")
		s.Error(err)
	})
}

func (s *TestSuiteTokenService) TestCheckToken() {
	jsonToken, _ := json.Marshal(entity.CachedToken{
		AccessID:  "123",
		RefreshID: "123",
	})

	s.Run("Success: access token", func() {
		s.mockRedis.On("Get", mock.Anything, mock.Anything).Return(string(jsonToken), nil)

		_, err := s.service.CheckToken(context.Background(), &jwt.MapClaims{
			"jti": "123",
			"uid": "some-uid",
			"cat": "access",
		})
		s.NoError(err)
	})
	s.TearDownTest()
	s.SetupTest()
	s.Run("Success: refresh token", func() {
		s.mockRedis.On("Get", mock.Anything, mock.Anything).Return(string(jsonToken), nil)

		_, err := s.service.CheckToken(context.Background(), &jwt.MapClaims{
			"jti": "123",
			"uid": "some-uid",
			"cat": "refresh",
		})
		s.NoError(err)
	})
	s.TearDownTest()
	s.SetupTest()
	s.Run("Fail: redis error", func() {
		s.mockRedis.On("Get", mock.Anything, mock.Anything).Return("", errors.New("redis error"))

		_, err := s.service.CheckToken(context.Background(), &jwt.MapClaims{
			"jti": "123",
			"uid": "some-uid",
			"cat": "access",
		})
		s.Error(err)
	})
	s.TearDownTest()
	s.SetupTest()
	s.Run("Fail: token not found", func() {
		s.mockRedis.On("Get", mock.Anything, mock.Anything).Return("", redis2.Nil)

		_, err := s.service.CheckToken(context.Background(), &jwt.MapClaims{
			"jti": "123",
			"uid": "some-uid",
			"cat": "access",
		})
		s.NoError(err)
	})
	s.TearDownTest()
	s.SetupTest()
	s.Run("Fail: token id does not match", func() {
		s.mockRedis.On("Get", mock.Anything, mock.Anything).Return(string(jsonToken), nil)

		_, err := s.service.CheckToken(context.Background(), &jwt.MapClaims{
			"jti": "567",
			"uid": "some-uid",
			"cat": "access",
		})
		s.NoError(err)
	})
}

func (s *TestSuiteTokenService) TestParseRefreshToken() {
	dummyToken, _, _ := s.service.generateRefreshToken(&entity.User{}, time.Now().Add(s.service.RefreshTokenExp))
	s.Run("Success", func() {

		refreshToken, err := s.service.ParseRefreshToken(dummyToken)
		s.NoError(err)
		s.NotNil(refreshToken)
	})
	s.TearDownTest()
	s.SetupTest()
	s.Run("Fail: invalid token", func() {
		refreshToken, err := s.service.ParseRefreshToken("invalid-token")
		s.Error(err)
		s.Nil(refreshToken)
	})
}
