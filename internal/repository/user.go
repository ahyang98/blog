package repository

import (
	"blog/internal/domain"
	"blog/internal/repository/dao"
	"context"
	"gorm.io/datatypes"
)

var (
	ErrDuplicateUser = dao.ErrDuplicateUser
	ErrUserNotFound  = dao.ErrRecordNotFound
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
}

type userRepository struct {
	userDao dao.UserDao
}

func NewUserRepository(userDao dao.UserDao) UserRepository {
	return &userRepository{userDao: userDao}
}

func (d *userRepository) Create(ctx context.Context, user *domain.User) error {
	return d.userDao.Insert(ctx, d.toEntity(user))
}

func (d *userRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := d.userDao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return toDomainUser(user), nil
}

func (d *userRepository) toEntity(user *domain.User) *dao.User {
	//now := time.Now()
	return &dao.User{
		//Model: gorm.Model{
		//	ID:        0,
		//	CreatedAt: now,
		//	UpdatedAt: now,
		//},
		Id:       user.Id,
		Email:    String2SqlNullString(user.Email),
		Password: String2SqlNullString(user.Password),
		Username: String2SqlNullString(user.Username),
	}
}

func String2SqlNullString(field string) datatypes.NullString {
	return datatypes.NullString{
		V:     field,
		Valid: len(field) > 0,
	}
}

func toDomainUser(user dao.User) domain.User {
	return domain.User{
		Id:       user.Id,
		Email:    user.Email.V,
		Password: user.Password.V,
		Username: user.Username.V,
	}
}
