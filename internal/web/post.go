package web

import (
	"blog/internal/domain"
	"blog/internal/service"
	"blog/internal/web/jwt"
	"blog/pkg/logger"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PostHandler struct {
	svc service.Post
	l   logger.LoggerV1
}

func NewPostHandler(svc service.Post, l logger.LoggerV1) *PostHandler {
	return &PostHandler{svc: svc, l: l}
}

func (h *PostHandler) RegisterRoute(server *gin.Engine) {
	group := server.Group("/posts")
	group.GET("", h.List)
	group.GET(":id", h.Detail)
	group.POST(":id", h.Edit)
	group.POST("", h.Create)
	group.DELETE(":id", h.Delete)
}

func (h *PostHandler) List(ctx *gin.Context) {
	userClaims := ctx.MustGet("user").(jwt.UserClaims)
	posts, err := h.svc.GetByAuthor(ctx.Request.Context(), userClaims.Uid)
	if err != nil {
		ctx.JSON(http.StatusOK, ErrSystemError)
		return
	}
	ctx.JSON(http.StatusOK, slice.Map[domain.Post, PostVo](posts, func(idx int, src domain.Post) PostVo {
		return PostVo{
			Id:           src.Id,
			Title:        src.Title,
			Content:      src.Content,
			AuthorId:     src.Author.Id,
			AuthorName:   src.Author.Name,
			CommentCount: src.CommentCount,
		}
	}))
}

func (h *PostHandler) Detail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusOK, "illegal number format")
		return
	}
	userClaims := ctx.MustGet("user").(jwt.UserClaims)
	post, err := h.svc.GetByPostId(ctx.Request.Context(), userClaims.Uid, uint(id))
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, PostVo{
		Id:           post.Id,
		Title:        post.Title,
		Content:      post.Content,
		AuthorId:     post.Author.Id,
		AuthorName:   post.Author.Name,
		CommentCount: post.CommentCount,
	})
}

func (h *PostHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Id      uint   `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var req Req
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, ErrSystemError)
		return
	}
	userClaims := ctx.MustGet("user").(jwt.UserClaims)
	id, err := h.svc.Save(ctx.Request.Context(), domain.Post{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author:  domain.Author{Id: userClaims.Uid},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, ErrSystemError)
		return
	}
	ctx.JSON(http.StatusOK, id)
}

func (h *PostHandler) Create(ctx *gin.Context) {
	type Req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var req Req
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, ErrSystemError)
		return
	}
	userClaims := ctx.MustGet("user").(jwt.UserClaims)
	id, err := h.svc.Save(ctx.Request.Context(), domain.Post{
		Title:   req.Title,
		Content: req.Content,
		Author:  domain.Author{Id: userClaims.Uid},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, ErrSystemError)
		return
	}
	ctx.JSON(http.StatusOK, id)
}

func (h *PostHandler) Delete(ctx *gin.Context) {
	type Req struct {
		Id uint `json:"id"`
	}
	var req Req
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, ErrSystemError)
		return
	}
	userClaims := ctx.MustGet("user").(jwt.UserClaims)
	err = h.svc.Delete(ctx.Request.Context(), userClaims.Uid, req.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, ErrSystemError)
		return
	}
	ctx.JSON(http.StatusOK, "success")
}
