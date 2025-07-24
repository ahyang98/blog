package middleware

import (
	ijwt "blog/internal/web/jwt"
	"github.com/gin-gonic/gin"
	cache "github.com/patrickmn/go-cache"
	"net/http"
)

type LoginJWTMiddlewareBuilder struct {
	jwtHdl ijwt.Handler
	c      *cache.Cache
	paths  map[string]string
}

func NewLoginJWTMiddlewareBuilder(c *cache.Cache, jwtHdl ijwt.Handler) *LoginJWTMiddlewareBuilder {
	m := make(map[string]string)
	m["/users/signup"] = "/users/signup"
	m["/users/login"] = "/users/login"
	m["/users/refresh_token"] = "/users/refresh_token"
	return &LoginJWTMiddlewareBuilder{c: c, jwtHdl: jwtHdl, paths: m}
}

func (b *LoginJWTMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if _, ok := b.paths[path]; ok {
			return
		}
		uc, err := b.jwtHdl.GetUserClaim(ctx)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		err = b.jwtHdl.CheckSession(ctx, uc.Ssid)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set("user", uc)
	}
}
