package service

import (
	"blog/internal/domain"
	"blog/internal/repository"
	"blog/pkg/logger"
	"context"
)

type Post interface {
	GetByAuthor(ctx context.Context, uid uint) (posts []domain.Post, err error)
	GetByPostId(ctx context.Context, uid uint, postId uint) (posts domain.Post, err error)
	Save(ctx context.Context, post domain.Post) (id uint, err error)
	Delete(ctx context.Context, uid uint, id uint) error
}

type PostService struct {
	repo repository.Post
	l    logger.LoggerV1
}

func (p *PostService) GetByAuthor(ctx context.Context, uid uint) (posts []domain.Post, err error) {
	return p.repo.GetByAuthor(ctx, uid)
}

func (p *PostService) GetByPostId(ctx context.Context, uid uint, postId uint) (posts domain.Post, err error) {
	return p.repo.GetByPostId(ctx, uid, postId)
}

func (p *PostService) Save(ctx context.Context, post domain.Post) (id uint, err error) {
	if post.Id > 0 {
		return post.Id, p.repo.Update(ctx, post)
	}
	return p.repo.Create(ctx, post)
}

func (p *PostService) Delete(ctx context.Context, uid uint, id uint) error {
	return p.repo.Delete(ctx, uid, id)
}

func NewPostService(repo repository.Post, l logger.LoggerV1) Post {
	return &PostService{repo: repo, l: l}
}
