package dao

import (
	"context"
	"gorm.io/gorm"
)

type Comment struct {
	Id       uint `gorm:"primaryKey, autoIncrement"`
	Content  string
	PostId   uint `gorm:"column:post_id"`
	AuthorId uint `gorm:"column:author_id"`
	gorm.DeletedAt
}

type CommentDao interface {
	GetByPost(ctx context.Context, id uint) ([]Comment, error)
	Insert(ctx context.Context, comment Comment) (uint, error)
}

type CommentGormDao struct {
	db *gorm.DB
}

func (c *CommentGormDao) GetByPost(ctx context.Context, id uint) ([]Comment, error) {
	var comments []Comment
	err := c.db.WithContext(ctx).Model(Comment{}).Where("post_id = ?", id).Find(&comments).Error
	return comments, err
}

func (c *CommentGormDao) Insert(ctx context.Context, comment Comment) (uint, error) {
	err := c.db.WithContext(ctx).Create(&comment).Error
	return comment.Id, err
}

func NewCommentGormDao(db *gorm.DB) CommentDao {
	return &CommentGormDao{db: db}
}
