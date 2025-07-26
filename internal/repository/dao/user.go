package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var (
	ErrDuplicateUser  = errors.New("email duplicated")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

const duplicateErr uint16 = 1062

type User struct {
	Id       uint                 `gorm:"primaryKey, autoIncrement"`
	Email    datatypes.NullString `gorm:"column:email;unique"`
	Password datatypes.NullString `gorm:"column:password"`
	Username datatypes.NullString `gorm:"column:username;unique"`
}

type UserDao interface {
	Insert(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindById(ctx context.Context, id uint) (User, error)
}

type GormUserDao struct {
	db *gorm.DB
}

func (g *GormUserDao) FindById(ctx context.Context, id uint) (User, error) {
	var user User
	err := g.db.WithContext(ctx).Model(&User{}).Where("id=?", id).First(&user).Error
	return user, err
}

func (g *GormUserDao) Insert(ctx context.Context, user *User) error {
	err := g.db.WithContext(ctx).Create(user).Error
	if sqlError, ok := err.(*mysql.MySQLError); ok {
		if sqlError.Number == duplicateErr {
			return ErrDuplicateUser
		}
	}
	return err
}

func (g *GormUserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := g.db.WithContext(ctx).Model(&User{}).Where("email=?", email).First(&user).Error
	return user, err
}

func NewGormUserDao(db *gorm.DB) UserDao {
	return &GormUserDao{db: db}
}
