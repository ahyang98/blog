package service

import (
	"blog/internal/domain"
	"blog/internal/repository"
	"blog/pkg/logger"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail        = repository.ErrDuplicateUser
	ErrInvalidUserOrPassword = errors.New("invalid user or password")
)

type UserService interface {
	SignUp(ctx context.Context, user *domain.User) error
	Login(ctx context.Context, email string, password string) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
	l    logger.LoggerV1
}

func (u *userService) SignUp(ctx context.Context, user *domain.User) error {
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(password)
	err = u.repo.Create(ctx, user)
	return err
}

func (u *userService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	user, err := u.repo.FindByEmail(ctx, email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return user, nil
}

func NewUserService(repo repository.UserRepository, l logger.LoggerV1) UserService {
	return &userService{repo: repo, l: l}
}
