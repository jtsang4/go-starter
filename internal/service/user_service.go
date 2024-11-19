package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jtsang4/go-stater/config"
	"github.com/jtsang4/go-stater/internal/model"
	"github.com/jtsang4/go-stater/internal/repository"
	"github.com/jtsang4/go-stater/pkg/auth"
	"github.com/jtsang4/go-stater/pkg/cache"
	"github.com/jtsang4/go-stater/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo  repository.UserRepositoryInterface
	cache cache.RedisCacheInterface
}

func NewUserService(repo repository.UserRepositoryInterface, cache cache.RedisCacheInterface) *UserService {
	return &UserService{
		repo:  repo,
		cache: cache,
	}
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Email    string `json:"email" binding:"required,email"`
}

func (s *UserService) CreateUser(req *CreateUserRequest) (*model.User, error) {
	// Check if username exists
	if _, err := s.repo.GetByUsername(req.Username); err == nil {
		return nil, errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	// 尝试从缓存获取
	var user *model.User
	cacheKey := fmt.Sprintf("user:%d", id)

	ctx := context.Background()
	err := s.cache.Get(ctx, cacheKey, &user)
	if err == nil {
		return user, nil
	}

	// 缓存未命中，从数据库获取
	user, err = s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 设置缓存
	if err := s.cache.Set(ctx, cacheKey, user, time.Hour); err != nil {
		logger.Logger.Warn("failed to set cache", zap.Error(err))
	}

	return user, nil
}

func (s *UserService) ValidateUser(username, password string) (*model.User, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	return user, nil
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (s *UserService) Login(req *LoginRequest, cfg config.JWTConfig) (string, error) {
	user, err := s.ValidateUser(req.Username, req.Password)
	if err != nil {
		return "", err
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, cfg)
	if err != nil {
		return "", err
	}

	return token, nil
}

type UpdateUserRequest struct {
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=6,max=32"`
}

func (s *UserService) UpdateUser(id uint, req *UpdateUserRequest) (*model.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Email != "" {
		user.Email = req.Email
	}

	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}
