package service

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"movie-rating-service/internal/application/models/request"
	"movie-rating-service/internal/application/models/response"
	"movie-rating-service/internal/domain"
	"movie-rating-service/internal/infrastructure/repository"
)

type UserService interface {
	Create(ctx context.Context, req request.CreateUser) (*response.CreateUser, error)
	Get(ctx context.Context, req request.GetUser) (*response.GetUser, error)
	IsAuthorized(ctx context.Context, req request.Login) (*response.GetUser, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (s *userService) Get(ctx context.Context, req request.GetUser) (*response.GetUser, error) {
	user, err := s.userRepository.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("error occurred while finding user: %w", err)
	}
	return user.GetUserResponse(), nil
}
func (s *userService) Create(ctx context.Context, req request.CreateUser) (*response.CreateUser, error) {
	user, err := s.userRepository.Create(ctx, domain.User{
		Username: req.Username,
		Password: req.Password,
		Name:     req.Name,
		Surname:  req.Surname,
		Email:    req.Email,
		Phone:    req.Phone,
		Address:  req.Address,
		IsAdmin:  req.IsAdmin,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user.CreateUserResponse(), nil
}

func (s *userService) IsAuthorized(ctx context.Context, req request.Login) (*response.GetUser, error) {
	user, err := s.userRepository.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("error occurred while logging in: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, err
	}

	return user.GetUserResponse(), nil
}
