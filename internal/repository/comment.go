package repository

import (
	"blog/internal/domain"
	"blog/internal/repository/dao"
	"context"
	"github.com/ecodeclub/ekit/slice"
)

type Comment interface {
	QueryByPost(ctx context.Context, id uint) ([]domain.Comment, error)
	Create(ctx context.Context, comment domain.Comment) (uint, error)
}

type CommentRepo struct {
	dao            dao.CommentDao
	userRepository UserRepository
}

func (c *CommentRepo) QueryByPost(ctx context.Context, id uint) ([]domain.Comment, error) {
	comments, err := c.dao.GetByPost(ctx, id)
	if err != nil {
		return nil, err
	}
	userName := ""
	user, err := c.userRepository.FindById(ctx, id)
	if err == nil {
		userName = user.Username
	}
	return slice.Map[dao.Comment, domain.Comment](comments, func(idx int, src dao.Comment) domain.Comment {
		return c.toDomain(userName, src)
	}), nil
}

func (c *CommentRepo) Create(ctx context.Context, comment domain.Comment) (uint, error) {
	id, err := c.dao.Insert(ctx, c.toEntity(comment))
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (c *CommentRepo) toDomain(userName string, comment dao.Comment) domain.Comment {
	return domain.Comment{
		Id:      comment.Id,
		Content: comment.Content,
		Author: domain.Author{
			Id:   comment.AuthorId,
			Name: userName,
		},
		PostId: comment.PostId,
	}
}

func (c *CommentRepo) toEntity(comment domain.Comment) dao.Comment {
	return dao.Comment{
		Content:  comment.Content,
		PostId:   comment.PostId,
		AuthorId: comment.Author.Id,
	}
}

func NewCommentRepo(dao dao.CommentDao, userRepository UserRepository) Comment {
	return &CommentRepo{dao: dao, userRepository: userRepository}
}
