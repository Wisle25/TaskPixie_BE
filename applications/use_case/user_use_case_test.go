package use_case_test

import (
	"github.com/wisle25/be-template/applications/file_statics"
	"github.com/wisle25/be-template/applications/use_case"
	"github.com/wisle25/be-template/domains/entity"
	"io"
	"mime/multipart"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wisle25/be-template/commons"
)

// Mocks for the dependencies
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) AddUser(payload *entity.RegisterUserPayload) string {
	args := m.Called(payload)

	return args.String(0)
}

func (m *MockUserRepository) GetUserForLogin(identity string) (string, string) {
	args := m.Called(identity)

	return args.String(0), args.String(1)
}

func (m *MockUserRepository) GetUserById(id string) *entity.User {
	args := m.Called(id)

	return args.Get(0).(*entity.User)
}

func (m *MockUserRepository) UpdateUserById(id string, payload *entity.UpdateUserPayload, avatarLink string) string {
	args := m.Called(id, payload, avatarLink)

	return args.String(0)
}

type MockPasswordHash struct {
	mock.Mock
}

func (m *MockPasswordHash) Hash(password string) string {
	args := m.Called(password)
	return args.String(0)
}

func (m *MockPasswordHash) Compare(password string, hashedPassword string) {
	m.Called(password, hashedPassword)
}

type MockValidateUser struct {
	mock.Mock
}

func (m *MockValidateUser) ValidateRegisterPayload(payload *entity.RegisterUserPayload) {
	m.Called(payload)
}

func (m *MockValidateUser) ValidateLoginPayload(payload *entity.LoginUserPayload) {
	m.Called(payload)
}

func (m *MockValidateUser) ValidateUpdatePayload(payload *entity.UpdateUserPayload) {
	m.Called(payload)
}

type MockToken struct {
	mock.Mock
}

func (m *MockToken) CreateToken(userID string, ttl time.Duration, privateKey string) *entity.TokenDetail {
	args := m.Called(userID, ttl, privateKey)
	return args.Get(0).(*entity.TokenDetail)
}

func (m *MockToken) ValidateToken(token string, publicKey string) *entity.TokenDetail {
	args := m.Called(token, publicKey)
	return args.Get(0).(*entity.TokenDetail)
}

type MockCache struct {
	mock.Mock
}

func (m *MockCache) SetCache(key string, value interface{}, expiration time.Duration) {
	m.Called(key, value, expiration)
}

func (m *MockCache) GetCache(key string) interface{} {
	args := m.Called(key)
	return args.Get(0)
}

func (m *MockCache) DeleteCache(key string) {
	m.Called(key)
}

type MockFileUpload struct {
	mock.Mock
}

func (m *MockFileUpload) UploadFile(buffer []byte, extension string) string {
	args := m.Called(buffer, extension)

	return args.String(0)
}

func (m *MockFileUpload) RemoveFile(oldFileLink string) {
	m.Called(oldFileLink)
}

type MockFileProcessing struct {
	mock.Mock
}

func (m *MockFileProcessing) CompressImage(buffer []byte, to file_statics.ConvertTo) ([]byte, string) {
	args := m.Called(buffer, to)

	return args.Get(0).([]byte), args.String(1)
}

func (m *MockFileProcessing) ResizeImage(fileHeader *multipart.FileHeader) {
	//TODO implement me
	panic("implement me")
}

func TestUserUseCase(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockPasswordHash := new(MockPasswordHash)
	mockValidator := new(MockValidateUser)
	mockConfig := &commons.Config{
		AccessTokenExpiresIn:  time.Hour,
		RefreshTokenExpiresIn: time.Hour * 24,
		AccessTokenPrivateKey: "any",
		RefreshTokenPublicKey: "any",
	}
	mockToken := new(MockToken)
	mockCache := new(MockCache)
	mockFileUpload := new(MockFileUpload)
	mockFileProcessing := new(MockFileProcessing)

	userUseCase := use_case.NewUserUseCase(
		mockUserRepo,
		mockFileProcessing,
		mockFileUpload,
		mockPasswordHash,
		mockValidator,
		mockConfig,
		mockToken,
		mockCache,
	)

	t.Run("Execute Add", func(t *testing.T) {
		// Arrange
		payload := &entity.RegisterUserPayload{
			Username: "testuser",
			Password: "password123",
			Email:    "test@example.com",
		}

		mockValidator.On("ValidateRegisterPayload", payload).Return(nil)
		mockPasswordHash.On("Hash", payload.Password).Return("hashedpassword")
		mockUserRepo.On("AddUser", payload).Return("userid123")

		// Action
		userId := userUseCase.ExecuteAdd(payload)

		// Assert
		assert.Equal(t, "userid123", userId)

		mockValidator.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockPasswordHash.AssertExpectations(t)
	})

	t.Run("Execute Login", func(t *testing.T) {
		// Arrange
		payload := &entity.LoginUserPayload{
			Identity: "testuser",
			Password: "password123",
		}

		user := &entity.User{
			Id:       "userid123",
			Username: "testuser",
			Email:    "test@example.com",
		}

		accessTokenDetail := &entity.TokenDetail{
			TokenId:   "access_token_id",
			ExpiresIn: time.Now().Add(time.Hour).Unix(),
			UserId:    "userid123",
			Token:     "access_token",
		}

		refreshTokenDetail := &entity.TokenDetail{
			TokenId:   "refresh_token_id",
			ExpiresIn: time.Now().Add(time.Hour * 24).Unix(),
			UserId:    "userid123",
			Token:     "refresh_token",
		}

		mockValidator.On("ValidateLoginPayload", payload).Return(nil)
		mockUserRepo.On("GetUserForLogin", payload.Identity).Return(user.Id, "hashedpassword")
		mockPasswordHash.On("Compare", payload.Password, "hashedpassword").Return(nil)
		mockToken.On("CreateToken", user.Id, mockConfig.AccessTokenExpiresIn, mockConfig.AccessTokenPrivateKey).Return(accessTokenDetail)
		mockToken.On("CreateToken", user.Id, mockConfig.RefreshTokenExpiresIn, mockConfig.RefreshTokenPrivateKey).Return(refreshTokenDetail)
		mockCache.On("SetCache", accessTokenDetail.TokenId, user.Id, mock.Anything).Return(nil)
		mockCache.On("SetCache", refreshTokenDetail.TokenId, user.Id, mock.Anything).Return(nil)

		// Action
		accessToken, refreshToken := userUseCase.ExecuteLogin(payload)

		// Assert
		assert.Equal(t, accessTokenDetail, accessToken)
		assert.Equal(t, refreshTokenDetail, refreshToken)

		mockValidator.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockPasswordHash.AssertExpectations(t)
		mockToken.AssertExpectations(t)
		mockCache.AssertExpectations(t)
	})

	t.Run("Execute Refresh Token", func(t *testing.T) {
		// Arrange
		refreshTokenCookie := "refresh_token123"

		accessTokenDetail := &entity.TokenDetail{
			TokenId:   "access_token_id",
			ExpiresIn: time.Now().Add(time.Hour).Unix(),
			UserId:    "userid123",
			Token:     "access_token",
		}
		refreshTokenDetail := &entity.TokenDetail{
			TokenId:   "refresh_token_id",
			ExpiresIn: time.Now().Add(time.Hour * 24).Unix(),
			UserId:    "userid123",
			Token:     "refresh_token",
		}

		mockToken.On("ValidateToken", refreshTokenCookie, mockConfig.RefreshTokenPublicKey).Return(refreshTokenDetail)
		mockCache.On("GetCache", refreshTokenDetail.TokenId).Return(refreshTokenDetail.UserId)
		mockToken.On("CreateToken", refreshTokenDetail.UserId, mockConfig.AccessTokenExpiresIn, mockConfig.AccessTokenPrivateKey).Return(accessTokenDetail)
		mockCache.On("SetCache", accessTokenDetail.TokenId, refreshTokenDetail.UserId, mock.Anything).Return(nil)

		// Action
		accessTokenResponse := userUseCase.ExecuteRefreshToken(refreshTokenCookie)

		// Assert
		assert.Equal(t, accessTokenDetail, accessTokenResponse)
		mockToken.AssertExpectations(t)
		mockCache.AssertExpectations(t)
	})

	t.Run("Execute Logout", func(t *testing.T) {
		// Arrange
		refreshTokenCookie := "refresh_token123"
		accessTokenId := "access_token_id"

		refreshTokenDetail := &entity.TokenDetail{
			TokenId: "refresh_token_id",
		}

		mockToken.On("ValidateToken", refreshTokenCookie, mockConfig.RefreshTokenPublicKey).Return(refreshTokenDetail)
		mockCache.On("DeleteCache", refreshTokenDetail.TokenId).Return(nil).Once()
		mockCache.On("DeleteCache", accessTokenId).Return(nil).Once()

		// Action
		userUseCase.ExecuteLogout(refreshTokenCookie, accessTokenId)

		// Assert
		mockToken.AssertExpectations(t)
		mockCache.AssertExpectations(t)
	})

	t.Run("Execute Guard", func(t *testing.T) {
		// Arrange
		accessToken := "access_token123"
		accessTokenDetail := &entity.TokenDetail{
			TokenId: "access_token123",
			UserId:  "userid123",
		}

		mockToken.On("ValidateToken", accessToken, mockConfig.AccessTokenPublicKey).Return(accessTokenDetail)
		mockCache.On("GetCache", accessTokenDetail.TokenId).Return(accessTokenDetail.UserId)

		// Action
		userIdCache, tokenDetail := userUseCase.ExecuteGuard(accessToken)

		assert.NotNil(t, userIdCache)
		assert.Equal(t, accessTokenDetail, tokenDetail)
		mockToken.AssertExpectations(t)
		mockCache.AssertExpectations(t)
	})

	t.Run("Execute GetUserById", func(t *testing.T) {
		// Arrange
		userId := "userid123"
		expectedUser := &entity.User{
			Id:         userId,
			Username:   "testuser",
			Email:      "test@example.com",
			AvatarLink: "anything",
		}

		mockUserRepo.On("GetUserById", userId).Return(expectedUser)

		// Actions
		user := userUseCase.ExecuteGetUserById(userId)

		// Assert
		assert.Equal(t, expectedUser, user)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Execute UpdateUserById", func(t *testing.T) {
		// Arrange
		userId := "userid123"
		payload := &entity.UpdateUserPayload{
			Username:        "username",
			Email:           "email",
			Password:        "any password",
			ConfirmPassword: "any password",
			Avatar: &multipart.FileHeader{
				Filename: "Any",
				Header:   nil,
				Size:     0,
			},
		}
		avatarLink := "avatar_link"
		oldAvatarLink := "old_avatar_link"
		file, _ := payload.Avatar.Open()
		avatarBuffer, _ := io.ReadAll(file)
		var compressBuffer []byte

		mockValidator.On("ValidateUpdatePayload", payload).Return(nil)
		mockPasswordHash.On("Hash", payload.Password).Return("hashedPassword")
		mockFileProcessing.On("CompressImage", avatarBuffer, mock.Anything).Return(compressBuffer, ".webp")
		mockFileUpload.On("UploadFile", compressBuffer, ".webp").Return(avatarLink)
		mockUserRepo.On("UpdateUserById", userId, payload, avatarLink).Return(oldAvatarLink)
		mockFileUpload.On("RemoveFile", oldAvatarLink).Return(nil)

		// Actions
		userUseCase.ExecuteUpdateUserById(userId, payload)

		// Assert
		mockValidator.AssertExpectations(t)
		mockFileUpload.AssertExpectations(t)
		mockFileProcessing.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})
}
