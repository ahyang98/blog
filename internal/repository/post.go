package repository

import (
	"blog/internal/domain"
	"blog/internal/repository/dao"
	"context"
	"github.com/ecodeclub/ekit/slice"
)

type Post interface {
	GetByAuthor(ctx context.Context, uid uint) ([]domain.Post, error)
	GetByPostId(ctx context.Context, uid uint, id uint) (domain.Post, error)
	Update(ctx context.Context, post domain.Post) error
	Create(ctx context.Context, post domain.Post) (uint, error)
	Delete(ctx context.Context, uid uint, id uint) error
	IncreaseComment(ctx context.Context, id uint)
}

type PostRepository struct {
	dao            dao.PostDao
	userRepository UserRepository
}

func (p *PostRepository) IncreaseComment(ctx context.Context, id uint) {
	p.dao.Add(ctx, id)
}

func (p *PostRepository) GetByAuthor(ctx context.Context, uid uint) ([]domain.Post, error) {
	posts, err := p.dao.GetByAuthor(ctx, uid)
	if err != nil {
		return nil, err
	}
	userName := p.getUserName(ctx, uid)
	res := slice.Map[dao.Post, domain.Post](posts, func(idx int, src dao.Post) domain.Post {
		return p.toDomain(userName, src)
	})
	return res, nil
}

func (p *PostRepository) getUserName(ctx context.Context, uid uint) string {
	userName := ""
	user, err := p.userRepository.FindById(ctx, uid)
	if err == nil {
		userName = user.Username
	}
	return userName
}

func (p *PostRepository) GetByPostId(ctx context.Context, uid uint, id uint) (domain.Post, error) {
	post, err := p.dao.GetById(ctx, uid, id)
	if err != nil {
		return domain.Post{}, err
	}
	res := p.toDomain(p.getUserName(ctx, uid), post)
	return res, nil
}

func (p *PostRepository) Update(ctx context.Context, post domain.Post) error {
	return p.dao.Update(ctx, p.toEntity(post))
}

func (p *PostRepository) Create(ctx context.Context, post domain.Post) (uint, error) {
	post.CommentCount = 0
	return p.dao.Insert(ctx, p.toEntity(post))
}

func (p *PostRepository) Delete(ctx context.Context, uid uint, id uint) error {
	return p.dao.Delete(ctx, uid, id)
}

func (p *PostRepository) toDomain(userName string, src dao.Post) domain.Post {
	return domain.Post{
		Id:      src.Id,
		Title:   src.Title,
		Content: src.Content,
		Author: domain.Author{
			Id:   src.AuthorId,
			Name: userName,
		},
		CommentCount: src.CommentCount,
	}
}

func (p *PostRepository) toEntity(post domain.Post) dao.Post {
	return dao.Post{
		Id:           post.Id,
		Title:        post.Title,
		Content:      post.Content,
		AuthorId:     post.Author.Id,
		CommentCount: post.CommentCount,
	}
}

func NewPostRepository(dao dao.PostDao, userRepository UserRepository) Post {
	return &PostRepository{dao: dao, userRepository: userRepository}
}
