package dao

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Post struct {
	Id           uint `gorm:"primaryKey, autoIncrement"`
	Title        string
	Content      string
	AuthorId     uint `gorm:"column:author_id"`
	CommentCount uint `gorm:"column:comment_count"`
	gorm.DeletedAt
}

type PostDao interface {
	GetByAuthor(ctx context.Context, uid uint) ([]Post, error)
	GetById(ctx context.Context, uid uint, id uint) (Post, error)
	Update(ctx context.Context, post Post) error
	Insert(ctx context.Context, post Post) (uint, error)
	Delete(ctx context.Context, uid uint, id uint) error
	Add(ctx context.Context, id uint)
}

type PostDaoGorm struct {
	db *gorm.DB
}

func (p *PostDaoGorm) Add(ctx context.Context, id uint) {
	var post Post
	p.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.Locking{
			Strength: "UPDATE",
		}).First(&post, id).Error
		if err != nil {
			return err
		}
		return tx.Exec("update posts set CommentCount = CommentCount + 1 where id = ?", id).Error
	})
}

func (p *PostDaoGorm) GetByAuthor(ctx context.Context, uid uint) ([]Post, error) {
	var posts []Post
	err := p.db.WithContext(ctx).Model(&Post{}).Where("author_id", uid).Find(&posts).Error
	return posts, err
}

func (p *PostDaoGorm) GetById(ctx context.Context, uid uint, id uint) (Post, error) {
	var post Post
	res := p.db.WithContext(ctx).Model(&Post{}).Where("author_id", uid).First(&post)
	if res.Error != nil {
		return post, res.Error
	}
	if res.RowsAffected == 0 {
		return post, ErrRecordNotFound
	}
	return post, nil
}

func (p *PostDaoGorm) Update(ctx context.Context, post Post) error {
	res := p.db.WithContext(ctx).Model(Post{}).Where("id = ? and author_id = ?", post.Id, post.AuthorId).
		Updates(map[string]any{
			"delete_at": sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (p *PostDaoGorm) Insert(ctx context.Context, post Post) (uint, error) {
	err := p.db.WithContext(ctx).Create(&post).Error
	return post.Id, err
}

func (p *PostDaoGorm) Delete(ctx context.Context, uid uint, id uint) error {
	res := p.db.WithContext(ctx).Model(Post{}).Where("id = ? and author_id = ?", id, uid).Updates(map[string]any{
		"delete_at": sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func NewPostDaoGorm(db *gorm.DB) PostDao {
	return &PostDaoGorm{db: db}
}
