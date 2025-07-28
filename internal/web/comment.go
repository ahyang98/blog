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

type CommentHandler struct {
	svc service.Comment
	l   logger.LoggerV1
}

func NewCommentHandler(svc service.Comment, l logger.LoggerV1) *CommentHandler {
	return &CommentHandler{svc: svc, l: l}
}

func (h *CommentHandler) RegisterRoute(server *gin.Engine) {
	group := server.Group("comment")
	group.GET("", h.ListByPost)
	group.POST("", h.Create)
}

func (h *CommentHandler) ListByPost(ctx *gin.Context) {
	var postId uint64
	if postIdStr, exists := ctx.GetQuery("post_id"); !exists {
		ctx.JSON(http.StatusOK, "error parameter")
		return
	} else {
		var err error
		postId, err = strconv.ParseUint(postIdStr, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusOK, ErrSystemError)
			return
		}
	}
	comments, err := h.svc.QueryByPost(ctx.Request.Context(), uint(postId))
	if err != nil {
		ctx.JSON(http.StatusOK, ErrSystemError)
		return
	}

	ctx.JSON(http.StatusOK, slice.Map[domain.Comment, CommentVo](comments, func(idx int, src domain.Comment) CommentVo {
		return CommentVo{
			Id:         src.Id,
			Content:    src.Content,
			AuthorId:   src.Author.Id,
			AuthorName: src.Author.Name,
		}
	}))

}

func (h *CommentHandler) Create(ctx *gin.Context) {
	type Req struct {
		Content string `json:"content"`
		PostId  uint   `json:"post_id"`
	}
	var req Req
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, ErrSystemError)
		return
	}
	userClaims := ctx.MustGet("user").(jwt.UserClaims)
	id, err := h.svc.Save(ctx.Request.Context(), domain.Comment{
		Content: req.Content,
		Author:  domain.Author{Id: userClaims.Uid},
		PostId:  req.PostId,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, ErrSystemError)
		return
	}
	ctx.JSON(http.StatusOK, id)
}
