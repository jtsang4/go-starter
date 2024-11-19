package service

import (
	"errors"
	"testing"

	"github.com/jtsang4/go-stater/config"
	"github.com/jtsang4/go-stater/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// Mock repository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id uint) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCache := new(MockCache)
	service := NewUserService(mockRepo, mockCache)

	tests := []struct {
		name    string
		req     *CreateUserRequest
		mock    func()
		wantErr bool
	}{
		{
			name: "success",
			req: &CreateUserRequest{
				Username: "testuser",
				Password: "password123",
				Email:    "test@example.com",
			},
			mock: func() {
				mockRepo.On("GetByUsername", "testuser").Return(nil, errors.New("not found"))
				mockRepo.On("Create", mock.AnythingOfType("*model.User")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "username exists",
			req: &CreateUserRequest{
				Username: "existinguser",
				Password: "password123",
				Email:    "test@example.com",
			},
			mock: func() {
				mockRepo.On("GetByUsername", "existinguser").Return(&model.User{}, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			user, err := service.CreateUser(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCache := new(MockCache)
	service := NewUserService(mockRepo, mockCache)

	tests := []struct {
		name    string
		id      uint
		mock    func()
		want    *model.User
		wantErr bool
	}{
		{
			name: "success from cache",
			id:   1,
			mock: func() {
				user := &model.User{ID: 1, Username: "testuser"}
				mockCache.On("Get", mock.Anything, "user:1", mock.AnythingOfType("*model.User")).
					Run(func(args mock.Arguments) {
						arg := args.Get(2).(*model.User)
						*arg = *user
					}).
					Return(nil)
			},
			want:    &model.User{ID: 1, Username: "testuser"},
			wantErr: false,
		},
		{
			name: "success from database",
			id:   2,
			mock: func() {
				user := &model.User{ID: 2, Username: "testuser2"}
				mockCache.On("Get", mock.Anything, "user:2", mock.AnythingOfType("*model.User")).
					Return(errors.New("cache miss"))
				mockRepo.On("GetByID", uint(2)).Return(user, nil)
				mockCache.On("Set", mock.Anything, "user:2", user, mock.Anything).Return(nil)
			},
			want:    &model.User{ID: 2, Username: "testuser2"},
			wantErr: false,
		},
		{
			name: "user not found",
			id:   3,
			mock: func() {
				mockCache.On("Get", mock.Anything, "user:3", mock.AnythingOfType("*model.User")).
					Return(errors.New("cache miss"))
				mockRepo.On("GetByID", uint(3)).Return(nil, errors.New("not found"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := service.GetUserByID(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCache := new(MockCache)
	service := NewUserService(mockRepo, mockCache)

	tests := []struct {
		name    string
		req     *LoginRequest
		mock    func()
		wantErr bool
	}{
		{
			name: "success",
			req: &LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			mock: func() {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				user := &model.User{
					ID:       1,
					Username: "testuser",
					Password: string(hashedPassword),
				}
				mockRepo.On("GetByUsername", "testuser").Return(user, nil)
			},
			wantErr: false,
		},
		{
			name: "invalid password",
			req: &LoginRequest{
				Username: "testuser",
				Password: "wrongpassword",
			},
			mock: func() {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				user := &model.User{
					ID:       1,
					Username: "testuser",
					Password: string(hashedPassword),
				}
				mockRepo.On("GetByUsername", "testuser").Return(user, nil)
			},
			wantErr: true,
		},
		{
			name: "user not found",
			req: &LoginRequest{
				Username: "nonexistent",
				Password: "password123",
			},
			mock: func() {
				mockRepo.On("GetByUsername", "nonexistent").Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			token, err := service.Login(tt.req, config.JWTConfig{Secret: "test", ExpireTime: 24})
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}
