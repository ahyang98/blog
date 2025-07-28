package service

import (
	"blog/internal/domain"
	"blog/internal/repository"
	"blog/pkg/logger"
	"context"
)

type Comment interface {
	QueryByPost(ctx context.Context, postId uint) ([]domain.Comment, error)
	Save(ctx context.Context, comment domain.Comment) (uint, error)
}

type CommentService struct {
	repo    repository.Comment
	l       logger.LoggerV1
	postSvc Post
}

func (c *CommentService) QueryByPost(ctx context.Context, postId uint) ([]domain.Comment, error) {
	return c.repo.QueryByPost(ctx, postId)
}

func (c *CommentService) Save(ctx context.Context, comment domain.Comment) (uint, error) {
	commentId, err := c.repo.Create(ctx, comment)
	if err != nil {
		return 0, err
	}
	c.postSvc.IncreaseComment(ctx, comment.PostId)
	return commentId, nil
}

func NewCommentService(repo repository.Comment, post Post, l logger.LoggerV1) Comment {
	return &CommentService{repo: repo, postSvc: post, l: l}
}
