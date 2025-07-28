package ioc

import (
	"blog/internal/web"
	ijwt "blog/internal/web/jwt"
	"blog/internal/web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler, postHdl *web.PostHandler, commentHdl *web.CommentHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	postHdl.RegisterRoute(server)
	commentHdl.RegisterRoute(server)
	return server
}

func InitGinMiddlewares(c *cache.Cache, ijwtHdl ijwt.Handler) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.NewLoginJWTMiddlewareBuilder(c, ijwtHdl).CheckLogin(),
	}
}
